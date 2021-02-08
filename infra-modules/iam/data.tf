data "aws_organizations_organization" "org" {}

data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

data "aws_iam_policy" "api-gateway-cloudwatch-logs-policy" {
  arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}

data "aws_iam_policy_document" "trust-assume-role-policy" {
  statement {

    actions = ["sts:AssumeRole", "sts:TagSession"]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_organizations_organization.org.master_account_id}:root"]
    }
  }
}

data "aws_iam_policy_document" "trust-api-gateway-assume-role-policy" {
  statement {

    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "trust-lambda-assume-role-policy" {
  statement {

    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "ec2-permissions-policy" {
  count = length(var.lambda_role_names)

  statement {
    effect = "Allow"

    actions = [
      "ec2:DescribeNetworkInterfaces",
      "ec2:CreateNetworkInterface",
      "ec2:DeleteNetworkInterface",
      "ec2:DescribeInstances",
      "ec2:AttachNetworkInterface"
    ]

    resources = [aws_iam_role.lambda_roles[count.index].arn]
  }
}

data "aws_iam_policy_document" "lambda-permissions-policy" {

  statement {
    effect = "Allow"

    actions = [
      "lambda:AddPermission",
      "lambda:InvokeFunction",
      "lambda:RemovePermission",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords"
    ]

    resources = ["arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:*"]
  }
}