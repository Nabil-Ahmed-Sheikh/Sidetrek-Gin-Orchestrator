terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.22"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 3.0"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 2.0"
    }
  }
}

provider "aws" {
  region = "us-west-1"
}

provider "docker" {
    registry_auth {
        address  = var.ecr_address
        username = var.ecr_user
        password = var.ecr_password
    }
}