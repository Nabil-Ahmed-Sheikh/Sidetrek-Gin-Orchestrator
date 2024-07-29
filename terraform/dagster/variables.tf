variable "cluster_name" {
    type = string
}

variable "additional_set" {
  type        = list(any)
  description = "Additional sets to Helm"
  default     = []
}