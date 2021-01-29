---
page_title: "NIFCLOUD: nifcloud_nat_table"
subcategory: "Network"
description: |-
  Provides a nat table resource.
---

# nifcloud_nat_table

Provides a nat table resource.

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

resource "nifcloud_nat_table" "nat" {
  snat {
    rule_number                   = "1"
    description                   = "memo"
    protocol                      = "TCP"
    source_address                = "192.0.2.1"
    source_port                   = 80
    translation_port              = 81
    outbound_interface_network_id = "net-COMMON_PRIVATE"
  }

  dnat {
    rule_number                    = "1"
    description                    = "memo"
    protocol                       = "ALL"
    translation_address            = "192.168.1.1"
    inbound_interface_network_name = nifcloud_private_lan.pri.private_lan_name
  }
}

resource "nifcloud_private_lan" "pri" {
  private_lan_name  = "pri"
  availability_zone = "east-11"
  cidr_block        = "192.168.1.0/24"
}

```

## Argument Reference

The following arguments are supported:

* `snat` - (Optional) A list of snat objects. see [snat](#snat).
* `dnat` - (Optional) A list of dnat objects. see [dnat](#dnat).

### snat

#### Arguments

* `rule_number` - (Required) The rule number.
* `description` - (Optional) The nat table rule description.
* `protocol` - (Required) The protocol.
  * Specifiable protocol: [ALL / TCP / UDP / TCP_UDP / ICMP]
* `source_address` - (Required) The source address.
* `source_port` - (Optional) The source port.
* `translation_port` - (Optional) The translation port.
* `outbound_interface_network_id` - (Optional) The outbound interface network id; `net-COMMON_GLOBAL` or `net-COMMON_PRIVATE` or private lan network id.
* `outbound_interface_network_name` - (Optional) The private lan name of target outbound interface network.

### dnat

#### Arguments

* `rule_number` - (Required) The rule number.
* `description` - (Optional) The nat table rule description.
* `protocol` - (Required) The protocol.
  * Specifiable protocol: [ALL / TCP / UDP / TCP_UDP / ICMP]
* `destination_port` - (Optional) The destination port.
* `translation_address` - (Required) The translation address.
* `translation_port` - (Optional) The translation port.
* `inbound_interface_network_id` - (Optional) The inbound interface network id; `net-COMMON_GLOBAL` or `net-COMMON_PRIVATE` or private lan network id.
* `inbound_interface_network_name` - (Optional) The private lan name of target inbound interface network.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `nat_table_id` - The id of nat table.

## Import

nifcloud_nat_table can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_nat_table.example foo
````
