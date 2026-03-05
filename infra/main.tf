provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      Terraform   = "true"
      Environment = "dev"
      Project     = "tech-challenge"
    }
  }
}

locals {
  cluster_name = "eks-tech-challenge"
}

module "vpc" {
  source       = "../modules/vpc"
  cluster_name = local.cluster_name
}

module "eks" {
  source       = "../modules/eks"
  cluster_name = local.cluster_name
  vpc_id       = module.vpc.vpc_id
  subnet_ids   = module.vpc.private_subnets
  role_arn     = var.cluster_role_arn
}

module "rds" {
  source                = "../modules/rds"
  db_password           = var.db_password
  vpc_id                = module.vpc.vpc_id
  subnet_ids            = module.vpc.private_subnets
  eks_security_group_id = module.eks.cluster_security_group_id
}
