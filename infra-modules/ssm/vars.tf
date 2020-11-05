variable "mongodb_uri" {
  type        = string
  description = "MongoDB connection string"
}

variable "common_tags" {
  type = map(string)
}

variable "env" {
  type = string
}