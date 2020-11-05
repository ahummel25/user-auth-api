resource "aws_ssm_parameter" "mongodb_uri" {
  name        = "/user-auth-api/dev/MONGODB_URI"
  description = "MongoDB connection string"
  type        = "SecureString"
  value       = var.mongodb_uri

  tags = {
    Application = var.all_tags["Application"]
    Project     = var.all_tags["Project"]
    Env         = var.env
  }
}