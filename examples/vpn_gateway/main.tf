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

resource "nifcloud_vpn_gateway" "vpngw" {
  accounting_type             = "1"
  description = "vpn gateway by terraform"
  name      = "vpngw001"
  type      = "small"
  availability_zone = "east-11"
  network_name = nifcloud_private_lan.basic.private_lan_name
  ip_address = "192.168.1.1"
  security_group = nifcloud_security_group.basic.group_name
  route_table_id = nifcloud_route_table.rt.route_table_id
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "example"
  availability_zone = "east-11"
  cidr_block        = "192.168.1.0/24"
}

resource "nifcloud_security_group" "basic" {
  group_name        = "example"
  availability_zone = "east-11"
}

resource "nifcloud_route_table" "rt" {
  route {
    cidr_block = "10.0.1.0/24"
    ip_address = "192.168.2.0"
  }

  route {
    cidr_block = "10.0.2.0/24"
    network_id = "net-COMMON_GLOBAL"
  }
}
