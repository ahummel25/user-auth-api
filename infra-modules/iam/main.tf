resource "aws_iam_role" "assumed_member_role" {
  name               = var.assumed_member_role_name
  assume_role_policy = data.aws_iam_policy_document.trust-assume-role-policy.json
}

resource "aws_iam_role" "lambda_roles" {
  count              = length(var.lambda_role_names)
  name               = var.lambda_role_names[count.index]
  assume_role_policy = data.aws_iam_policy_document.trust-lambda-assume-role-policy.json
}

resource "aws_iam_role" "api_gateway_logs_role" {
  name               = var.api_gateway_logs_role_name
  assume_role_policy = data.aws_iam_policy_document.trust-api-gateway-assume-role-policy.json
}

resource "aws_iam_role_policy" "lambda_role_policy" {
  count  = length(var.lambda_role_names)
  name   = format("%s-%s", var.lambda_role_names[count.index], "role-policy")
  role   = aws_iam_role.lambda_roles[count.index].id
  policy = data.aws_iam_policy_document.lambda-permissions-policy.json
}