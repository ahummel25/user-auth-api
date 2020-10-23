output "arn" {
  description = "ARN of the certificate"
  value       = aws_acm_certificate.acm_certificate.arn
}
