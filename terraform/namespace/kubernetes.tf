# Configure the AWS Provider
provider "aws" {
  region = "us-west-1"
}

data "aws_eks_cluster_auth" "cluster-auth" {
  name       = var.cluster_name
}

data "aws_eks_cluster" "cluster" {
  name       = var.cluster_name
}


provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  # token                  = data.aws_eks_cluster_auth.cluster-auth.token
  
  exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      args        = [
        "eks", "get-token", 
        "--cluster-name", data.aws_eks_cluster.cluster.name, 
        # "--role-arn", var.iam_admin_role_arn #add role-arn
      ]
      command     = "aws"
  }
}