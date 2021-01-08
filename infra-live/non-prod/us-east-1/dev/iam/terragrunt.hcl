locals {
  # Automatically load environment-level variables
  environment_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  # Extract out common variables for reuse
  env = local.environment_vars.locals.environment

  lambda_roles = ["user-auth-api-${local.env}-auth-lambdaRole"]
}

terraform {
  source = "../../../../../infra-modules/iam"
}

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders()
}

inputs = {
  env               = local.env
  lambda_role_names = local.lambda_roles
}