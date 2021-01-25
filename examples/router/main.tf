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

resource "nifcloud_router" "basic" {
  name              = "example"
  description       = "memo"
  availability_zone = "east-21"
  accounting_type   = "2"
  type              = "small"
  security_group    = nifcloud_security_group.basic.group_name

  network_interface {
    network_name = nifcloud_private_lan.basic.private_lan_name
    ip_address   = "192.168.1.1"
  }
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "example"
  availability_zone = "east-21"
  cidr_block        = "192.168.1.0/24"
}

resource "nifcloud_security_group" "basic" {
  group_name        = "example"
  availability_zone = "east-21"
}
