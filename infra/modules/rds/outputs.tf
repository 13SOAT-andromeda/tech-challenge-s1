output "db_endpoint" {
  description = "Endpoint da instância RDS"
  value       = aws_db_instance.main.endpoint
}

output "db_address" {
  description = "Hostname"
  value       = aws_db_instance.main.address
}

output "db_port" {
  description = "Porta"
  value       = aws_db_instance.main.port
}

output "db_name" {
  description = "Database name"
  value       = aws_db_instance.main.db_name
}

output "db_username" {
  description = "Database user"
  value       = aws_db_instance.main.username
  sensitive   = true
}