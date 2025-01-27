package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Custom List item
type certItem struct {
	arn string
}

func (c certItem) Title() string       { return c.arn }
func (c certItem) Description() string { return "Certificate ARN" }
func (c certItem) FilterValue() string { return c.arn }

// Struct to Fetch the certificate from aws
type fetchCertsMsg struct {
	Certificates []string
}

// Struct for Config Data
type ConfigData struct {
	InstanceType   string
	CertificateArn string
	PublicIP       string
	Region         string
	Domain         string
	RootURL        string
	SsmKmsKeyArn   *string
}

// for Error Message
type errMsg error

type model struct {
	textInput     textinput.Model
	instanceType  string
	step          int
	selectedARN   string
	kmsARN        *string
	err           error
	list          list.Model
	keyMap        KeyMap
	publicIP      string
	region        string
	domain        string
	rootURL       string
	spinner       spinner.Model
	statusMessage string
	quitting      bool
}

// Struct for the keys which such as arrow keys and enter keys which will be used in Selection of AWS ACM Certificate ARN
type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Esc   key.Binding
}

// type sturct for terraform output
type TerraformOutput struct {
	GiteaInstanceID struct {
		Value string `json:"value"`
	} `json:"Gitea_instance_id"`
	NfsSharePath struct {
		Value string `json:"value"`
	} `json:"nfs_share_path"`
	Region struct {
		Value string `json:"value"`
	} `json:"region"`
	StorageGatewayIP struct {
		Value string `json:"value"`
	} `json:"storage_gateway_private_ip"`
}

// for the placeholder and input style for TUI with color
var (
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("147"))            //146 and 147
	inputStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))             // 69
	headingStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true) // 204 Red like teal peach
	pinkBold         = "\033[1;38;5;213m"
	reset            = "\033[0m"
)

// To Start the TUI
func main() {
	p := tea.NewProgram(initialModel())
	_m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	m := _m.(model)

	if m.quitting {
		return
	}

	// Display message afterthe TUI exits
	fmt.Printf("\n%sTUI exited. Starting Terraform setup...%s\n", pinkBold, reset)
	//fmt.Println("\nTUI exited. Starting Terraform setup...")

	//Step 1: Run `terraform init`
	if err := runTerraform("init"); err != nil {
		fmt.Printf("\nError during %sterraform init%s: %v\n", pinkBold, reset, err)
		os.Exit(1)

	}

	// Step 2: Run 'terraform apply -auto-approve'
	if err := runTerraform("apply"); err != nil {
		fmt.Printf("\nError during %sterraform apply%s: %v\n", pinkBold, reset, err)
		os.Exit(1)
	}

	tfOutput, err := getTerraformOutput()
	if err != nil {
		fmt.Printf("\nError getting Terraform output: %v\n", err)
		os.Exit(1)
	}
	replacements := map[string]string{
		"Region":       tfOutput.Region.Value,
		"InstanceID":   tfOutput.GiteaInstanceID.Value,
		"NfsSharePath": tfOutput.NfsSharePath.Value,
		"NfsServer":    tfOutput.StorageGatewayIP.Value,
		"Fqdn":         m.domain,
		"RootURL":      m.rootURL,
	}
	inventoryPath := "ansible/inventory_aws_ec2.yaml.gotmpl"
	playbookPath := "ansible/playbook.yaml.gotmpl"

	if err := updateFile(inventoryPath, replacements); err != nil {
		fmt.Printf("\nError updating inventory file: %v\n", err)
		os.Exit(1)
	}
	if err := updateFile(playbookPath, replacements); err != nil {
		fmt.Printf("\nError updating playbook file: %v\n", err)
		os.Exit(1)
	}

	// Completion message for Terraform setup
	fmt.Printf("\n%sTerraform setup completed sucessfully 🎉 🎇%s\n", pinkBold, reset)

	// Waits for 2 Minutes before Running Ansible Playbook
	fmt.Printf("\n%sWaiting for 2 minutes before running the Ansible playbook 👻 🕙...%s\n", pinkBold, reset)
	time.Sleep(2 * time.Minute) // Change kar haan

	// Run the Ansible Playbook
	if err := runAnsiblePlaybook(); err != nil {
		fmt.Printf("\nError during %sAnsible Playbook%s execution: %v\n", pinkBold, reset, err)
		os.Exit(1)

	}

	fmt.Printf("\n%sAnsible playbook executed successfully 🎉 🎇%s\n", pinkBold, reset)

}

