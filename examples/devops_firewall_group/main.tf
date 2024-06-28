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

resource "nifcloud_devops_firewall_group" "example" {
  name              = "example"
  availability_zone = "east-11"
  description       = "memo"

  rule {
    protocol    = "TCP"
    port        = 443
    cidr_ip     = "192.168.1.0/24"
    description = "https from pri"
  }

  rule {
    protocol    = "TCP"
    port        = 22
    cidr_ip     = "192.168.1.0/24"
    description = "ssh from pri"
  }

  rule {
    protocol    = "ICMP"
    cidr_ip     = "192.168.1.0/24"
    description = "icmp from pri"
  }
}
