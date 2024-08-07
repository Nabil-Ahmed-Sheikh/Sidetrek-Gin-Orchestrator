provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster-auth.token
  }
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "dagster"
  }
}


resource "helm_release" "cluster_dagster" {
  name       = "cluster-dagster"
  repository = "https://dagster-io.github.io/helm"
  chart      = "dagster"
  namespace  = kubernetes_namespace.example.metadata[0].name
  create_namespace = false
  
  dynamic "set" {
    for_each = var.additional_set
    content {
      name  = set.value.name
      value = set.value.value
      type  = lookup(set.value, "type", null)
    }
  }
}