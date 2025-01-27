data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
  filter {
    name   = "name"
    values = ["al2023-ami-2023*"]
  }
}

# Create the EC2 Instance in the Private Subnet
resource "aws_instance" "private" {
  depends_on                  = [aws_vpc.main]
  ami                         = data.aws_ami.amazon_linux.id
  instance_type               = var.instance_type
  subnet_id                   = aws_subnet.private.id
  associate_public_ip_address = false

  iam_instance_profile   = aws_iam_instance_profile.ssm_profile.name
  vpc_security_group_ids = [aws_security_group.private_instance_sg.id]

  tags = {
    Name        = "Gitea Private Instance"
    Service     = "gitea"
    Environment = "development"
  }

  lifecycle {
    ignore_changes = [ami]
  }
}


// Module to Create ec2 Storage Gateway 
module "ec2_sgw" {
  source                        = "aws-ia/storagegateway/aws//modules/ec2-sgw"
  vpc_id                        = aws_vpc.main.id
  subnet_id                     = aws_subnet.public1.id
  name                          = "my-storage-gateway_test"
  availability_zone             = data.aws_availability_zones.available.names[0]
  create_security_group         = true
  ingress_cidr_blocks           = "10.2.3.0/24"
  ingress_cidr_block_activation = var.ingress_cidr_block_activation // Important needs user Public IP of the PC
  ssh_key_name                  = local.ssh_key_name
}


# Module for NFS SHARE 
module "nfs_share" {
  source        = "aws-ia/storagegateway/aws//modules/s3-nfs-share"
  share_name    = "nfs_share_test"
  gateway_arn   = module.sgw.storage_gateway.arn
  bucket_arn    = module.s3_bucket.s3_bucket_arn
  role_arn      = aws_iam_role.sgw.arn
  client_list   = ["10.2.3.0/24"]
  log_group_arn = ""
}


