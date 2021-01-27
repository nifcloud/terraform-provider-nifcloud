---
page_title: "NIFCLOUD: nifcloud_web_proxy"
subcategory: "Computing"
description: |-
  Provides a web proxy resource.
---

# nifcloud_web_proxy

Provides a web proxy resource.

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
  region = "jp-east-2"
}

resource "nifcloud_router" "example" {
  availability_zone = "east-21"

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }

  network_interface {
    network_id = nifcloud_private_lan.example.id
  }
}

resource "nifcloud_web_proxy" "example" {
  router_id                   = nifcloud_router.example.id
  bypass_interface_network_id = nifcloud_private_lan.example.id
  listen_interface_network_id = "net-COMMON_GLOBAL"
  listen_port                 = "8080"
  description                 = "memo2"
  name_server                 = "1.1.1.1"
}

resource "nifcloud_private_lan" "example" {
  private_lan_name  = "example"
  availability_zone = "east-21"
  cidr_block        = "192.168.100.0/24"
}

```

## Argument Reference

The following arguments are supported:


* `bypass_interface_network_id` - (Optional) The id for the by pass network.
* `bypass_interface_network_name` - (Optional) The name for the by pass network
* `description` - (Optional) The web proxy description.
* `listen_interface_network_id` - (Optional) The id for the listen network. listen_interface_network_name and either is required.
* `listen_interface_network_name` - (Optional) The name for the listen network. listen_interface_network_id and either is required.
* `listen_port` - (Required) The port of web proxy.
* `name_server` - (Optional) The ip address for dns server.
* `router_id` - (Optional) The id for the router. router_name and either is required.
* `router_name` - (Optional) The name for the router. route_id and either is required.

## Import

nifcloud_web_proxy can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_web_proxy.example foo
```
