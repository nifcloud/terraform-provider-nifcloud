---
page_title: "NIFCLOUD: nifcloud_customer_gateway"
subcategory: "Computing"
description: |-
  Provides a customer gateway resource.
---

# nifcloud_customer_gateway

Provides a customer gateway resource.

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

resource "nifcloud_customer_gateway" "example" {
  customer_gateway_name        = "cgw002"
  customer_gateway_description = "memo"
  type                         = "IPsec"
  ip_address                   = "192.168.0.1"
  lan_side_ip_address          = "192.168.0.1"
  lan_side_cidr_block          = "192.168.0.0/28"
}
```

## Argument Reference

The following arguments are supported:


* `customer_gateway_name` - (Optional) The customer gateway name.
* `customer_gateway_description` - (Optional) The customer gateway description.
* `type` - (Optional) The type.
* `ip_address` - (Required) The IP address.
* `lan_side_ip_address` - (Optional) The LAN side IP address.
* `lan_side_cidr_block` - (Optional) The LAN side CIDR block.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `customer_gateway_id` - The customer gateway ID.
* `state` - The state.


## Import

nifcloud_customer_gateway can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_customer_gateway.example foo
```
