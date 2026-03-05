variable "db_name" {
  description = "Name of the database"
  type        = string
  default     = "garagedb"
}

variable "db_user" {
  description = "Username for the database"
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "Password for the database"
  type        = string
  sensitive   = true
}

variable "vpc_id" {
  description = "VPC ID where the RDS will be deployed"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet IDs for the RDS instance"
  type        = list(string)
}

variable "eks_security_group_id" {
  description = "Security group ID of the EKS cluster to allow ingress"
  type        = string
}
