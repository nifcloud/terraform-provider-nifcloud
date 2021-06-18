---
page_title: "NIFCLOUD: nifcloud_hatoba_cluster"
subcategory: "Hatoba"
description: |-
  Provides a Hatoba cluster resource.
---

# nifcloud_hatoba_cluster

Provides a Hatoba cluster resource.

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

resource "nifcloud_hatoba_cluster" "example" {
  name           = "cluster001"
  description    = "memo"
  firewall_group = nifcloud_hatoba_firewall_group.example.name
  locations      = ["east-11"]

  addons_config {
    http_load_balancing {
      disabled = true
    }
  }

  node_pools {
    name          = "default"
    instance_type = "medium"
    node_count    = 1
  }
}

resource "nifcloud_hatoba_firewall_group" "example" {
  name = "group001"
}
```

## Argument Reference

The following arguments are supported:

* `addons_config` - (Optional) The configs for Kubernetes addons. see [addons_config](#addons_config)
* `description` - (Optional) The Hatoba cluster description.
* `firewall_group` - (Required) The firewall group name to associate with; which can be managed using the nifcloud_hatoba_firewall_group resource.
* `kubernetes_version` - (Optional) The version of Kubernetes.
* `locations` - (Required) The cluster location.
* `name` - (Required) The name for the Hatoba cluster.
* `network_config` - (Optional) The configs for network. see [network_config](#network_config)
* `node_pools` - (Required) The node pool config. see [node_pools](#node_pools)

### addons_config

#### Arguments

* `http_load_balancing` - (Optional) The configs for HTTP load balancer. It can set only `disabled` parameter.

### network_config

#### Arguments

* `network_id` - (Optional) The ID of private LAN.

### node_pools

#### Arguments

* `instance_type` - (Required) The instance type for node pool.
* `name` - (Required) The name of node pool.
* `node_count` - (Required) The desired node count in this node pool.

## Attributes Reference

* `cluster.node_pools.*.nodes` - The list of node information. see [node](#node)

### node

* `availability_zone` - The availability zone where the node located.
* `name` - The name of the node.
* `public_ip_address` - The public IP address of the node.
* `private_ip_address` - The private IP address of the node.

## Import

nifcloud_hatoba_cluster can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_hatoba_cluster.example foo
```
