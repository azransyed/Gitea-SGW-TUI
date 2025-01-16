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

🔗 Python

    Install Python (version 3.8 or higher).
    Set up a virtual environment:

    python3 -m venv venv  
    source venv/bin/activate  # On Linux/Mac  
    venv\Scripts\activate     # On Windows  

🔗 Ansible

    Install Ansible and its dependencies:

        pip install ansible boto3 botocore  

🔗 Dependencies

    Python Libraries:
    Make sure to install the following Python libraries inside your virtual environment:
    pip install boto3 botocore ansible

  

🔗 Bubble Tea TUI:
Install the required Go dependencies for the TUI:

    go mod tidy  

⚠️ Disclaimers

    Enable Virtual Environment:
    Always activate your Python virtual environment before running the TUI or executing any Ansible playbooks to avoid dependency conflicts.

    ACM Region-Specific Certificates:
    Ensure your AWS CLI is set to the correct region where your ACM certificate is located, as it is not shared across regions.

▶️ Setup and Run

    Clone the repository:
