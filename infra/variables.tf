variable "cluster_role_arn" {
  description = "IAM role ARN for the EKS cluster"
  type        = string
  default     = "LabEksClusterRole"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}
