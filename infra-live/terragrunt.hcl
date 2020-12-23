locals {
  # Automatically load account-level variables
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl"))

  # Automatically load region-level variables
  region_vars = read_terragrunt_config(find_in_parent_folders("region.hcl"))

  # Automatically load environment-level variables
  environment_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  tag_vars = read_terragrunt_config(find_in_parent_folders("tags.hcl"))

  # Extract the variables we need for easy access
  account_name = local.account_vars.locals.account_name
  account_id   = local.account_vars.locals.aws_account_id
  aws_region   = local.region_vars.locals.aws_region
  environment  = local.environment_vars.locals.environment
  common_tags  = local.tag_vars.locals.common_tags
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
	provider "aws" {
	region = "${local.aws_region}"
	  # Only these AWS Account IDs may be operated on by this template
	  allowed_account_ids = ["${local.account_id}"]
	}
	EOF
}

remote_state {
  backend = "s3"
  config = {
    encrypt                = true
    bucket                 = "user-auth-api-terragrunt-state-${local.aws_region}-${local.environment}"
    key                    = "${path_relative_to_include()}/terraform.tfstate"
    region                 = local.aws_region
    dynamodb_table         = "terraform-locks"
  }
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}

inputs = merge(
  local.account_vars.locals,
  local.region_vars.locals,
  local.environment_vars.locals,
  local.common_tags,
)
