output "role_name" {
  description = "The name of the IAM role created for the Step Function"
  value       = module.step_function.this_role_name
}
