locals {
  env                        = get_env("TF_WORKSPACE")
  api_gateway_logs_role_name = "api-apigw-role"
  lambda_roles               = ["api-lambda-role"]
}

terraform {
  source = "${get_parent_terragrunt_dir()}/../modules//iam"
}

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders()
}

inputs = {
  api_gateway_logs_role_name = local.api_gateway_logs_role_name
  lambda_role_names          = local.lambda_roles
}
