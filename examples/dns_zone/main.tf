terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_dns_zone" "example" {
  name    = "example.test"
  comment = "memo"
}
