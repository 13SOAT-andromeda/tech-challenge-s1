variable "cluster_name" {
  type        = string
  default     = "tech-challenge-api"
  description = "The name of the project"
}

variable "environment" {
  type        = string
  default     = "dev"
  description = "The deployment environment (e.g., dev, staging, prod)"
}

variable "region" {
  type        = string
  default     = "us-east-1"
  description = "The AWS region to deploy resources in"
}

variable "instance_name" {
  type        = string
  default     = "tech-challenge-api-node-group"
  description = "The name of the EC2 instance"
}

variable "instance_types" {
  type        = list(string)
  default     = ["t4g.medium"]
  description = "The type of EC2 instance"
}