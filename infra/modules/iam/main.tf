resource "aws_iam_role" "lambda_roles" {
  count              = length(var.lambda_role_names)
  name               = var.lambda_role_names[count.index]
  assume_role_policy = data.aws_iam_policy_document.trust-lambda-assume-role-policy.json
  tags               = var.common_tags
}

resource "aws_iam_role" "api_gateway_logs_role" {
  name               = var.api_gateway_logs_role_name
  assume_role_policy = data.aws_iam_policy_document.trust-api-gateway-assume-role-policy.json
  tags               = var.common_tags
}

resource "aws_iam_role_policy" "api_gateway_logs_role_policy" {
  name   = data.aws_iam_policy.api-gateway-cloudwatch-logs-policy.name
  role   = aws_iam_role.api_gateway_logs_role.id
  policy = data.aws_iam_policy.api-gateway-cloudwatch-logs-policy.policy
}

resource "aws_iam_role_policy" "lambda_vpc_role_policy" {
  count  = length(var.lambda_role_names)
  name   = data.aws_iam_policy.lambda-vpc-policy.name
  role   = aws_iam_role.lambda_roles[count.index].id
  policy = data.aws_iam_policy.lambda-vpc-policy.policy
}

resource "aws_iam_role_policy" "lambda_role_policy" {
  count  = length(var.lambda_role_names)
  name   = format("%s-%s", var.lambda_role_names[count.index], "policy")
  role   = aws_iam_role.lambda_roles[count.index].id
  policy = data.aws_iam_policy_document.lambda-permissions-policy.json
}
