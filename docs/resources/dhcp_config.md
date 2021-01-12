---
page_title: "NIFCLOUD: nifcloud_dhcp_config"
subcategory: "Computing"
description: |-
  Provides a dhcp config resource.
---

# nifcloud_dhcp_config

Provides a dhcp config resource.

# Example Usage

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

resource "nifcloud_dhcp_config" "example" {
    static_mapping {
        static_mapping_ipaddress = "192.168.1.10"
        static_mapping_macaddress = "00:00:5e:00:53:00"
        static_mapping_description = "static-mapping-memo"
    }
    ipaddress_pool {
        ipaddress_pool_start = "192.168.2.1"
        ipaddress_pool_stop = "192.168.2.100"
        ipaddress_pool_description = "ipaddress-pool-memo"
    }
}
```

## Argument Reference

The following arguments are supported:

* `static_mapping` - (Optional) A list of static mapping ip address. see [static_mapping](#static_mapping)
* `ipaddress_pool` - (Optional) A list of ipaddress pool. see [ipaddress_pool](#ipaddress_pool)

### static_mapping

#### Arguments

* `static_mapping_ipaddress` - (Required) The static mapping IP address.
* `static_mapping_macaddress` - (Required) The static mapping MAC address.
* `static_mapping_descreption` - (Optional) The static mapping IP address description.

### ipaddress_pool

#### Arguments

* `ipaddress_pool_start` - (Required) The start IP address of ipAddressPool.
* `ipaddress_pool_stop` - (Required) The stop IP address of ipAddressPool.
* `ipaddress_pool_description` - (Optional) The ipaddress pool description.

## Import

nifcloud_dhcp_config can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_dhcp_config.example foo
```
