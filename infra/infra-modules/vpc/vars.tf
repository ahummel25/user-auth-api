variable "env" {
  type = string
}

variable "name" {
  description = "The name prefix"
  default     = "user-auth"
}

variable "common_tags" {
  type = map(string)
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

variable "private_subnet" {
  type        = string
  description = "private subnet cidr"
  default     = "172.31.1.0/20"
}

variable "public_subnet" {
  type        = string
  description = "public subnet cidr"
  default     = "172.31.101.0/20"
}

variable "vpc_cidr" {
  type        = string
  description = "vpc cidr block"
  default     = "172.31.0.0/16"
}
