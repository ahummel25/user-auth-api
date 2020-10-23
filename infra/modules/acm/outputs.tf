output "arn" {
  description = "ARN of the certificate"
  # value       = local.is_prod ? aws_acm_certificate.acm_certificate.*.arn : null
  value = aws_acm_certificate.acm_certificate.arn
}
