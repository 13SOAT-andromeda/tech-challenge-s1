resource "aws_db_subnet_group" "main" {
  name = "${var.db_name}-subnet-group"  
  subnet_ids = var.private_subnets

  tags = {
    Name = "${var.db_name}-subnet-group"
  }
}

resource "aws_security_group" "rds" {
  name = "${var.db_name}-sg"
  vpc_id = var.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [var.eks_cluster_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  tags = {
    Name = "${var.db_name}-sg"
  }
}

resource "aws_db_instance" "main" {
  identifier = "db-${var.db_name}"
  
  db_name = "${var.db_name}"
  username = "${var.db_user}"
  password = "${var.db_pass}"
  
  engine = "postgres"
  engine_version = "${var.engine_version}"
  instance_class = var.instance_class

  allocated_storage = var.allocated_storage
  storage_type          = "gp3"
  storage_encrypted      = true

  db_subnet_group_name = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  publicly_accessible = false

  skip_final_snapshot = false
  final_snapshot_identifier = "${var.db_name}-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  backup_retention_period = 7


   tags = {
    Name = "db-${var.db_name}"
  }
}
