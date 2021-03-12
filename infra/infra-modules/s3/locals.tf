locals {
  bucket_name = "stepfunctions-emrproject-${data.aws_caller_identity.current.account_id}"
}
