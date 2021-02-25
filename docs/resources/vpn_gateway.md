---
page_title: "NIFCLOUD: nifcloud_vpn_gateway"
subcategory: "Network"
description: |-
  Provides a vpn gateway resource.
---

# nifcloud_vpn_gateway

Provides a vpn gateway resource.

## Example Usage

```hd
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

resource "nifcloud_vpn_gateway" "vpngw" {
  accounting_type             = "1"
  name      = "vpngw001"
  type      = "small"
  availability_zone = "east-11"
  network_name = nifcloud_private_lan.basic.private_lan_name
  ip_address = "192.168.1.1"
  security_group = nifcloud_security_group.basic.group_name
  route_table_id = nifcloud_route_table.rt.route_table_id
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "example"
  availability_zone = "east-11"
  cidr_block        = "192.168.1.0/24"
}

resource "nifcloud_security_group" "basic" {
  group_name        = "example"
  availability_zone = "east-11"
}

resource "nifcloud_route_table" "rt" {
  route {
    cidr_block = "10.0.1.0/24"
    ip_address = "192.168.2.0"
  }

  route {
    cidr_block = "10.0.2.0/24"
    network_id = "net-COMMON_GLOBAL"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name for the vpn gateway.
* `type` - (Optional) The type of vpn gateway.
* `availability_zone` - (Optional) The availability zone.
* `accounting_type` - (Optional) The accounting type.
* `description` - (Optional) The vpn gateway description.
* `network_id` - (Optional) The id for the network.
* `network_name` - (Optional) The name for the network.
* `ip_address` - (Optional) The private ip address.
* `security_group` - (Optional) The name of firewall group.
* `route_table_id` - (Optional) The ID of the route table to attach.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `vpn_gateway_id` - The id for the vpn gateway.
* `route_table_association_id` - The ID of the route table association.

## Import

nifcloud_vpn_gateway can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_vpn_gateway.example foo
```
