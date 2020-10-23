# ----------------------------------------------------------------------------------------------------------------------
# REQUIRE A SPECIFIC TERRAFORM VERSION OR HIGHER
# This module has been updated with 0.14 syntax, which means it is no longer compatible with any versions below 0.14.
# ----------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 1.3.7"
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "1.8.1"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE AN ATLAS PROJECT THAT THE CLUSTER WILL RUN INSIDE
# ---------------------------------------------------------------------------------------------------------------------

resource "mongodbatlas_project" "project" {
  name   = var.project_name
  org_id = var.org_id

  #Associate teams and privileges if passed, if not - run with an empty object
  dynamic "teams" {
    for_each = var.teams

    content {
      team_id    = mongodbatlas_teams.team[teams.key].team_id
      role_names = [teams.value.role]
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE NETWORK WHITE-LISTS FOR ACCESSING THE PROJECT
# ---------------------------------------------------------------------------------------------------------------------

#Optionall, if no variable is passed, the loop will run on an empty object.
resource "mongodbatlas_project_ip_access_list" "whitelists" {
  for_each = var.white_lists

  project_id = mongodbatlas_project.project.id
  comment    = each.key
  cidr_block = each.value
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE MONGODB ATLAS CLUSTER IN THE PROJECT
# ---------------------------------------------------------------------------------------------------------------------

resource "mongodbatlas_cluster" "cluster" {
  project_id                   = mongodbatlas_project.project.id
  provider_name                = "TENANT"
  backing_provider_name        = "AWS"
  provider_region_name         = var.region
  name                         = var.cluster_name
  provider_instance_size_name  = var.instance_type
  mongo_db_major_version       = var.mongodb_major_ver
  cluster_type                 = var.cluster_type
  num_shards                   = var.num_shards
  replication_factor           = var.replication_factor
  pit_enabled                  = var.pit_enabled
  disk_size_gb                 = var.disk_size_gb
  auto_scaling_disk_gb_enabled = var.auto_scaling_disk_gb_enabled
  provider_disk_iops           = var.provider_disk_iops
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE AWS PEER REQUESTS TO AWS VPC
# ---------------------------------------------------------------------------------------------------------------------
resource "mongodbatlas_network_peering" "mongo_peer" {
  for_each = var.vpc_peer

  accepter_region_name   = each.value.region
  project_id             = mongodbatlas_project.project.id
  container_id           = mongodbatlas_cluster.cluster.container_id
  provider_name          = "AWS"
  route_table_cidr_block = each.value.route_table_cidr_block
  vpc_id                 = each.value.vpc_id
  aws_account_id         = each.value.aws_account_id
}

# ---------------------------------------------------------------------------------------------------------------------
# ACCEPT THE PEER REQUESTS ON THE AWS SIDE
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_vpc_peering_connection_accepter" "peer" {
  for_each = var.vpc_peer

  vpc_peering_connection_id = mongodbatlas_network_peering.mongo_peer[each.key].connection_id
  auto_accept               = true
}
