locals {
  account_name   = "dev-account"
  aws_account_id = get_env("AWS_DEV_ACCOUNT_ID")
  aws_profile    = "dev-account-profile"
}
