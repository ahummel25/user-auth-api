variable "env" {
  type = string
}


variable "api_gateway_logs_role_name" {
  description = "The role name for the API Gateway logs."
  type        = string
}
variable "assumed_member_role_name" {
  description = "The role name belonging to the assumed account"
  type        = string
  default     = "adminAssumeRole"
}

variable "lambda_role_names" {
  description = "The role names belonging individual lambdas."
  type        = list(string)
}