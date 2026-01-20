resource "aws_eks_cluster" "main" {
  name     = var.cluster_name
  role_arn = var.lab_role_arn
  version  = "1.30"

  vpc_config {
    subnet_ids              = concat(var.public_subnets, var.private_subnets)
    endpoint_public_access  = true
    endpoint_private_access = true
  }

  access_config {
    authentication_mode                         = "API_AND_CONFIG_MAP"
    bootstrap_cluster_creator_admin_permissions = true
  }
}

resource "aws_eks_node_group" "main" {
  cluster_name    = aws_eks_cluster.main.name
  node_group_name = "${var.cluster_name}-nodes"
  node_role_arn   = var.lab_role_arn
  subnet_ids      = var.private_subnets

  instance_types = ["t3.medium"]
  capacity_type  = "ON_DEMAND"

  scaling_config {
    desired_size = 2
    max_size     = 3
    min_size     = 1
  }

  update_config {
    max_unavailable = 1
  }
}