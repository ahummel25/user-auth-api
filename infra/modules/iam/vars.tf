variable "api_gateway_logs_role_name" {
  description = "The role name for the API Gateway logs."
  type        = string
}

variable "common_tags" {
  description = "A mapping of tags to assign to the resource"
  type        = map(string)
  default     = {}
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
