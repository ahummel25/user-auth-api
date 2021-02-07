locals {
  account_name   = "non-prod-account"
  aws_account_id = get_env("AWS_NON_PROD_ACCOUNT_ID")
  aws_profile    = "non-prod-account-profile"
}
