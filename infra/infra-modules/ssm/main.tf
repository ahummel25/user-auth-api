resource "aws_ssm_parameter" "mongodb_uri" {
  name        = var.param_names["MongoDB_URI"]
  description = "MongoDB connection string"
  type        = "SecureString"
  value       = var.mongodb_uri

  tags = {
    Application = var.common_tags["Application"]
    Project     = var.common_tags["Project"]
    Env         = var.env
  }
}