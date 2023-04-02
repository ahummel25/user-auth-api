variable "certificate_transparency_logging_preference" {
  description = "Specifies whether certificate details should be added to a certificate transparency log"
  type        = bool
  default     = true
}

variable "validation_option" {
  description = "The domain name that you want ACM to use to send you validation emails. This domain name is the suffix of the email addresses that you want ACM to use."
  type        = any
  default     = {}
}

variable "common_tags" {
  description = "A mapping of tags to assign to the resource"
  type        = map(string)
  default     = {}
}
