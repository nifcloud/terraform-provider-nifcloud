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

resource "nifcloud_hatoba_firewall_group" "example" {
  name        = "group001"
  description = "memo"

  rule {
    protocol    = "TCP"
    direction   = "IN"
    from_port   = 80
    to_port     = 80
    cidr_ip     = "0.0.0.0/0"
    description = "http access"
  }
}
