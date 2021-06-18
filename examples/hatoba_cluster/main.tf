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
