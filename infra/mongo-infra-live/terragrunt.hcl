locals {
  # Automatically load region-level variables
  region_vars = read_terragrunt_config(find_in_parent_folders("region.hcl"))

  # Automatically load environment-level variables
  environment_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  aws_region   = local.region_vars.locals.aws_region
  environment  = local.environment_vars.locals.environment
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
	provider "mongodbatlas" {
	  public_key = "${get_env("MONGO_PUBLIC_KEY")}"
	  private_key = "${get_env("MONGO_PRIVATE_KEY")}"
	}
	EOF
}

remote_state {
  backend = "s3"
  config = {
    encrypt        = true
    bucket         = "user-auth-api-mongodb-terragrunt-state-${local.aws_region}-${local.environment}"
    key            = "${path_relative_to_include()}/terraform.tfstate"
    region         = local.aws_region
    dynamodb_table = "terraform-locks-mongodb"
  }
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}
