resource "aws_iam_policy" "session_manager_policy" {
  for_each    = var.ssm_kms_key_arn != null ? toset(["with_kms_encryption"]) : toset([])
  name        = "SessionManagerPolicy"
  path        = "/Gitea/"
  description = "Permissions for Session Manager logging"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "kms:Decrypt",
        ],
        Effect   = "Allow",
        Resource = var.ssm_kms_key_arn
      },
      {
        Action = [
          "kms:GenerateDataKey",
        ],
        Effect   = "Allow",
        Resource = "*"
      }
    ]
  })
}


# IAM Role for SSM Access 
resource "aws_iam_role" "ssm_role" {
  name = "SSM-Role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "ec2.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}


#Attach the SSM policy to the Role
resource "aws_iam_role_policy_attachment" "session_manager_attach" {
  for_each   = var.ssm_kms_key_arn != null ? toset(["with_kms_encryption"]) : toset([])
  role       = aws_iam_role.ssm_role.name
  policy_arn = aws_iam_policy.session_manager_policy[each.key].arn
}

resource "aws_iam_role_policy_attachment" "ssm_attach" {
  role       = aws_iam_role.ssm_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"

}

resource "aws_iam_role_policy_attachment" "cloudwatch_role_policy" {
  for_each   = var.ssm_kms_key_arn != null ? toset(["with_kms_encryption"]) : toset([])
  role       = aws_iam_role.ssm_role.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"
}


# IAM Instance Profile for the EC2 Instance
resource "aws_iam_instance_profile" "ssm_profile" {
  name = "SSM-Instance-Profile"
  role = aws_iam_role.ssm_role.name
}

# IAM Role for Storage Gateway S3 Access
resource "aws_iam_role" "sgw" {
  name               = "SGW-S3-Role"
  assume_role_policy = data.aws_iam_policy_document.sgw.json
}

# S3 Policy for Storage Gateway
resource "aws_iam_policy" "sgw" {
  name   = "SGW-S3-Policy"
  policy = data.aws_iam_policy_document.bucket_sgw.json
}

# Attach the S3 Policy to SGW Role
resource "aws_iam_role_policy_attachment" "sgw_attach" {
  role       = aws_iam_role.sgw.name
  policy_arn = aws_iam_policy.sgw.arn
}

# Policy Document for SGW Role
data "aws_iam_policy_document" "sgw" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["storagegateway.amazonaws.com"]
    }
  }
}

# S3 Bucket Policy Document for SGW
data "aws_iam_policy_document" "bucket_sgw" {
  statement {
    sid       = "AllowStorageGatewayBucketTopLevelAccess"
    effect    = "Allow"
    resources = [module.s3_bucket.s3_bucket_arn]
    actions = [
      "s3:GetAccelerateConfiguration",
      "s3:GetBucketLocation",
      "s3:GetBucketVersioning",
      "s3:ListBucket",
      "s3:ListBucketVersions",
      "s3:ListBucketMultipartUploads"
    ]
  }

  statement {
    sid       = "AllowStorageGatewayBucketObjectLevelAccess"
    effect    = "Allow"
    resources = ["${module.s3_bucket.s3_bucket_arn}/*"]
    actions = [
      "s3:AbortMultipartUpload",
      "s3:DeleteObject",
      "s3:DeleteObjectVersion",
      "s3:GetObject",
      "s3:GetObjectAcl",
      "s3:GetObjectVersion",
      "s3:ListMultipartUploadParts",
      "s3:PutObject",
      "s3:PutObjectAcl"
    ]
  }
}




