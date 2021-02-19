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

resource "aws_iam_role_policy" "api_gateway_logs_role_policy" {
  name   = format("%s-%s", var.api_gateway_logs_role_name, "policy")
  role   = aws_iam_role.api_gateway_logs_role.id
  policy = data.aws_iam_policy.api-gateway-cloudwatch-logs-policy.policy
}

resource "aws_iam_role_policy" "lambda_vpc_role_policy" {
  count  = length(var.lambda_role_names)
  name   = format("%s-%s", var.lambda_role_names[count.index], "lambda-vpc-role-policy")
  role   = aws_iam_role.lambda_roles[count.index].id
  policy = data.aws_iam_policy.lambda-vpc-policy.policy
}

resource "aws_iam_role_policy" "lambda_role_policy" {
  count  = length(var.lambda_role_names)
  name   = format("%s-%s", var.lambda_role_names[count.index], "lambda-role-policy")
  role   = aws_iam_role.lambda_roles[count.index].id
  policy = data.aws_iam_policy_document.lambda-permissions-policy.json
}