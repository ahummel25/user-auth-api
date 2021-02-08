##################################
# VPC and all related components #
##################################
# https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/1.0.0

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.name}-${var.env}-vpc"

  cidr = var.vpc_cidr

  azs             = var.azs
  private_subnets = cidrsubnets(var.private_subnet, 4, 4, 4, 4, 4, 4)
  public_subnets  = cidrsubnets(var.public_subnet, 4, 4, 4, 4, 4, 4)

  create_database_subnet_group = false

  enable_dns_hostnames = true
  enable_dns_support   = true

  enable_ipv6 = false

  enable_nat_gateway = true
  single_nat_gateway = true

  manage_default_security_group  = true
  default_security_group_name    = "${var.name}-${var.env}-security-group-default-not-used"
  default_security_group_ingress = [{}]
  default_security_group_egress  = [{}]

  # -- VPC Flow Logs (Cloudwatch log group and IAM role will be created) -- #
  enable_flow_log                                 = true
  create_flow_log_cloudwatch_log_group            = true
  create_flow_log_cloudwatch_iam_role             = true
  flow_log_cloudwatch_log_group_name_prefix       = "${var.name}-vpc-flow-logs"
  flow_log_cloudwatch_log_group_retention_in_days = 90
  flow_log_traffic_type                           = "ALL"

  # -- TAGS -- #
  tags = {
    Application = var.common_tags["Application"]
    Project     = var.common_tags["Project"]
    Env         = var.env
  }
}

######
# SG #
######
module "security_group" {
  source = "terraform-aws-modules/security-group/aws"

  name        = "${var.name}-${var.env}-security-group"
  description = "Security group for user auth api"

  vpc_id = module.vpc.vpc_id

  #ingress_cidr_blocks = [module.vpc.vpc_cidr_block]
  ingress_cidr_blocks = [module.vpc.vpc_cidr_block]

  # Prefix list ids to use in all ingress rules in this module.
  # Open for all CIDRs defined in ingress_cidr_blocks
  ingress_rules = ["https-443-tcp"]

  # Prefix list ids to use in all egress rules in this module.
  # Open for all CIDRs defined in egress_cidr_blocks
  egress_rules = ["all-all"]

  tags = {
    Application = var.common_tags["Application"]
    Project     = var.common_tags["Project"]
    Env         = var.env
  }
}
