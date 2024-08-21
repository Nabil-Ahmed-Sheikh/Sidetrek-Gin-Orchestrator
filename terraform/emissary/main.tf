resource "kubernetes_manifest" "emissary_host" {
  manifest = {
    apiVersion = "getambassador.io/v2"
    kind       = "Host"
    metadata = {
      name      = "${var.host_name}"
      namespace = "emissary"
    }
    spec = {
      hostname = "${var.host_name}"
      tlsSecret = {
        name = "tls-subdomain-wildcard"
      }
    }
  }
}

resource "kubernetes_manifest" "emissary_mapping" {
  manifest = {
    apiVersion = "getambassador.io/v3alpha1"
    kind       = "Mapping"
    metadata = {
      name      = "${var.host_name}-mapping-name"
      namespace = "emissary"
    }
    spec = {
      prefix    = "/dagster"
      rewrite   = ""
      service   = "${var.svc_name}.${var.namespace_name}"
      host      = "${var.host_name}"
      tls       = true
      cors = {
        origins = ["*"]
        methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"]
        headers = ["Authorization", "Content-Type", "x-csrf-token"]
        credentials = true
      }
    }
  }
}