resource "aws_iam_role" "assumed_member_role" {
  name               = var.role_name
  assume_role_policy = data.aws_iam_policy_document.trust-assume-role-policy.json
}