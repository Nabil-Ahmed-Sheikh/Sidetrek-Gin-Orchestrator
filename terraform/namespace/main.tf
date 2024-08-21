resource "kubernetes_namespace" "namespace" {
  count = var.namespace_count
  metadata {
    name = var.namespace_name
  }
}