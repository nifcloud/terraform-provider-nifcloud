terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_dns_record" "example1" {
  zone_id = nifcloud_dns_zone.example.id
  name    = "test1"
  type    = "A"
  ttl     = 300
  record  = "192.168.0.1"
  comment = "memo"
}

resource "nifcloud_dns_record" "example2" {
  zone_id = nifcloud_dns_zone.example.id
  name    = "test2.example.test"
  type    = "A"
  ttl     = 300
  record  = "192.168.0.2"
  comment = "memo"
}

resource "nifcloud_dns_record" "example3" {
  zone_id = nifcloud_dns_zone.example.id
  name    = "@"
  type    = "A"
  ttl     = 300
  record  = "192.168.0.3"
  comment = "memo"
}

resource "nifcloud_dns_zone" "example" {
  name    = "example.test"
  comment = "memo"
}
