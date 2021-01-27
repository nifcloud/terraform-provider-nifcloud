---
page_title: "NIFCLOUD: nifcloud_route_table"
subcategory: "Network"
description: |-
  Provides a route table resource.
---

# nifcloud_route_table

Provides a route table resource.

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

resource "nifcloud_route_table" "r" {
  route {
    cidr_block = "10.0.1.0/24"
    ip_address = "192.168.0.1"
  }

  route {
    cidr_block = "10.0.2.0/24"
    network_id = "net-COMMON_GLOBAL"
  }
}

```

## Argument Reference

The following arguments are supported:

* `route` - (Optional) A list of route objects. see [route](#route).

### route

#### Arguments

* `cidr_block` - (Required) The destination IP address or CIDR.
* `ip_address` - (Optional) The target IP address.
* `network_id` - (Optional) The id of target network; 'net-COMMON_GLOBAL' or `net-COMMON_PRIVATE` or private lan network id.
* `network_name` - (Optional) The private lan name of target network.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `route_table_id` - The id of route table.


## Import

nifcloud_route_table can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_route_table.example foo
```
