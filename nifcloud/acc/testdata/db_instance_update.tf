provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_db_instance" "basic" {
  accounting_type                = "2"
  availability_zone              = "east-21"
  instance_class                 = "db.large8"
  db_name                        = "baz"
  username                       = "for"
  password                       = "barbarbarupd"
  engine                         = "MySQL"
  engine_version                 = "5.7.15"
  allocated_storage              = 100
  storage_type                   = 0
  identifier                     = "%supd"
  backup_retention_period        = 2
  binlog_retention_period        = 2
  custom_binlog_retention_period = true
  backup_window                  = "00:00-09:00"
  maintenance_window             = "sun:22:00-sun:22:30"
  multi_az                       = false
  port                           = 3306
  publicly_accessible            = true
  final_snapshot_identifier      = "%s"
  skip_final_snapshot            = false
  db_security_group_name         = nifcloud_db_security_group.upd.id
  parameter_group_name           = nifcloud_db_parameter_group.upd.id
  apply_immediately              = true
}

resource "nifcloud_db_instance" "base" {
  accounting_type                = "1"
  availability_zone              = "east-21"
  instance_class                 = "db.large"
  db_name                        = "baz"
  username                       = "for"
  password                       = "barbarbar"
  engine                         = "MySQL"
  engine_version                 = "5.7.15"
  allocated_storage              = 50
  storage_type                   = 0
  identifier                     = "%s-base"
  backup_retention_period        = 1
  binlog_retention_period        = 1
  custom_binlog_retention_period = true
  backup_window                  = "00:00-08:00"
  maintenance_window             = "sun:23:00-sun:23:30"
  multi_az                       = true
  port                           = 3306
  publicly_accessible            = true
  db_security_group_name         = nifcloud_db_security_group.basic.id
  parameter_group_name           = nifcloud_db_parameter_group.basic.id
  skip_final_snapshot            = true
  apply_immediately              = true
}

resource "nifcloud_db_instance" "replica" {
  identifier          = "%s-replica"
  replicate_source_db = nifcloud_db_instance.base.id
  storage_type        = nifcloud_db_instance.base.storage_type
  instance_class      = nifcloud_db_instance.base.instance_class
  accounting_type     = nifcloud_db_instance.base.accounting_type
  skip_final_snapshot = true
}

resource "nifcloud_db_instance" "restore" {
  restore_to_point_in_time {
    source_db_instance_identifier = nifcloud_db_instance.base.id
    use_latest_restorable_time    = true
  }

  identifier             = "%s-restore"
  availability_zone      = nifcloud_db_instance.base.availability_zone
  instance_class         = nifcloud_db_instance.base.instance_class
  accounting_type        = nifcloud_db_instance.base.accounting_type
  storage_type           = nifcloud_db_instance.base.storage_type
  multi_az               = nifcloud_db_instance.base.multi_az
  publicly_accessible    = nifcloud_db_instance.base.publicly_accessible
  db_security_group_name = nifcloud_db_instance.base.db_security_group_name
  parameter_group_name   = nifcloud_db_instance.base.parameter_group_name
  port                   = nifcloud_db_instance.base.port
  apply_immediately      = true
  skip_final_snapshot    = true
}

resource "nifcloud_db_parameter_group" "basic" {
  name   = "%s"
  family = "mysql5.7"
}

resource "nifcloud_db_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}

resource "nifcloud_db_parameter_group" "upd" {
  name   = "%supd"
  family = "mysql5.7"
}

resource "nifcloud_db_security_group" "upd" {
  group_name        = "%supd"
  availability_zone = "east-21"
}
