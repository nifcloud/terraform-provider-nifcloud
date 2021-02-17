---
page_title: "NIFCLOUD: nifcloud_db_instance"
subcategory: "RDB"
description: |-
  Provides a rdb instance resource.
---

# nifcloud_db_instance

Provides a rdb instance resource.

## Example Usage

```hcl
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

```

## Argument Reference

The following arguments are supported:


* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `allocated_storage` - (Optional) The allocated storage in gibibytes.
* `apply_immediately` - (Optional) Specifies whether any database modifications are applied immediately, or during the next maintenance window. Default is `false`
* `availability_zone` - (Optional) The AZ for the DB instance.
* `backup_retention_period` - (Optional) The days to retain backups for. If `0` automatic backup will be off
* `backup_window` - (Optional) The daily time range (in UTC) during which automated backups are created if they are enabled. Example: `09:46-10:16`
* `binlog_retention_period` - (Optional) The days to retain binlog for. Be sure to specify `custom_binlog_retention_period = true` as a set
* `ca_cert_identifier` - (Optional) The identifier of the CA certificate for the DB instance.
* `custom_binlog_retention_period` - (Optional) The flag of set binary log retention period. Only MySQL can be specified
* `db_name` - (Optional) The name of the database to create when the DB instance is created. If this parameter is not specified, no database is created.
* `db_security_group_name` - (Optional) The security group name to associate with; which can be managed using the nifcloud_db_security_group resource.
* `engine` - (Optional) The database engine. `MySQL` or `postgres` or `MariaDB`
* `engine_version` - (Optional) The database engine version.
* `final_snapshot_identifier` - (Optional) The name of your final DB snapshot when this DB instance is deleted. Must be provided if `skip_final_snapshot` is set to false.
* `identifier` - (Required) The name of the DB instance.
* `instance_class` - (Required) The instance type of the DB instance.
* `maintenance_window` - (Optional) The weekly time range (in UTC) the instance maintenance window. Example: `Sun:05:00-Sun:06:00`
* `master_private_address` - (Optional) Private IP address for master DB.
* `multi_az` - (Optional) If the DB instance is multi AZ enabled.
* `multi_az_type` - (Optional) The type of multi AZ. (0: Data priority, 1: Performance priority) default `0`
* `network_id` - (Optional) The id of private lan.
* `parameter_group_name` - (Optional) Name of the DB parameter group to associate; which can be managed using the nifcloud_db_parameter_group resource.
* `password` - (Optional) Password for the master DB user.
* `port` - (Optional) The database port.
* `publicly_accessible` - (Optional) Bool to control if instance is publicly accessible. Default is `true`
* `read_replica_identifier` - (Optional) The DB instance name for read replica.
* `read_replica_private_address` - (Optional) Private IP address for read replica.
* `replicate_source_db` - (Optional) Specifies that this resource is a Replicate database, and to use this value as the source database.
* `restore_to_point_in_time` - (Optional) A configuration block for restoring a DB instance to an arbitrary point in time See [this](#restore-to-point-in-time).
* `skip_final_snapshot` - (Optional) Determines whether a final DB snapshot is created before the DB instance is deleted. Defaults to `true` no DBSnapshot is created
* `slave_private_address` - (Optional) Private IP address for master DB.
* `snapshot_identifier` - (Optional) Specifies whether or not to create this database from a snapshot.
* `storage_type` - (Optional) One of `0` (HDD), or `1` (Flash drive). The default is `0`
* `username` - (Optional) Username for the master DB user.
* `virtual_private_address` - (Optional) Private IP address for virtual load balancer.

### `restore_to_point_in_time`

* `restore_time` - (Optional) The date and time to restore from. Value must be a time in Universal Coordinated Time (UTC) format.
* `source_db_instance_identifier` - (Required) The identifier of the source DB instance from which to restore.
* `use_latest_restorable_time` - (Optional) A boolean value that indicates whether the DB instance is restored from the latest backup time. Defaults to `false`

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `address` - The hostname of the DB instance.


## Import

nifcloud_db_instance can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_db_instance.example foo
```
