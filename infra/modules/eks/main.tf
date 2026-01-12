data "aws_caller_identity" "current" {}

# Conditional creation: create role only when create_cluster_role is true
resource "aws_iam_role" "eks_cluster" {
  count = var.create_cluster_role ? 1 : 0
  name  = "${var.cluster_name}-cluster-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "eks.amazonaws.com" }
    }]
  })
}

# Compute role name from ARN when user provides existing_cluster_role_arn
locals {
  existing_role_name_from_arn = var.existing_cluster_role_arn != "" ? element(split("/", var.existing_cluster_role_arn), length(split("/", var.existing_cluster_role_arn)) - 1) : ""
  cluster_role_arn  = var.existing_cluster_role_arn != "" ? var.existing_cluster_role_arn : (var.create_cluster_role ? aws_iam_role.eks_cluster[0].arn : "")
  cluster_role_name = var.existing_cluster_role_arn != "" ? local.existing_role_name_from_arn : (var.create_cluster_role ? aws_iam_role.eks_cluster[0].name : "")
}

# Attach cluster policy only when attach_cluster_policies is true. Role name uses local.cluster_role_name (from ARN or created role)
resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  for_each   = var.attach_cluster_policies ? { "cluster" = local.cluster_role_name } : {}
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = each.value
}

resource "aws_iam_role" "eks_node_group" {
  name = "${var.cluster_name}-node-group-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}
resource "aws_iam_role_policy_attachment" "eks_worker_node_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node_group.name
}
resource "aws_iam_role_policy_attachment" "eks_cni_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_node_group.name
}
resource "aws_iam_role_policy_attachment" "eks_container_registry_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.eks_node_group.name
}

resource "aws_eks_cluster" "main" {
  name     = var.cluster_name
  role_arn = local.cluster_role_arn
  version  = var.kubernetes_version

  vpc_config {
    subnet_ids              = concat(var.private_subnet_ids, var.public_subnet_ids)
    endpoint_private_access = true
    endpoint_public_access  = true
    public_access_cidrs     = ["0.0.0.0/0"]
  }

  access_config {
    authentication_mode                         = "API_AND_CONFIG_MAP"
    bootstrap_cluster_creator_admin_permissions = true
  }

  # Depend on the attachment instance (static reference) so attachment is created before cluster
  depends_on = [aws_iam_role_policy_attachment.eks_cluster_policy["cluster"]]
}

# Use lookup(...) to access the 'issuer' attribute to avoid static analysis unresolved-reference errors
locals {
  eks_oidc_issuer = lookup(try(aws_eks_cluster.main.identity[0].oidc[0], {}), "issuer", "")
}

data "tls_certificate" "eks" {
  url = local.eks_oidc_issuer
}

resource "aws_iam_openid_connect_provider" "eks" {
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = [data.tls_certificate.eks.certificates[0].sha1_fingerprint]
  url             = local.eks_oidc_issuer
}

locals {
  oidc_issuer_url = replace(local.eks_oidc_issuer, "https://", "")
}

resource "aws_iam_role" "vpc_cni_role" {
  name = "${var.cluster_name}-vpc-cni-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRoleWithWebIdentity"
      Effect    = "Allow"
      Principal = { Federated = aws_iam_openid_connect_provider.eks.arn }
      Condition = {
        StringEquals = {
          "${local.oidc_issuer_url}:sub" = "system:serviceaccount:kube-system:aws-node",
          "${local.oidc_issuer_url}:aud" = "sts.amazonaws.com"
        }
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "vpc_cni_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.vpc_cni_role.name
}

resource "aws_eks_node_group" "main" {
  cluster_name    = aws_eks_cluster.main.name
  node_group_name = "${var.cluster_name}-nodes"
  node_role_arn   = aws_iam_role.eks_node_group.arn
  subnet_ids      = var.private_subnet_ids

  instance_types = ["t4g.medium"]
  ami_type       = "AL2023_ARM_64_STANDARD"
  capacity_type  = "SPOT"
  disk_size      = 20

  scaling_config {
    desired_size = 2
    max_size     = 4
    min_size     = 1
  }

  tags = {
    "k8s.io/cluster-autoscaler/${var.cluster_name}" = "owned"
    "k8s.io/cluster-autoscaler/enabled"             = "true"
  }

  update_config {
    max_unavailable = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks_worker_node_policy,
    aws_iam_role_policy_attachment.eks_cni_policy,
    aws_iam_role_policy_attachment.eks_container_registry_policy,
  ]
}

resource "aws_eks_addon" "vpc_cni" {
  cluster_name                = aws_eks_cluster.main.name
  addon_name                  = "vpc-cni"
  resolve_conflicts_on_create = "OVERWRITE"
  resolve_conflicts_on_update = "OVERWRITE"
  service_account_role_arn    = aws_iam_role.vpc_cni_role.arn
  depends_on                  = [aws_eks_node_group.main]
}

resource "aws_eks_addon" "coredns" {
  cluster_name                = aws_eks_cluster.main.name
  addon_name                  = "coredns"
  resolve_conflicts_on_create = "OVERWRITE"
  depends_on                  = [aws_eks_addon.vpc_cni, aws_eks_node_group.main]
}

resource "aws_eks_addon" "kube_proxy" {
  cluster_name                = aws_eks_cluster.main.name
  addon_name                  = "kube-proxy"
  resolve_conflicts_on_create = "OVERWRITE"
}

resource "aws_eks_access_entry" "admin" {
  cluster_name  = aws_eks_cluster.main.name
  principal_arn = data.aws_caller_identity.current.arn
  type          = "STANDARD"
}

resource "aws_eks_access_policy_association" "admin" {
  cluster_name  = aws_eks_cluster.main.name
  policy_arn    = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
  principal_arn = aws_eks_access_entry.admin.principal_arn
  access_scope {
    type = "cluster"
  }
}