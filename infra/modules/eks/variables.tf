variable "cluster_name" {}
variable "kubernetes_version" {}
variable "vpc_id" {}
variable "private_subnet_ids" { type = list(string) }
variable "public_subnet_ids" { type = list(string) }

# New variables to support using an existing Lab EKS cluster role instead of creating one
variable "existing_cluster_role_arn" {
  description = "Optional ARN of an existing IAM role to use for the EKS cluster (e.g. LabEksClusterRole). If set, the module will use this role instead of creating a new one."
  type        = string
  default     = ""
}

variable "create_cluster_role" {
  description = "When true the module will create the IAM role for the EKS cluster. Set to false to prevent role creation (use existing_cluster_role_arn)."
  type        = bool
  default     = true
}

variable "attach_cluster_policies" {
  description = "When true the module will attach required managed policies to the cluster role. Set to false if your Labs account cannot attach policies or if attachments are managed elsewhere."
  type        = bool
  default     = true
}
