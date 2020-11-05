locals {
  # Automatically load environment-level variables
  environment_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  tag_vars = read_terragrunt_config(find_in_parent_folders("tags.hcl"))

  # Extract out common variables for reuse
  env      = local.environment_vars.locals.environment
  all_tags = local.tag_vars.locals.all_tags
}

terraform {
  source = "../../../../../infra-modules/ssm"
}

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders()
}


inputs = {
  all_tags = local.all_tags
  env      = local.env
}
