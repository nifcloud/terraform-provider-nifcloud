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
  type                               = "IPsec"
  ip_address                         = "192.168.0.1"
  nifty_customer_gateway_name        = "cgw002"
  nifty_lan_side_ip_address          = "192.168.0.1"
  nifty_lan_side_cidr_block          = "192.168.0.0/28"
  nifty_customer_gateway_description = "memo"
}
```

## Argument Reference

The following arguments are supported:


* `type` - (Optional) The type.
* `ip_address` - (Required) The IP address.
* `nifty_customer_gateway_name` - (Optional) The nifty customer gateway name.
* `nifty_lan_side_ip_address` - (Optional) The nifty LAN side IP address.
* `nifty_lan_side_cidr_block` - (Optional) The nifty LAN side CIDR block.
* `nifty_customer_gateway_description` - (Optional) The nifty customer gateway description.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `customer_gateway_id` - The customer gateway ID.
* `state` - The state.


## Import

nifcloud_customer_gateway can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_customer_gateway.example foo
```
