variable "aws_region" {
  description = "AWS region for resources"
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "Name of the EKS cluster"
  default     = "tech-challenge-api"
}

variable "kubernetes_version" {
  description = "Kubernetes version"
  default     = "1.31" # Updated to a stable, supported version
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  default     = "10.0.0.0/16"
}