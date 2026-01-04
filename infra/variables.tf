variable "aws_region" {
  default = "us-east-1"
}

variable "cluster_name" {
  default = "tech-challenge-api"
}

variable "kubernetes_version" {
  default = "1.33"
}

variable "vpc_cidr" {
  default = "10.0.0.0/16"
}

variable "vpc_cni_version" {
  description = "Version of the VPC CNI add-on"
  type        = string
  default     = "v1.19.0-eksbuild.1"
}
variable "coredns_version" {
  description = "Version of the CoreDNS add-on"
  type        = string
  default     = "v1.11.1-eksbuild.9" # Versão estável comum
}
variable "kube_proxy_version" {
  description = "Version of the kube-proxy add-on"
  type        = string
  default     = "v1.30.0-eksbuild.1" # Ajustado para compatibilidade geral
}
variable "ebs_csi_version" {
  description = "Version of the EBS CSI driver add-on"
  type        = string
  default     = "v1.31.0-eksbuild.1"
}