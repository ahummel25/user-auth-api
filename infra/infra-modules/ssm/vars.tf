variable "common_tags" {
  type = map(string)
}

variable "env" {
  type = string
}

variable "mongodb_uri" {
  type        = string
  description = "MongoDB connection string"
}

variable "param_names" {
  type = map(string)
}