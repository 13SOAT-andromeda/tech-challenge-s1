module "network" {
  source       = "./modules/network"
  vpc_cidr     = var.vpc_cidr
  cluster_name = var.cluster_name
}

module "eks" {
  source             = "./modules/eks"
  cluster_name       = var.cluster_name
  kubernetes_version = var.kubernetes_version
  vpc_id             = module.network.vpc_id
  public_subnet_ids  = module.network.public_subnet_ids
  private_subnet_ids = module.network.private_subnet_ids

  # Pass the lab role ARN if provided; tell module not to create role when using lab role
  existing_cluster_role_arn = var.lab_cluster_role_arn
  create_cluster_role       = var.lab_cluster_role_arn == "" ? true : false
  attach_cluster_policies   = var.lab_attach_cluster_policies
}