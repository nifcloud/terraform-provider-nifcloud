---
page_title: "NIFCLOUD: nifcloud_private_lan"
subcategory: "Computing"
description: |-
  Provides a private LAN resource.
---

## nifcloud_private_lan

Provides a private LAN resource.

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

resource "nifcloud_private_lan" "pri" {
  private_lan_name  = "pri"
  description       = "pri lan 01"
  availability_zone = "east-11"
  cidr_block        = "192.168.1.0/24"
  accounting_type   = "2"
}
```

## Argument Reference

The following arguments are supported:

* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `availability_zone` - (Optional) The availability zone.
* `cidr_block ` - (Required) The CIDR IP Address Block.
* `description` - (Optional) The private LAN description.
* `private_lan_name` - (Optional) The license name.

## Import

nifcloud_private_lan can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_private_lan.example foo
```
