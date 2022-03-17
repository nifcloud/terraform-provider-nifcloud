resource "nifcloud_dns_record" "basic" {
  zone_id        = nifcloud_dns_zone.basic.id
  name           = var.dns_record_name
  type           = "A"
  ttl            = 60
  record         = "192.0.2.1"
  comment        = "tfacc-memo"
  failover_routing_policy {
    type = "PRIMARY"
    health_check {
      protocol = "HTTPS"
      ip_address = "192.0.2.2"
      port = 443
      resource_path = "test"
      resource_domain = "example.test"
    }
  }
}

resource "nifcloud_dns_zone" "basic" {
  name    = var.dns_zone_name
  comment = "tfacc-memo"
}

variable "dns_record_name" {
    description = "test dns record"
    type        = string
}

variable "dns_zone_name" {
    description = "test dns zone"
    type        = string
}
