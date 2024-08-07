resource "kubernetes_namespace" "namespace" {

  count      = 1
  metadata {
    name        = var.namespace_name
  }

}

