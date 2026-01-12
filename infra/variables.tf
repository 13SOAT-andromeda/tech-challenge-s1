variable "aws_region" {
  default = "us-east-1"
}

variable "vpc_cidr" {
  default = "10.0.0.0/16"
}

variable "cluster_name" {
  default = "tech-challenge-api"
}

variable "kubernetes_version" {
  default = "1.31"
}

# New root-level variable to supply the Lab EKS cluster role ARN (leave empty to let module create the role)
variable "lab_cluster_role_arn" {
  description = "Optional ARN of a pre-provisioned Lab EKS cluster role (LabEksClusterRole). When set, the module will use this role and won't create a new cluster role."
  type        = string
  default     = ""
}

# New root-level toggle to indicate whether the module should attempt to attach cluster policies
variable "lab_attach_cluster_policies" {
  description = "When false the module will not attempt to attach managed policies to the cluster role (useful when using restricted Lab accounts)."
  type        = bool
  default     = false
}
