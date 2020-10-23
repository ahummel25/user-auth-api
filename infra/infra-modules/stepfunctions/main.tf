module "step_function" {
  source = "terraform-aws-modules/step-functions/aws"

  name = "emr-spark-processing-jobs"

  definition = local.definition_template

  attach_policy_jsons = true
  policy_jsons = [<<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "iam:PassRole"
            ],
            "Resource": [
                "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/EMR_DefaultRole",
                "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/EMR_EC2_DefaultRole"
            ],
            "Effect": "Allow"
        },
        {
            "Action": [
				"cloudwatch:*",
                "elasticmapreduce:RunJobFlow",
                "elasticmapreduce:TerminateJobFlows",
                "elasticmapreduce:DescribeCluster",
                "elasticmapreduce:AddJobFlowSteps",
                "elasticmapreduce:DescribeStep",
				"logs:*"
            ],
            "Resource": "*",
            "Effect": "Allow"
        }
    ]
}
EOF
  ]

  number_of_policy_jsons = 1
}
