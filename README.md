# Gitea Deployment Simplified: AWS VPC, NFS Storage Gateway, and TUI Automation

## 🌟 Introduction

Welcome to the ultimate solution for deploying Gitea, a lightweight and self-hosted Git service, on AWS! This project combines state-of-the-art technologies to deliver a secure, scalable, and fully automated deployment experience.
🚀 What is this project about?

This project uses a combination of advanced tools and services to deploy a private Git hosting solution:

    🌐 AWS VPC: Creates an isolated, secure environment for your resources.
    📦 NFS Storage Gateway: Provides scalable, reliable, and persistent storage for repositories and files.
    🛠️ Terraform: Automates the creation of AWS infrastructure like EC2, security groups, and more.
    ⚙️ Ansible: Manages the seamless installation and configuration of Gitea.
    🎨 Bubble Tea TUI: Provides an intuitive, interactive interface for configuring and deploying your setup with ease.

❓ What problem does it solve?

Managing a self-hosted Git service is often challenging due to the need for:

    🔒 Security: Ensuring your repositories are private and protected.
    📈 Scalability: Handling increasing storage and performance demands.
    ⏱️ Efficiency: Reducing manual setup and repetitive configuration tasks.

This project solves these issues by automating the entire process, creating a highly secure and scalable setup with minimal effort.
💡 Why is it useful?

This project is designed to simplify your life while ensuring top-notch performance and security:

    🔑 Comprehensive Deployment: A one-stop solution for secure and private Git hosting on AWS.
    🤝 User-Friendly: The interactive TUI interface makes setup simple for everyone.
    ✨ Feature-Rich: Supports NFS-based storage, private VPC networking, and customizable configurations.
    💵 Cost-Effective: Optimized infrastructure to minimize unnecessary expenses.

Whether you’re an individual developer or part of a large organization, this project delivers everything you need to securely host Git repositories while scaling effortlessly as your needs grow.


## 🛠️ Installation

Follow these steps to set up and run your project successfully:
📋 Prerequisites

🔗 AWS CLI
        Ensure the AWS CLI is installed and configured with appropriate credentials.
        Run aws configure to set the default region and credentials.

    ⚠️ Important Note: The ACM certificates displayed in the TUI are region-specific. Make sure your AWS CLI is configured to the region where your ACM certificates exist (e.g., us-east-1). If you switch to a different region, ensure you have ACM certificates available in that region.

🔗 Terraform       
Install Terraform from the official Terraform website.
        Verify the installation with:

    terraform --version

🔗 `uvx`

Install [`uv`](https://github.com/astral-sh/uv)

```sh
curl -LsSf https://astral.sh/uv/install.sh | sh 
```


🔗 Bubble Tea TUI:
Install the required Go dependencies for the TUI:

    go mod tidy  

⚠️ Disclaimers

    ACM Region-Specific Certificates:
    Ensure your AWS CLI is set to the correct region where your ACM certificate is located, as it is not shared across regions.




▶️ Setup and Run
1. Clone the Repository

Download the project to your local machine:
```
git clone https://github.com/azransyed/Gitea-SGW-TUI.git  
cd Gitea-Terraform-Ansible-TUI
```

2. Run the Bubble Tea TUI

With the virtual environment activated, navigate to the project directory and start the TUI:
```
go mod tidy 
go run main.go  
```


🎯 What Happens Next?

The TUI will guide you through configuring the AWS setup with the following inputs:

    1.Instance Type: Choose the instance type for your EC2 instance (e.g., t2.micro).
    2.Public IP of the User's PC: Enter your PC's public IP for activating the NFS Storage Gateway.
    3.Region: Specify the AWS region where the VPC will be deployed (e.g., us-east-1).
    4.Gitea Base Domain: Provide the base domain for accessing the Gitea website (e.g., gitea.example.com).
    5.OIDC Authentication (OAuth2 Setup): Configure authentication using OAuth2 and Instructions are provided below.
    6.KMS ARN (Optional): If using AWS KMS for encryption, provide the KMS ARN. If not required, leave this blank.
    7.ACM ARN Certificate: Select the ARN of the ACM certificate for securing the website.


🔑 Setting Up OIDC Authentication

To enable OAuth2 authentication, you need to enter the following details correctly:
 
    1. Authentication Name / OAuth Application Name
    - This should match the identity provider configuration. If using Auth0, enter `auth0`.
    
    2. OAuth2 Provider
    - The type of OAuth provider. For this setup, use `openidConnect` enter this value in same format.

    3. Client ID (Key)
     - This is the unique identifier provided by your OAuth2 provider when you create an application.

    4. Client Secret
    - This is the secret key provided by the OAuth2 provider, used to authenticate requests.

    5. Auto Discover URL
    - This is the full URL required for auto-discovery of authentication settings.
    - Example: `https://login.example.org/.well-known/openid-configuration`
    - Typically, the provider only gives you the middle part (`login.example.org`), so make sure you enter the full URL correctly.
    

🔗 Tips for OAuth2 Setup:
https://docs.gitea.com/development/oauth2-provider#endpoints


🚨 Important Notice:

1.  If you enter incorrect values,the Ansible job will fail and use the same formating as given in the example. Double-check all values before proceeding.

2.  An admin user will also be created during setup.

3.  You can modify the default admin username and password in ansible/playbook.yaml.gotmpl before running the playbook.

Following these steps correctly ensures a smooth setup process with authentication enabled! 🚀
    


⚙️ Important Once You Complete the Configuration

    Terraform:
    The tool will deploy the necessary AWS infrastructure, including the VPC, EC2 instance, Application Load Balancer (ALB), NFS Storage Gateway, and other resources.

    Ansible:
    The tool will configure Gitea with all required settings, including attaching the NFS storage, setting up the domain, and securing it with the ACM certificate.

    DNS Configuration (GoDaddy Setup):
    - Once Terraform and Ansible have completed their tasks, go to the AWS Management Console.
    - Navigate to the Load Balancer section and locate the ALB that was created.
    - Copy the DNS name of the ALB.
    - Go to your GoDaddy DNS settings and create a CNAME record, pointing your domain (gitea.example.com) to the ALB's DNS name.
    - This step ensures that your domain resolves correctly to the Gitea instance.
