---
page_title: "NIFCLOUD: nifcloud_load_balancer"
subcategory: "Computing"
description: |-
  Provides a load balancer resource.
---

## nifcloud_load_balancer

Provides a load balancer resource.

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

resource "nifcloud_load_balancer" "l4lb" {
  accounting_type = "1"
  load_balancer_name = "nl4lb"
}
```

## Argument Reference

The following arguments are supported:

* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `availability_zone` - (Optional) The availability zone.
* `ip_version` - (Optional) The IP version. (v4 or v6).
* `load_balancer_name` - (Required) The load balancer name.
* `network_volume` - (Optional) The network volume.
* `policy_type` - (Optional) The policy type. (standard or ats).

## Import

nifcloud_load_balancer can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_load_balancer.example foo
```
