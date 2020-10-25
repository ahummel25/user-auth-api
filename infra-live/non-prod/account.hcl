locals {
  account_name   = "dev-account"
  aws_account_id = get_env("AWS_ACCOUNT_ID_DEV")
  aws_profile    = "dev-account-profile"
}
