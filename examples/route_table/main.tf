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
