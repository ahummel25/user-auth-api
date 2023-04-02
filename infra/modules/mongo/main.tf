# ---------------------------------------------------------------------------------------------------------------------
# CREATE AN ATLAS PROJECT THAT THE CLUSTER WILL RUN INSIDE
# ---------------------------------------------------------------------------------------------------------------------

resource "mongodbatlas_project" "project" {
  name   = "${title(local.env)} DB"
  org_id = data.mongodbatlas_roles_org_id.org.org_id
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
  provider_region_name         = "US_EAST_1"
  name                         = "${title(local.env)}Cluster"
  provider_instance_size_name  = "M0"
  mongo_db_major_version       = 5.0
  cluster_type                 = "REPLICASET"
  num_shards                   = 1
  replication_factor           = var.replication_factor
  pit_enabled                  = false
  disk_size_gb                 = var.disk_size_gb
  auto_scaling_disk_gb_enabled = var.auto_scaling_disk_gb_enabled
  provider_disk_iops           = var.provider_disk_iops
}

# ---------------------------------------------------------------------------------------------------------------------
# DB USERS
# ---------------------------------------------------------------------------------------------------------------------
resource "mongodbatlas_database_user" "db_user" {
  for_each = { for user in local.db_user_data.users : user.username => user }

  username           = each.value.username
  password           = each.value.password
  project_id         = mongodbatlas_project.project.id
  auth_database_name = "admin"

  roles {
    role_name     = each.value.rolename
    database_name = each.value.database
  }

  labels {
    key   = each.value.labelkey
    value = each.value.labelvalue
  }

  scopes {
    name = "${title(local.env)}Cluster"
    type = "CLUSTER"
  }
}
