##################################
# VPC and all related components #
##################################
# https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/1.0.0

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.name}-vpc"

  cidr = var.vpc_cidr

  azs             = var.azs
  private_subnets = cidrsubnets(var.private_subnet, 4, 4, 4, 4, 4, 4)
  public_subnets  = cidrsubnets(var.public_subnet, 4, 4, 4, 4, 4, 4)

  create_database_subnet_group = false

  enable_dns_hostnames = true
  enable_dns_support   = true

  enable_ipv6 = true

  enable_nat_gateway     = true
  single_nat_gateway     = true
  one_nat_gateway_per_az = false

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

  name        = "${var.name}-security-group"
  description = "Security group for user auth api"

  vpc_id = module.vpc.vpc_id

  #ingress_cidr_blocks = [module.vpc.vpc_cidr_block]
  ingress_cidr_blocks = ["0.0.0.0/0"]

  # Prefix list ids to use in all ingress rules in this module.
  # ingress_prefix_list_ids = ["pl-123456"]
  # Open for all CIDRs defined in ingress_cidr_blocks
  ingress_rules = ["https-443-tcp"]

  # Use computed value here (eg, `${module...}`). Plain string is not a real use-case for this argument.
  computed_ingress_rules           = ["ssh-tcp"]
  number_of_computed_ingress_rules = 1

  ingress_with_self = [{
    rule = "all-all"
  }]

  egress_cidr_blocks = ["0.0.0.0/0"]

  # Prefix list ids to use in all egress rules in this module.
  # egress_prefix_list_ids = ["pl-123456"]
  # Open for all CIDRs defined in egress_cidr_blocks
  egress_rules = ["https-443-tcp"]

  egress_with_self = [{
    rule = "all-all"
  }]

  computed_egress_rules           = ["ssh-tcp"]
  number_of_computed_egress_rules = 1

  tags = {
    Application = var.common_tags["Application"]
    Project     = var.common_tags["Project"]
    Env         = var.env
  }
}
