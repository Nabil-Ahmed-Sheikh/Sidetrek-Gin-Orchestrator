provider "aws" {
  region = "us-west-1"
}

data "aws_eks_cluster_auth" "cluster-auth" {
  name       = "sidetrek-app-cluster"
}

data "aws_eks_cluster" "cluster" {
  name       = "sidetrek-app-cluster"
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster-auth.token
}


resource "kubernetes_deployment" "my_app" {
  metadata {
    name = "my-app"
    namespace = "dagster"
    labels = {
      app = "my-app"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "my-app"
      }
    }
    template {
      metadata {
        labels = {
          app = "my-app"
        }
      }
      spec {
        container {
          image = "447632895027.dkr.ecr.us-west-1.amazonaws.com/dagster-test-project:20240814121246"
          name  = "my-app"
          port {
            container_port = 3000
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "my_app" {
  metadata {
    name = "my-app"
    namespace = "dagster"
  }
  spec {
    selector = {
      app = "my-app"
    }
    port {
      port        = 3000
      target_port = 3000
    }
    # type = "LoadBalancer"
  }
}