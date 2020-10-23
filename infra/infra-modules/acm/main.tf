resource "aws_acm_certificate" "acm_certificate" {
  domain_name       = "*.andrewhummel.dev"
  validation_method = "DNS"
  key_algorithm     = "RSA_2048"

  options {
    certificate_transparency_logging_preference = var.certificate_transparency_logging_preference ? "ENABLED" : "DISABLED"
  }

  dynamic "validation_option" {
    for_each = var.validation_option

    content {
      domain_name       = try(validation_option.value["domain_name"], validation_option.key)
      validation_domain = validation_option.value["validation_domain"]
    }
  }

  tags = var.common_tags

  lifecycle {
    create_before_destroy = true
  }
}
