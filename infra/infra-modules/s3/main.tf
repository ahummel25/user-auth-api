# resource "aws_kms_key" "objects" {
#   description             = "KMS key is used to encrypt bucket objects"
#   deletion_window_in_days = 7
# }

module "emr_logging_s3_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"

  bucket        = local.bucket_name
  acl           = "private"
  force_destroy = true

  #   attach_policy = true
  #   policy        = data.aws_iam_policy_document.bucket_policy.json

  tags = {
    Owner = "EMR Step Functions"
  }

  versioning = {
    enabled = true
  }

  #   server_side_encryption_configuration = {
  #     rule = {
  #       apply_server_side_encryption_by_default = {
  #         kms_master_key_id = aws_kms_key.objects.arn
  #         sse_algorithm     = "aws:kms"
  #       }
  #     }
  #   }

  #   object_lock_configuration = {
  #     object_lock_enabled = "Enabled"
  #     rule = {
  #       default_retention = {
  #         mode = "GOVERNANCE"
  #         days = 1
  #       }
  #     }
  #   }

  # S3 bucket-level Public Access Block configuration
  #   block_public_acls       = false
  #   block_public_policy     = false
  #   ignore_public_acls      = false
  #   restrict_public_buckets = false
}