// First Page of TUI
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "eg: t2.micro"
	ti.Focus()
	ti.Width = 30
	ti.PlaceholderStyle = placeholderStyle
	ti.TextStyle = inputStyle

	items := []list.Item{}
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 90, 10)

	return model{
		textInput: ti,
		list:      l,
		step:      1,
		keyMap: KeyMap{
			Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
			Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
			Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("↵", "select")),
			Esc:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
		},
	}
}

// Custom delegate for full-width text display
type customDelegate struct{}

func (d customDelegate) Height() int                             { return 1 }
func (d customDelegate) Spacing() int                            { return 0 }
func (d customDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d customDelegate) Render(w, h int, item list.Item, _ bool) string {
	return item.FilterValue() // Render the full text without any cut of value shown
}

// THE BLINK we will add the color
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update Function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	//Handle Keyboard inputs
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			switch m.step {
			case 1:
				//	m.textInput.Placeholder = "Enter your desired EC2 instance type"
				m.instanceType = m.textInput.Value()
				m.step = 2
				m.textInput.SetValue("")
				m.textInput.Placeholder = "eg: 82.129.80.111/32"
				return m, nil
			case 2:
				m.publicIP = m.textInput.Value()
				m.step = 3
				m.textInput.SetValue("")
				m.textInput.Placeholder = "eg: us-east-1"
				return m, nil
			case 3:
				m.region = m.textInput.Value()
				m.step = 4
				m.textInput.SetValue("")
				m.textInput.Placeholder = "eg: example.com)"
				return m, nil
			case 4:
				m.domain = m.textInput.Value()
				m.textInput.Reset() // Clear input for the next use
				m.step = 5
				m.textInput.SetValue("")
				m.textInput.Placeholder = "eg: arn:aws:acm:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
				return m, nil
			case 5:
				kmsARN := m.textInput.Value()
				if kmsARN != "" {
					m.kmsARN = &kmsARN
				}
				m.textInput.Reset() // Clear input for the next use
				m.step = 6
				m.textInput.Blur()
				return m, fetchCertificateCmd() // Fetch certificate
			case 6:
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					m.selectedARN = selectedItem.FilterValue()
					//Populate ConfigData with Instnace type and selected ARN
					data := ConfigData{
						InstanceType:   m.instanceType,
						CertificateArn: m.selectedARN,
						PublicIP:       m.publicIP,
						Region:         m.region,
						Domain:         m.domain,
						RootURL:        fmt.Sprintf("http://%s/", m.domain),
						SsmKmsKeyArn:   m.kmsARN,
					}
					m.textInput.PlaceholderStyle = placeholderStyle
					m.textInput.TextStyle = inputStyle

					// Call a function to generate the terraform file

					if err := generateTerraformFile(data); err != nil {
						fmt.Println("Error generating Terraform file:", err)
						m.quitting = true
						return m, tea.Quit
					}
					return m, tea.Quit

				}
			}

		case key.Matches(msg, m.keyMap.Esc):
			m.quitting = true
			return m, tea.Quit // Quit on esc

		}

	///////////////////////////////////////
	case fetchCertsMsg:
		items := make([]list.Item, len(msg.Certificates))
		for i, cert := range msg.Certificates {
			items[i] = certItem{arn: cert}
		}
		m.list.SetItems(items)
		m.step = 5

	case errMsg:
		m.err = msg //Handle error
		return m, tea.Quit
	}

	// Update LIst or text or text input depending on the step
	if m.step < 6 {
		m.textInput, cmd = m.textInput.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd

}

// Display Differnt views based on steps
func (m model) View() string {
	switch m.step {
	case 1:
		return fmt.Sprintf(
			headingStyle.Render("Enter your desired instance type 😅🖥️ :\n\n%s\n\n%s"),
			m.textInput.View(),
			"(Press Enter to confirm, Esc to quit)",
		)
	case 2:
		return fmt.Sprintf(
			headingStyle.Render("Enter your Public IP (e.g., 12.23.45.67/32) add /32 🙌🌐:\n\n%s\n\n%s"), // we need to discaler to add /32
			m.textInput.View(),
			"(Press Enter to confirm, Esc to quit)",
		)
	case 3:
		return fmt.Sprintf(
			headingStyle.Render("Enter your AWS Region (eg., us-east-1) 😏📍:\n\n%s\n\n%s"),
			m.textInput.View(),
			"(Press Enter to confirm, Esc to quit)",
		)

	case 4:
		return fmt.Sprintf(
			headingStyle.Render("Enter Gitea base domain (e.g., example.com) 😄🌸:\n\n%s\n\n%s"),
			m.textInput.View(),
			"(Press Enter to confirm, Esc to quit)",
		)
	case 5:
		return fmt.Sprintf(
			headingStyle.Render("Enter KMS ARN (optional) 🔒:\n\n%s\n\n%s"),
			m.textInput.View(),
			"(Press Enter to confirm, Esc to quit)",
		)
	case 6:
		return fmt.Sprintf(
			headingStyle.Render("Select a certificate ARN 🔒:\n\n%s\n\n%s"),
			m.list.View(),
			"(Use ↑/↓ to navigate, Enter to select, Esc to quit)",
		)
	default:
		return "Terraform file created successfully! Press Esc to exit."
	}
}

