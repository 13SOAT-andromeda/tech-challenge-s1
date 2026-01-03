provider "aws" {
  region = var.region
}

resource "aws_iam_role" "ebs_csi_driver" {
  name = "${var.cluster_name}-ebs-csi-driver"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "pods.eks.amazonaws.com"
        }
        Action = [
          "sts:AssumeRole",
          "sts:TagSession"
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ebs_csi_driver" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
  role       = aws_iam_role.ebs_csi_driver.name
}

resource "aws_eks_pod_identity_association" "ebs_csi" {
  cluster_name    = module.eks.cluster_name
  namespace       = "kube-system"
  service_account = "ebs-csi-controller-sa"
  role_arn        = aws_iam_role.ebs_csi_driver.arn

  depends_on = [module.eks]
}

resource "aws_eks_addon" "vpc_cni" {
  cluster_name = module.eks.cluster_name
  addon_name   = "vpc-cni"

  # Ensure we get the latest version suitable for the cluster
  // v1.20.4-eksbuild.2
  addon_version = "v1.19.2-eksbuild.1" # OR set to null/remove to auto-detect, but hardcoding is safer for stability

  # The "Nuclear" Overwrite flags
  resolve_conflicts_on_create = "OVERWRITE"
  resolve_conflicts_on_update = "OVERWRITE"

  service_account_role_arn = null

  # ---------------------------------------------------------
  # FIX: The "Depends On" you wanted
  # ---------------------------------------------------------
  # This forces Terraform to wait until the Control Plane AND Node Groups
  # are fully finished before even attempting to touch the CNI.
  depends_on = [
    module.eks.eks_managed_node_groups,
    module.eks
  ]
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 6.0"

  name = "${var.cluster_name}-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["us-east-1a", "us-east-1b"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]

  enable_dns_hostnames = true
  enable_nat_gateway   = true
  single_nat_gateway   = true

  public_subnet_tags = {
    "kubernetes.io/role/elb" = 1
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
  }
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  name               = var.cluster_name
  kubernetes_version = "1.33"

  endpoint_public_access                   = true
  enable_cluster_creator_admin_permissions = true

  addons = {
    eks-pod-identity-agent = {
      most_recent = true
    }

    aws-ebs-csi-driver = {
      most_recent = true

      service_account_role_arn = null
    }
  }

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  eks_managed_node_groups = {
    eks_nodes = {
      # NOTE: enabled only in local env - Use a specific AMI ID for skip validation to localstack
      # ami_id = "ami-0123456789abcdef0"
      ami_type  = "AL2023_ARM_64_STANDARD"
      disk_size = 50

      #In v21, this setting enables the EKS Pod Identity agent. This is the modern replacement for IRSA (IAM Roles for Service Accounts). It allows your applications (pods) to assume IAM roles without needing to manage complex OIDC provider trust relationships or service account annotations.
      create_pod_identity_association = true

      iam_role_attach_cni_policy = true
      iam_role_additional_policies = {
        AmazonEKS_CNI_Policy = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
      }

      capacity_type = "SPOT"

      desired_capacity = 2
      max_capacity     = 3
      min_capacity     = 1

      instance_types = var.instance_types

      tags = {
        Name = var.instance_name
      }
    }
  }

  tags = {
    Environment = var.environment
    Project     = var.cluster_name
  }
}