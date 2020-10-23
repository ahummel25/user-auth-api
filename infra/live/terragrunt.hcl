locals {
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl"))
  region_vars  = read_terragrunt_config(find_in_parent_folders("region.hcl"))
  tag_vars     = read_terragrunt_config(find_in_parent_folders("tags.hcl"))
  # Extract the variables we need for easy access
  account_name = local.account_vars.locals.account_name
  account_id   = local.account_vars.locals.aws_account_id
  aws_region   = local.region_vars.locals.aws_region
  env          = get_env("TF_WORKSPACE")
  backend_state_tags = {
    Description = "State management resource for all infrastructure as code"
  }
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
  provider "mongodbatlas" {
	  public_key = "${get_env("MONGO_PUBLIC_KEY")}"
	  private_key = "${get_env("MONGO_PRIVATE_KEY")}"
	}
  terraform {
    required_providers {
      mongodbatlas = {
        source  = "mongodb/mongodbatlas"
        version = "1.9.0"
      }
    }
  }
	EOF
}

terraform {
  before_hook "tflint" {
    commands = ["apply", "plan"]
    execute  = ["tflint"]
  }

  before_hook "init_reconfigure" {
    commands = ["apply", "destroy", "plan"]
    execute  = ["terraform", "init", "-reconfigure"]
  }
}

remote_state {
  backend = "s3"
  config = {
    encrypt             = true
    bucket              = "${local.env}-personal-terraform-state"
    key                 = "${path_relative_to_include()}/terraform.tfstate"
    region              = local.aws_region
    s3_bucket_tags      = local.backend_state_tags
    dynamodb_table      = "${local.env}-personal-terraform-state"
    dynamodb_table_tags = local.backend_state_tags
  }
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}

inputs = merge(
  local.account_vars.locals,
  local.region_vars.locals,
  local.tag_vars.locals,
)
