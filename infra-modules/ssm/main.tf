resource "aws_ssm_parameter" "mongodb_uri" {
  name        = "/user-auth-api/dev/MONGODB_URI"
  description = "MongoDB connection string"
  type        = "SecureString"
  value       = var.mongodb_uri

  tags = {
    Application = var.common_tags["Application"]
    Project     = var.common_tags["Project"]
    Env         = var.env
  }
}