locals {
  env = get_env("TF_WORKSPACE")
}

terraform {
  source = "${get_parent_terragrunt_dir()}/../modules//mongo"
}

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders()
}