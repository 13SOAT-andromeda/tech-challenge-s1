output "instance_hostname" {
  value       = aws_instance.app_server.private_dns     # The actual value to be outputted
  description = "Private DNS name of the EC2 instance." # Description of what this output represents
}