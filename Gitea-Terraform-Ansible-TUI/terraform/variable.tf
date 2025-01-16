# Variable for Instance Type
variable "instance_type" {
  description = "EC2 instance type"
}
#Variable for ingress_cidr_block_activation
variable "ingress_cidr_block_activation" {
  description = "User Public IP"

}
#Variable for Region
variable "region" {
  description = "Region for AWS"
}

# Variable for ACM Certificate ARN
variable "acm_certificate_arn" {
  description = "The ARN of the ACM certificate"
}

# Variable for ssh public key path either define the key path from your Terrform Runner Means your PC or inside the terraform Directory
variable "ssh_public_key_path" {
  type        = string
  description = "(Optional) Absolute file path to the the public key for the EC2 Key pair. If ommitted, the EC2 key pair resource will not be created"
  default     = ""
}

# Variable for ssh key name either keep default empty "" or put in your own custom name
variable "ssh_key_name" {
  type        = string
  description = "(Optional) The name of an existing EC2 Key pair for SSH access to the EC2 Storage Gateway"
  default     = ""
}

//sgw-key

