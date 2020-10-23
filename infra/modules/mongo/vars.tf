variable "teams" {
  description = "An object that contains all the groups that should be created in the project"
  type        = map(any)
  default     = {}
}

variable "white_lists" {
  description = "An object that contains all the network white-lists that should be created in the project"
  type        = map(any)
  default = {
    "All access" : "0.0.0.0/0",
  }
}

variable "replication_factor" {
  description = "The Number of replica set members, possible values are 3/5/7"
  type        = number
  default     = null
}

variable "provider_backup" {
  description = "Indicating if the cluster uses Cloud Backup for backups"
  type        = bool
  default     = true
}

variable "pit_enabled" {
  description = "Indicating if the cluster uses Continuous Cloud Backup, if set to true - provider_backup must also be set to true"
  type        = bool
  default     = false
}

variable "disk_size_gb" {
  description = "Capacity,in gigabytes,of the host’s root volume"
  type        = number
  default     = null
}

variable "auto_scaling_disk_gb_enabled" {
  description = "Indicating if disk auto-scaling is enabled"
  type        = bool
  default     = true
}

variable "provider_disk_iops" {
  description = "The maximum IOPS the system can perform"
  type        = number
  default     = null
}