// Command to Fetch ARN's
func fetchCertificateCmd() tea.Cmd {
	return func() tea.Msg {
		arns, err := fetchCertificates() // Calls the AWS SDK here
		if err != nil {
			return errMsg(err)
		}
		return fetchCertsMsg{Certificates: arns}
	}
}

// Function to Fetch the AWS ACM Certificate
func fetchCertificates() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := acm.NewFromConfig(cfg) // Assuming AWS config is set
	output, err := client.ListCertificates(context.TODO(), &acm.ListCertificatesInput{})
	if err != nil {
		return nil, err
	}

	var arns []string
	for _, cert := range output.CertificateSummaryList {
		arns = append(arns, *cert.CertificateArn)
	}

	return arns, nil
}

// Function which will Generate Terraform tf vars file which will be used by terraform the data are user inputs
func generateTerraformFileCmd(data ConfigData) tea.Cmd {
	return func() tea.Msg {
		err := generateTerraformFile(data)
		if err != nil {
			fmt.Printf("Error generating Terraform file: %v\n", err)
			return errMsg(err)
		}

		fmt.Println("\nTerraform file generated successfully at terraform/terraform.tfvars!")

		return tea.Quit // Exit after file generation
	}
}

// Function to Generate Terraform File after the input from user
func generateTerraformFile(data ConfigData) error {
	return updateFile("terraform/terraform.tfvars.gotmpl", data)
}

// //////////////////////////////////////////////////////////
// Function to runTerraform execution
func runTerraform(command string) error {

	// Define bold pink style with ANSI escape code
	pinkBold := "\033[1;38;5;213m"
	reset := "\033[0m"

	// The full command which shows in terminal for terraform
	fullCommand := fmt.Sprintf("%sterraform %s ✨%s", pinkBold, command, reset)

	// To Print the Styled Command
	fmt.Printf("\nRunning: %s\n\n", fullCommand)

	// Prepare and executes the command
	cmd := exec.Command("terraform", strings.Fields(command)...)
	cmd.Dir = "./terraform"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//fmt.Printf("\nRunning: %s\n\n", fullCommand)

	// Run the command and return any error encounterd
	return cmd.Run()
}

// function to get Terraform output
func getTerraformOutput() (*TerraformOutput, error) {
	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = "./terraform" // To ensure we are in the correct directory
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get Terraform output: %w", err)
	}

	var tfOutput TerraformOutput
	if err := json.Unmarshal(output, &tfOutput); err != nil {
		return nil, fmt.Errorf("failed to parse Terraform output: %w", err)
	}
	return &tfOutput, nil
}

// Function to Update the Ansible playbook
func updateFile(filepath string, replacements any) error {
	// Check if the template file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("template file %s does not exist", filepath)
	}

	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return fmt.Errorf("error parsing template file: %w", err)
	}

	outputFile, err := os.Create(strings.TrimSuffix(filepath, ".gotmpl"))
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	// Execute the template with the provided data
	err = tmpl.Execute(outputFile, replacements)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

// Function which runs the Ansible playbook command
func runAnsiblePlaybook() error {
	// Styled Playbook command Message
	fullCommand := fmt.Sprintf("%sansible-playbook -i inventory_aws_ec2.yaml playbook.yaml -vvv ✨%s", pinkBold, reset)

	// Print the styled message
	fmt.Printf("\nRunning: %s\n\n", fullCommand)

	// Prepare and execute the playbook command
	cmd := exec.Command("bash", "-c", "ansible-playbook -i inventory_aws_ec2.yaml playbook.yaml -vvv ")
	cmd.Dir = "./ansible"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
