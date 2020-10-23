data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "bucket_policy" {
  statement {
    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/EMR_DefaultRole",
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/EMR_EC2_DefaultRole",
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/emr-spark-processing-jobs"
      ]
    }

    actions = [
      "s3:*",
    ]

    resources = [
      "arn:aws:s3:::${local.bucket_name}",
      "arn:aws:s3:::${local.bucket_name}/*",
      "arn:aws:s3:::${local.bucket_name}/libs/*",
      "arn:aws:s3:::${local.bucket_name}/libs/spark/*",
      "arn:aws:s3:::${local.bucket_name}/resources/*",
      "arn:aws:s3:::${local.bucket_name}/resources/spark/*",
    ]
  }
}
