#VPC and Subnet Creation
provider "aws" {
  region = var.region
}

// Aws Availability Zone Data Source 
data "aws_availability_zones" "available" {
  state = "available"

}

# VPC Creation
resource "aws_vpc" "main" {
  cidr_block           = "10.2.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name = "Gitea VPC"
  }
}
#Public Subnet 1 (for NAT GATWAY it will store here)
resource "aws_subnet" "public1" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.2.1.0/24"
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true
  tags = {
    Name = "Public Subnet1"
  }

}

# Public Subnet 2
resource "aws_subnet" "public2" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.2.2.0/24"
  availability_zone       = data.aws_availability_zones.available.names[1]
  map_public_ip_on_launch = true
  tags = {
    Name = "Public Subnet2"
  }

}


# Private Subnet
resource "aws_subnet" "private" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.2.3.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]
  tags = {
    Name = "Private Subnet"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "InternetGateway"
  }
}

# NAT Gateway for Private Subnet Internet Access
resource "aws_eip" "nat" {
}
resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public1.id
  tags = {
    Name = "NATGateway"
  }
}




#Route Table Public Subnet (Internet Access)
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
  tags = {
    Name = "PublicRouteTable"
  }
}

# Associate Public Route Table with Public Subent 1
resource "aws_route_table_association" "public_association1" {
  subnet_id      = aws_subnet.public1.id
  route_table_id = aws_route_table.public.id
}

# Associate Public Route Table with Public Subent 2
resource "aws_route_table_association" "public_association2" {
  subnet_id      = aws_subnet.public2.id
  route_table_id = aws_route_table.public.id
}


#Route Table for Private Subnet (NAT Gatewat Access)
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat.id
  }
  tags = {
    Name = "PrivateRouteTable"
  }
}
#Associate Private Route Table with Private Subnet
resource "aws_route_table_association" "private" {
  subnet_id      = aws_subnet.private.id
  route_table_id = aws_route_table.private.id
}

# Storage Gateway Module
module "sgw" {
  depends_on                         = [module.ec2_sgw]
  source                             = "aws-ia/storagegateway/aws//modules/aws-sgw"
  gateway_name                       = "my-storage-gateway"
  gateway_ip_address                 = module.ec2_sgw.public_ip
  join_smb_domain                    = false
  gateway_type                       = "FILE_S3"
  create_vpc_endpoint                = true
  create_vpc_endpoint_security_group = true
  vpc_id                             = aws_vpc.main.id
  vpc_endpoint_subnet_ids            = [aws_subnet.private.id]   // Talk
  gateway_private_ip_address         = module.ec2_sgw.private_ip // we talk
}


// For the Keypair for the storage gateway Ec2 please define the key name and path and the key in the variables.tf
locals {
  ssh_key_name = length(var.ssh_public_key_path) > 0 ? aws_key_pair.ec2_sgw_key_pair["ec2_sgw_key_pair"].key_name : null
}

resource "aws_key_pair" "ec2_sgw_key_pair" {

  for_each = length(var.ssh_public_key_path) > 0 ? toset(["ec2_sgw_key_pair"]) : toset([])

  key_name   = var.ssh_key_name
  public_key = file(var.ssh_public_key_path)
}











