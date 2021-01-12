---
page_title: "NIFCLOUD: nifcloud_dhcp_option"
subcategory: "Computing"
description: |-
  Provides a dhcp option resource.
---

# nifcloud_dhcp_option

Provides a dhcp option resource.

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

resource "nifcloud_dhcp_option" "example" {
    default_router = "192.168.0.1"
    domain_name = "example.com"
    domain_name_servers = ["192.168.0.1", "192.168.0.2"]
    ntp_servers = ["192.168.0.1"]
    netbios_name_servers = ["192.168.0.1", "192.168.0.2"]
    netbios_node_type = "1"
    lease_time = "600"
}
```

## Argument Reference

The following arguments are supported:

* `default_router` - (Optional) The IP address of default gateway.
* `domain_name` - (Optional) The domain name used by the client in host name resolution.
* `domain_name_servers` - (Optional) The IP address list of the DNS server.
* `ntp_servers` - (Optional) The IP address list of the NTP server.
* `netbios_name_servers` - (Optional) The IP address list of the NetBIOS server.
* `netbios_node_type` - (Optional) The NetBIOS node type. (1: Don't use WINS, 2: Don't use broadcast, 4: Priorirtize broadcasting, 8: Prioritize WINS)
* `lease_time` - (Optional) The IP address lease time.　(Unit: second)

## Import

nifcloud_dhcp_option can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_dhcp_option.example foo
```
