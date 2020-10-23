resource "aws_secretsmanager_secret" "crnunit_secrets" {
  name        = "andrewhummel/api/${var.environment}"
  description = "Configuration info necessary for the API"
}

resource "aws_secretsmanager_secret_version" "crnunit_secrets_version" {
  secret_id     = aws_secretsmanager_secret.crnunit_secrets.id
  secret_string = file("api.${var.environment}.secrets.json")
}
