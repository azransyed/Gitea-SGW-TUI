# To see Gitea EC2 instance ID or so called Hostname for ansibel 
output "Gitea_instance_id" {
  value = aws_instance.private.id
}

# For the region output
output "region" {
  value = var.region

}

# For Private IP of the storage gateway ec2 instance
output "storage_gateway_private_ip" {
  value       = module.ec2_sgw.private_ip
  description = "Private IP of the Storage Gateway EC2 Instance"
}

# STorage gateway id
output "storage_gateway_id" {
  value       = module.sgw.storage_gateway.gateway_id
  description = "Storage Gateway ID"
  sensitive   = true
}

# Name of the S3 bucket
output "s3_bucket_id" {
  value       = module.s3_bucket.s3_bucket_id
  description = "The name of the bucket."
}

output "s3_bucket_arn" {
  value       = module.s3_bucket.s3_bucket_arn
  description = "The ARN of the bucket. Will be of format arn:aws:s3:::bucketname."
}

output "nfs_share_arn" {
  value       = module.nfs_share.nfs_share_arn
  description = "ARN of the created NFS share"
}

# The mount path for nfs mount
output "nfs_share_path" {
  value       = module.nfs_share.nfs_share_path
  description = "NFS share mountpoint path"
}

