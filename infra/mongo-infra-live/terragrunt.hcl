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
