locals {
  account_name   = "non-prod-account"
  aws_account_id = get_env("AWS_ACCOUNT_ID")
  aws_profile    = "non-prod-account-profile"
}
