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

resource "nifcloud_elastic_ip" "example" {
  ip_type           = false
  availability_zone = "east-11"
  description       = "memo"
}
