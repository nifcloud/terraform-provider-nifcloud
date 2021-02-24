terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_db_instance" "example" {
  accounting_type                = "2"
  availability_zone              = "east-11"
  instance_class                 = "db.large8"
  db_name                        = "baz"
  username                       = "for"
  password                       = "barbarbarupd"
  engine                         = "MySQL"
  engine_version                 = "5.7.15"
  allocated_storage              = 100
  storage_type                   = 0
  identifier                     = "example"
  backup_retention_period        = 2
  binlog_retention_period        = 2
  custom_binlog_retention_period = true
  backup_window                  = "00:00-09:00"
  maintenance_window             = "sun:22:00-sun:22:30"
  multi_az                       = true
  multi_az_type                  = 1
  port                           = 3306
  publicly_accessible            = true
  final_snapshot_identifier      = "example"
  skip_final_snapshot            = false
  read_replica_identifier        = "example-read"
  apply_immediately              = true
}
