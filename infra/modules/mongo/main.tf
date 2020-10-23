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

  is_collect_database_specifics_statistics_enabled = true
  is_data_explorer_enabled                         = true
  is_performance_advisor_enabled                   = true
  is_realtime_performance_panel_enabled            = true
  is_schema_advisor_enabled                        = true
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
  mongo_db_major_version       = "6.0"
  cluster_type                 = "REPLICASET"
  num_shards                   = 1
  replication_factor           = var.replication_factor
  pit_enabled                  = false
  disk_size_gb                 = var.disk_size_gb
  auto_scaling_disk_gb_enabled = var.auto_scaling_disk_gb_enabled
  provider_disk_iops           = var.provider_disk_iops
}

# ---------------------------------------------------------------------------------------------------------------------
# Serverless Instance - Just for testing, not used in the API
# ---------------------------------------------------------------------------------------------------------------------
resource "mongodbatlas_serverless_instance" "serverless" {
  count      = local.env == "dev" ? 1 : 0
  project_id = mongodbatlas_project.project.id
  name       = "${title(local.env)}Serverless"

  provider_settings_backing_provider_name = "AWS"
  provider_settings_provider_name         = "SERVERLESS"
  provider_settings_region_name           = "US_EAST_1"
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

  dynamic "roles" {
    for_each = each.value.roles
    content {
      role_name     = roles.value.rolename
      database_name = roles.value.database
    }
  }

  labels {
    key   = each.value.labelkey
    value = each.value.labelvalue
  }

  dynamic "scopes" {
    for_each = each.value.scopes

    content {
      name = scopes.value.name
      type = scopes.value.type
    }
  }
}
