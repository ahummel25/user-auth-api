output "lambda_roles_ids" {
  description = "The name of the role."
  value       = aws_iam_role.lambda_roles.*.id
}

output "lambda_roles_names" {
  description = "The name of the role."
  value       = aws_iam_role.lambda_roles.*.name
}

output "lambda_roles_create_date" {
  description = "The creation date of the IAM role."
  value       = aws_iam_role.lambda_roles.*.create_date
}

output "lambda_roles_arn" {
  description = "The Amazon Resource Name (ARN) specifying the role."
  value       = aws_iam_role.lambda_roles.*.arn
}

output "api_gateway_roles_arn" {
  description = "The Amazon Resource Name (ARN) specifying the role."
  value       = aws_iam_role.api_gateway_logs_role.arn
}