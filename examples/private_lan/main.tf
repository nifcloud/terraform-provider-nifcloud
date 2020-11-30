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

resource "nifcloud_private_lan" "pri" {
  private_lan_name  = "pri"
  description       = "pri lan 01"
  availability_zone = "east-11"
  cidr_block        = "192.168.1.0/24"
  accounting_type   = "1"
}
