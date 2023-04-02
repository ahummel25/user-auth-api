terraform {
  source = "${get_parent_terragrunt_dir()}/../modules//acm"
}

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders()
}
