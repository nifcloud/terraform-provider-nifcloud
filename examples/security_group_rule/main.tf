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

resource "nifcloud_security_group" "example" {
  group_name        = "allowtcp"
  availability_zone = "east-11"
}

resource "nifcloud_security_group_rule" "example" {
  security_group_names = [nifcloud_security_group.example.group_name]
  type                 = "IN"
  from_port            = 0
  to_port              = 65535
  protocol             = "TCP"
  cidr_ip              = "0.0.0.0/0"
  description          = "メモです"
}
