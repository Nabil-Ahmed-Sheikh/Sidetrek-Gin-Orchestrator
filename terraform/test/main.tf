variable "repository" {
    type = string
}

locals {
  rendered_template = templatefile("${path.module}/opt/values.yaml", {
    repository = var.repository
  })
}

output "rendered_values_yaml" {
  value = local.rendered_template
}