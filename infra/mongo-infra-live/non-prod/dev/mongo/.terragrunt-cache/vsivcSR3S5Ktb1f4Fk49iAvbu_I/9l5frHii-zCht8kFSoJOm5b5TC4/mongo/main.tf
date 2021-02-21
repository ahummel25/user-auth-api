terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "0.8.2"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

locals {
  env = title(var.env)
}

module "atlas_cluster" {
  source = "./source"

  project_name = "User Auth ${local.env}"
  org_id       = var.org_id

  white_lists = {
    "All access" : "0.0.0.0/0",
  }

  region = "US_EAST_1"

  cluster_name = "UserAuthMongoCluster${local.env}"

  instance_type     = "M0"
  mongodb_major_ver = 4.4
  cluster_type      = "REPLICASET"
  num_shards        = 1
  provider_backup   = true
  pit_enabled       = false
}