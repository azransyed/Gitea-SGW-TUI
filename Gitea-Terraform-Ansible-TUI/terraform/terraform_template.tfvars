# Template file for terraform vars this template is a refernce template for creating the Terraform tfvars

instance_type                 = "{{ .InstanceType }}"
ingress_cidr_block_activation = "{{ .PublicIP }}"
region                        = "{{ .Region }}"
acm_certificate_arn           = "{{ .CertificateArn }}"
