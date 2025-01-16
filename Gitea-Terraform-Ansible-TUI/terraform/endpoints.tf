#S3 Bucket Creation
module "s3_bucket" {
  source                   = "terraform-aws-modules/s3-bucket/aws"
  version                  = ">=3.5.0"
  bucket                   = "chai-gitea-nahi" // Change to your desired name and change also in ansible playbook
  control_object_ownership = true
  object_ownership         = "BucketOwnerEnforced"
  block_public_acls        = true
  block_public_policy      = true
  ignore_public_acls       = true
  restrict_public_buckets  = true

  logging = {
    target_bucket = module.log_delivery_bucket.s3_bucket_id
    target_prefix = "log/"
  }

  versioning = {
    enabled = false
  }

}


module "log_delivery_bucket" {
  source                   = "terraform-aws-modules/s3-bucket/aws"
  version                  = ">=3.5.0"
  bucket                   = "logs-gitea-nahi"
  control_object_ownership = true
  object_ownership         = "BucketOwnerEnforced"
  block_public_acls        = true
  block_public_policy      = true
  ignore_public_acls       = true
  restrict_public_buckets  = true

  versioning = {
    enabled = false // for later if ammar wants we turn it true
  }
}

# S3 Gateway VPC Endpoint
resource "aws_vpc_endpoint" "s3" {
  vpc_id       = aws_vpc.main.id
  service_name = "com.amazonaws.${var.region}.s3"
}




# VPC Endpoints for SSM
resource "aws_vpc_endpoint" "ssm_endpoint" {
  vpc_id              = aws_vpc.main.id
  service_name        = "com.amazonaws.${var.region}.ssm"
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [aws_security_group.ssm_https.id]
  private_dns_enabled = true
  subnet_ids          = [aws_subnet.private.id]
}

resource "aws_vpc_endpoint" "ssm_messages_endpoint" {
  vpc_id              = aws_vpc.main.id
  service_name        = "com.amazonaws.${var.region}.ssmmessages"
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [aws_security_group.ssm_https.id]
  private_dns_enabled = true
  subnet_ids          = [aws_subnet.private.id]
}

resource "aws_vpc_endpoint" "ec2_endpoint" {
  vpc_id              = aws_vpc.main.id
  service_name        = "com.amazonaws.${var.region}.ec2"
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [aws_security_group.ssm_https.id]
  private_dns_enabled = true
  subnet_ids          = [aws_subnet.private.id]
}

