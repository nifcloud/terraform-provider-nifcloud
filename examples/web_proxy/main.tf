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
