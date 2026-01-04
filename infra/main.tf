provider "aws" {
  region  = var.aws_region
  profile = "terraform-admin"
}

data "aws_availability_zones" "available" {
  state = "available"
}
data "aws_caller_identity" "current" {}

# ==========================================
# 1. VPC NETWORK (Mantido igual)
# ==========================================
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Name                                        = "${var.cluster_name}-vpc"
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  tags   = { Name = "${var.cluster_name}-igw" }
}

resource "aws_subnet" "public" {
  count                   = 2
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  cidr_block              = "10.0.${count.index + 1}.0/24"
  vpc_id                  = aws_vpc.main.id
  map_public_ip_on_launch = true
  tags = {
    Name                                        = "${var.cluster_name}-public-${count.index + 1}"
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                    = "1"
  }
}

resource "aws_subnet" "private" {
  count             = 2
  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = "10.0.${count.index + 10}.0/24"
  vpc_id            = aws_vpc.main.id
  tags = {
    Name                                        = "${var.cluster_name}-private-${count.index + 1}"
    "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    "kubernetes.io/role/internal-elb"           = "1"
  }
}

resource "aws_eip" "nat" {
  count  = 2
  domain = "vpc"
  tags   = { Name = "${var.cluster_name}-nat-${count.index + 1}" }
}

resource "aws_nat_gateway" "main" {
  count         = 2
  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id
  tags          = { Name = "${var.cluster_name}-nat-${count.index + 1}" }
  depends_on    = [aws_internet_gateway.main]
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
  tags = { Name = "${var.cluster_name}-public" }
}

resource "aws_route_table" "private" {
  count  = 2
  vpc_id = aws_vpc.main.id
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main[count.index].id
  }
  tags = { Name = "${var.cluster_name}-private-${count.index + 1}" }
}

resource "aws_route_table_association" "public" {
  count          = 2
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count          = 2
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

# ==========================================
# 2. IAM ROLES (Cluster & Node)
# ==========================================

resource "aws_iam_role" "eks_cluster" {
  name = "${var.cluster_name}-cluster-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "eks.amazonaws.com" }
    }]
  })
}
resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_cluster.name
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

# ==========================================
# 3. EKS CLUSTER
# ==========================================

resource "aws_eks_cluster" "main" {
  name     = var.cluster_name
  role_arn = aws_iam_role.eks_cluster.arn
  version  = var.kubernetes_version

  vpc_config {
    subnet_ids              = concat(aws_subnet.private[*].id, aws_subnet.public[*].id)
    endpoint_private_access = true
    endpoint_public_access  = true
    public_access_cidrs     = ["0.0.0.0/0"]
  }

  access_config {
    authentication_mode                         = "API_AND_CONFIG_MAP"
    bootstrap_cluster_creator_admin_permissions = true
  }

  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  depends_on = [aws_iam_role_policy_attachment.eks_cluster_policy]
}

# ==========================================
# 4. OIDC & IRSA CONFIG (AWS Recommended)
# ==========================================

# Setup OIDC Provider
data "tls_certificate" "eks" {
  url = aws_eks_cluster.main.identity[0].oidc[0].issuer
}
resource "aws_iam_openid_connect_provider" "eks" {
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = [data.tls_certificate.eks.certificates[0].sha1_fingerprint]
  url             = aws_eks_cluster.main.identity[0].oidc[0].issuer
}

# Helper local to clean up the URL for policy conditions
locals {
  oidc_issuer_url = replace(aws_eks_cluster.main.identity[0].oidc[0].issuer, "https://", "")
}

# ------------------------------------
# Role for VPC CNI (Networking)
# ------------------------------------
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

# ------------------------------------
# Role for EBS CSI (Storage)
# ------------------------------------
resource "aws_iam_role" "ebs_csi_role" {
  name = "${var.cluster_name}-ebs-csi-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRoleWithWebIdentity"
      Effect    = "Allow"
      Principal = { Federated = aws_iam_openid_connect_provider.eks.arn }
      Condition = {
        StringEquals = {
          "${local.oidc_issuer_url}:sub" = "system:serviceaccount:kube-system:ebs-csi-controller-sa",
          "${local.oidc_issuer_url}:aud" = "sts.amazonaws.com"
        }
      }
    }]
  })
}
resource "aws_iam_role_policy_attachment" "ebs_csi_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
  role       = aws_iam_role.ebs_csi_role.name
}

# ==========================================
# 5. NODE GROUP
# ==========================================

resource "aws_eks_node_group" "main" {
  cluster_name    = aws_eks_cluster.main.name
  node_group_name = "${var.cluster_name}-nodes"
  node_role_arn   = aws_iam_role.eks_node_group.arn
  subnet_ids      = aws_subnet.private[*].id

  instance_types = ["t4g.medium"]
  ami_type       = "AL2023_ARM_64_STANDARD"
  capacity_type  = "SPOT"
  disk_size      = 20

  scaling_config {
    desired_size = 2
    max_size     = 4
    min_size     = 1
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

# ==========================================
# 6. ADDONS (With Priority Order)
# ==========================================

# 1. VPC CNI - Most Critical (Install First)
resource "aws_eks_addon" "vpc_cni" {
  cluster_name                = aws_eks_cluster.main.name
  addon_name                  = "vpc-cni"
  addon_version               = var.vpc_cni_version
  resolve_conflicts_on_create = "OVERWRITE"
  resolve_conflicts_on_update = "OVERWRITE"
  service_account_role_arn    = aws_iam_role.vpc_cni_role.arn

  depends_on = [
    aws_eks_node_group.main
  ]
}

# 2. CoreDNS - Depends on CNI
resource "aws_eks_addon" "coredns" {
  cluster_name                = aws_eks_cluster.main.name
  addon_name                  = "coredns"
  addon_version               = var.coredns_version
  resolve_conflicts_on_create = "OVERWRITE"

  depends_on = [
    aws_eks_addon.vpc_cni,
    aws_eks_node_group.main
  ]
}

# 3. Kube Proxy - Depends on CNI
resource "aws_eks_addon" "kube_proxy" {
  cluster_name = aws_eks_cluster.main.name
  addon_name   = "kube-proxy"
  #addon_version               = var.kube_proxy_version
  resolve_conflicts_on_create = "OVERWRITE"

  depends_on = [
    aws_eks_addon.vpc_cni,
    aws_eks_node_group.main
  ]
}

# 4. EBS CSI - Storage (Depends on CNI)
resource "aws_eks_addon" "ebs_csi_driver" {
  cluster_name                = aws_eks_cluster.main.name
  addon_name                  = "aws-ebs-csi-driver"
  addon_version               = var.ebs_csi_version
  resolve_conflicts_on_create = "OVERWRITE"
  service_account_role_arn    = aws_iam_role.ebs_csi_role.arn

  depends_on = [
    aws_eks_addon.vpc_cni,
    aws_eks_node_group.main
  ]
}

# ==========================================
# 7. ACCESS ENTRY (Admin User)
# ==========================================
# resource "aws_eks_access_entry" "admin" {
#   cluster_name  = aws_eks_cluster.main.name
#   principal_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/terraform-admin"
#   type          = "STANDARD"
# }

# resource "aws_eks_access_policy_association" "admin" {
#   cluster_name  = aws_eks_cluster.main.name
#   policy_arn    = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
#   principal_arn = aws_eks_access_entry.admin.principal_arn
#   access_scope {
#     type = "cluster"
#   }
# }