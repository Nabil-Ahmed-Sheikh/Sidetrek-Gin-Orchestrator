provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster-auth.token
  }
}


resource "helm_release" "cluster_dagster" {
  name       = "c-dag"
  repository = "https://dagster-io.github.io/helm"
  chart      = "dagster"
  version = "1.8.3"
  namespace  = "dagster"
  create_namespace = false
  values = [
    # "${file("./opt/values.yaml")}"
    "${templatefile("${path.module}/opt/values.yaml", {
      repository = var.repository
    })}" 
  ]

  
  dynamic "set" {
    for_each = var.additional_set
    content {
      name  = set.value.name
      value = set.value.value
      type  = lookup(set.value, "type", null)
    }
  }

}