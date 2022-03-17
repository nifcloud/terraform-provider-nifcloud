terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_dns_record" "example" {
  zone_id = nifcloud_dns_zone.example.id
  name    = "test.example.test"
  type    = "A"
  ttl     = "300"
  record  = ["192.168.0.1"]
  comment = "memo"
}

resource "nifcloud_dns_zone" "example" {
  name    = "example.test"
  comment = "memo"
}

