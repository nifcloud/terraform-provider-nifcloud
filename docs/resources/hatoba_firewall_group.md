---
page_title: "NIFCLOUD: nifcloud_hatoba_firewall_group"
subcategory: "Hatoba"
description: |-
  Provides a Hatoba firewall group resource.
---

# nifcloud_hatoba_firewall_group

Provides a Hatoba firewall group resource.

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

resource "nifcloud_hatoba_firewall_group" "example" {
  name        = "group001"
  description = "memo"

  rule {
    protocol    = "TCP"
    direction   = "IN"
    from_port   = 80
    to_port     = 80
    cidr_ip     = "0.0.0.0/0"
    description = "http access"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The Hatoba firewall group description.
* `name` - (Required) The name for the Hatoba firewall group.
---
* `rule` - (Optional) A list of the Hatoba firewall group rule objects. see [rule](#rule).

### rule

#### Arguments

* `cidr_ip` - (Required) The CIDR IP address that allow access.
* `description` - (Optional) The firewall group rule description.
* `direction` - (Optional) The direction of rule being created. Valid options are IN (Incoming) or OUT (Outgoing).
* `from_port` - (Optional) The start port.
* `protocol` - (Optional) The protocol.
* `to_port` - (Optional) The end port.

## Attributes Reference

* `rule.*.id` - The identifier of each firewall group rule.

## Import

nifcloud_hatoba_firewall_group can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_hatoba_firewall_group.example foo
```
