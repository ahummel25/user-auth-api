locals {
  env          = terraform.workspace
  db_user_data = jsondecode(data.aws_s3_object.db_users.body)
}
