##################################
# VPC and all related components #
##################################
# https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/1.0.0

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.name}-vpc"

  cidr = var.vpc_cidr

  azs             = var.azs
  private_subnets = var.private_subnets
  public_subnets  = var.public_subnets

  create_database_subnet_group = false

  enable_dns_hostnames = true
  enable_dns_support   = true

  enable_ipv6 = true

  enable_nat_gateway = true
  single_nat_gateway = true

  # -- VPC Flow Logs (Cloudwatch log group and IAM role will be created) -- #
  enable_flow_log                                 = true
  create_flow_log_cloudwatch_log_group            = true
  create_flow_log_cloudwatch_iam_role             = true
  flow_log_cloudwatch_log_group_name_prefix       = "${var.name}-vpc-flow-logs"
  flow_log_cloudwatch_log_group_retention_in_days = 90
  flow_log_traffic_type                           = "ALL"

  # -- TAGS -- #
  tags = {
    Application = var.all_tags["Application"]
    Project     = var.all_tags["Project"]
    Env         = var.all_tags["Env"]
  }
}
