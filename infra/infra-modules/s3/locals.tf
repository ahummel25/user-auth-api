locals {
  bucket_name = "step-functions-emr-${data.aws_caller_identity.current.account_id}"
}
