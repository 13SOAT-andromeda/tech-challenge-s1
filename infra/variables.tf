variable "instance_name" {
  type        = string                         # The type of the variable, in this case a string
  default     = "learn-terraform"              # Default value for the variable
  description = "The name of the EC2 instance" # Description of what this variable represents
}

variable "instance_type" {
  type        = string                     # The type of the variable, in this case a string
  default     = "t2.micro"                 # Default value for the variable
  description = "The type of EC2 instance" # Description of what this variable represents
}