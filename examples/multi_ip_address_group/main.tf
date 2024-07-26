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

resource "nifcloud_multi_ip_address_group" "basic" {
  name              = "basic"
  description       = "memo"
  availability_zone = "east-11"
  ip_address_count  = 1
}
