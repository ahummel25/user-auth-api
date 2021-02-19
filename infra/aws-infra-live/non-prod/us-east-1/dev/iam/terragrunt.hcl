locals {
  # Automatically load environment-level variables
  environment_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  # Extract out common variables for reuse
  env = local.environment_vars.locals.environment

  api_gateway_logs_role_name = "user-auth-api-${local.env}-apigw-role"
  lambda_roles               = ["user-auth-api-${local.env}-auth-lambdaRole"]
}

terraform {
  source = "../../../../../infra-modules//iam"
}

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders()
}

inputs = {
  api_gateway_logs_role_name = local.api_gateway_logs_role_name
  env                        = local.env
  lambda_role_names          = local.lambda_roles
}