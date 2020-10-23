data "aws_s3_object" "db_users" {
  bucket = "${local.env}-personal-terraform-state"
  key    = "mongo/db_users.json"
}

data "mongodbatlas_roles_org_id" "org" {
}
