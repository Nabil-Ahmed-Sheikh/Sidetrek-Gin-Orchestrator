provider "aws" {
  profile = "default"
  region = "us-west-1"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}


data "aws_eks_cluster" "eks_cluster" {
  name = "sidetrek-app-cluster"
}

data "aws_eks_cluster_auth" "eks_cluster" {
  name = "sidetrek-app-cluster"
}



# Convert the subnet_ids set to a list
locals {
  subnet_ids_list = tolist(data.aws_eks_cluster.eks_cluster.vpc_config[0].subnet_ids)
}

resource "aws_security_group" "eks_node_sg" {
  name        = "eks-node-sg"
  description = "Allow communication between nodes and EKS cluster"
  vpc_id      = data.aws_eks_cluster.eks_cluster.vpc_config[0].vpc_id

  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}


resource "aws_instance" "app_server" {
  ami = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"

  # Use the first subnet from the list
  subnet_id = local.subnet_ids_list[0]

  # Optional: Attach a security group that allows communication with the EKS cluster
  security_groups = [aws_security_group.eks_node_sg.id]


  tags = {
    Name = "NabilTerraformInstance"
  }
}

# sidetrek-app-cluster