resource "aws_secretsmanager_secret" "crnunit_secrets" {
  name        = "api/${local.env}"
  description = "Configuration info necessary for the API"
  tags        = var.common_tags
}

resource "aws_secretsmanager_secret_version" "crnunit_secrets_version" {
  secret_id     = aws_secretsmanager_secret.crnunit_secrets.id
  secret_string = file("${local.env}_api_secrets.json")
}
