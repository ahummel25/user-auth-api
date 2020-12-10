locals {
  # Automatically load environment-level variables
  environment_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  tag_vars = read_terragrunt_config(find_in_parent_folders("tags.hcl"))

  # Extract out common variables for reuse
  env         = local.environment_vars.locals.environment
  common_tags = local.tag_vars.locals.common_tags

  param_names = {
    MongoDB_URI = "/user-auth-api/${local.env}/MONGODB_URI"
  }
}

terraform {
  source = "../../../../../infra-modules/ssm"
}

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders()
}


inputs = {
  common_tags = local.common_tags
  env         = local.env
  mongodb_uri = get_env("MONGODB_URI")
  param_names = local.param_names
}
