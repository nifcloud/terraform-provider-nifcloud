---
page_title: "NIFCLOUD: nifcloud_db_security_group"
subcategory: "RDB"
description: |-
  Provides a RDB security group resource.
---

# nifcloud_db_security_group

Provides a RDB security group resource.

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

resource "nifcloud_db_security_group" "example" {
  group_name        = "dbgroup001"
  description       = "memo"
  availability_zone = "east-11"

  rule {
    cidr_ip = "0.0.0.0/0"
  }

  rule {
    security_group_name = nifcloud_security_group.example.group_name
  }
}

resource "nifcloud_security_group" "example" {
  group_name        = "group001"
  availability_zone = "east-11"
}

```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) The availability zone.
* `description` - (Optional) The db security group description.
* `group_name` - (Required) The name for the db security group.
---
* `rule` - (Optional) A list of the db security group rule objects. see [rule](#rule).

### rule

#### Arguments

* `cidr_ip` - (Optional) The CIDR IP Address that allow access. Cannot be specified with `security_group_name` .
* `security_group_name` - (Optional) The security group name that allow access. Cannot be specified with `cidr_ip` .

## Import

nifcloud_db_security_group can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_db_security_group.example foo
```