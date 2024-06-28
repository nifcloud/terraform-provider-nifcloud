---
page_title: "NIFCLOUD: nifcloud_devops_firewall_group"
subcategory: "DevOps with GitLab"
description: |-
  Provides a DevOps firewall group resource.
---

# nifcloud_devops_firewall_group

Provides a DevOps firewall group resource.

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

resource "nifcloud_devops_firewall_group" "example" {
  name              = "example"
  availability_zone = "east-11"
  description       = "memo"

  rule {
    protocol    = "TCP"
    port        = 443
    cidr_ip     = "192.168.1.0/24"
    description = "https from pri"
  }

  rule {
    protocol    = "TCP"
    port        = 22
    cidr_ip     = "192.168.1.0/24"
    description = "ssh from pri"
  }

  rule {
    protocol    = "ICMP"
    cidr_ip     = "192.168.1.0/24"
    description = "icmp from pri"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The availability zone for the DevOps firewall group.
* `description` - (Optional) Description of the DevOps firewall group.
* `name` - (Required) The name of the DevOps firewall group.
* `rule` - (Optional) List of the DevOps firewall rules. see [rule](#rule).

### rule

* `cidr_ip` - (Required) CIDR block or IPv4 address.
* `description` - (Optional) Description of the rule.
* `port` - (Optional, but required if `protocol` is `TCP`) Port. Valid values are `22` or `443`.
* `protocol` - (Required) Protocol. Valid values are `TCP` or `ICMP`.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

### rule

* `id` - ID of the rule.

## Import

nifcloud_devops_firewall_group can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_devops_firewall_group.example foo
```
