data "aws_iam_policy_document" "trust-assume-role-policy" {
  statement {

    actions = ["sts:AssumeRole", "sts:TagSession"]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_organizations_organization.org.master_account_id}:root"]
    }
  }
}

resource "aws_iam_role" "assumed_member_role" {
  name               = var.role_name
  assume_role_policy = data.aws_iam_policy_document.trust-assume-role-policy.json
}