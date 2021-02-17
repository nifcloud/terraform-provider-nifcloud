---
page_title: "NIFCLOUD: nifcloud_db_parameter_group"
subcategory: "RDB"
description: |-
  Provides a DB parameter group resource.
---

# nifcloud_db_parameter_group

Provides a DB parameter group resource.

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

resource "nifcloud_db_parameter_group" "example" {
  name        = "example"
  family      = "mysql5.7"
  description = "memo"

  parameter {
    name  = "character_set_server"
    value = "utf8"
  }

  parameter {
    name  = "character_set_client"
    value = "utf8"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DB parameter group.
* `family` - (Required) The DB parameter group family name.
* `description` - (Optional) The description for the DB parameter group.
* `parameter` - (Optional) A list of parameters. see [parameter](#parameter)

### parameter

* `name` - (Required) The name of the parameter.
* `value` - (Required) The value of the parameter.
* `apply_method` - (Optional) The time to apply the paramete updates. Valid methods are `immediate` or `pending-reboot`. Default is `immediate`.

## Import

nifcloud_db_parameter_group can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_db_parameter_group.example foo
```
