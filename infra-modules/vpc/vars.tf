variable "name" {
  description = "The name prefix"
  default     = "user-auth"
}

variable "all_tags" {
  type        = map(string)
  description = "All other tags"
  default = {
    Application = "User Auth"
    Project     = "User Auth"
    Env         = "dev"
  }
}

variable "azs" {
  type        = list(string)
  description = "Availability zones to be used for subnets"
  default     = ["us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d", "us-east-1e", "us-east-1f"]
}

variable "sg_cidr_block" {
  type        = string
  description = "Ingress IP with CIDR for the bastion host security group."
  default     = "0.0.0.0/0"
}

variable "private_subnets" {
  type        = list(string)
  description = "private subnet cidrs"
  default     = ["172.31.1.0/24", "172.31.2.0/24", "172.31.3.0/24", "172.31.4.0/24", "172.31.5.0/24", "172.31.6.0/24"]
}

variable "public_subnets" {
  type        = list(string)
  description = "public subnet cidrs"
  default     = ["172.31.101.0/24", "172.31.102.0/24", "172.31.103.0/24", "172.31.104.0/24", "172.31.105.0/24", "172.31.106.0/24"]
}

variable "vpc_cidr" {
  type        = string
  description = "vpc cidr block"
  default     = "172.31.0.0/16"
}
