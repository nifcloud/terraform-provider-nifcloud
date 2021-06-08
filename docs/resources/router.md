---
page_title: "NIFCLOUD: nifcloud_router"
subcategory: "Network"
description: |-
  Provides a router resource.
---

# nifcloud_router

Provides a router resource.

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

resource "nifcloud_router" "router" {
  name       = "router"
  availability_zone = "east-12"
  security_group    = nifcloud_security_group.router.group_name
  type              = "small"
  accounting_type   = "2"

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }

  network_interface {
    network_name = nifcloud_private_lan.pri.private_lan_name
    ip_address   = "192.168.1.1"
    dhcp         = false
  }
}

resource "nifcloud_security_group" "router" {
  group_name        = "routerfw"
  availability_zone = "east-12"
}

resource "nifcloud_private_lan" "pri" {
  private_lan_name  = "pri"
  availability_zone = "east-12"
  cidr_block        = "192.168.1.0/24"
  accounting_type   = "2"
}
```

## Argument Reference

The following arguments are supported:

* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `availability_zone` - (Optional) The availability zone.
* `description` - (Optional) The router description.
* `name` - (Optional) The router name.
* `nat_table_id` - (Optional) The ID of the NAT table to attach.
* `network_interface` - (Required) The network interface list. see [network interface](#network-interface).
* `route_table_id` - (Optional) The ID of associated route table.
* `security_group` - (Optional) The security group name to associate with; which can be managed using the nifcloud_security_group resource.
* `type` - (Optional) The type of the router. Valid types are `small`, `medium`, `large`.

### network interface

* `dhcp` - (Optional) The flag to enable or disable DHCP.
* `dhcp_config_id` - (Optional) The ID of the DHCP config to attach.
* `dhcp_options_id` - (Optional) The ID of the DHCP options to attach.
* `ip_address` - (Optional) The IP address to select from `static` or `elastic IP address` or `static IP address`; Default(null) is DHCP.
* `network_id` - (Optional) The ID of the network to attach; `net-COMMON_GLOBAL` or `net-COMMON_PRIVATE` or `private lan network id` .
* `network_name` - (Optional) The private lan name of the network to attach.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `router_id` - The unique ID of the router.
* `route_table_association_id` - The ID of the route table association.
* `nat_table_association_id` - The ID of the NAT table association.

## Import

nifcloud_router can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_router.example foo
```

## Notice

This provider does not support upgrading router.
If you want to upgrade the router, please upgrade from control panel and also you must delete backup of router after upgrade.
