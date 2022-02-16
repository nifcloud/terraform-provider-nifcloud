resource "nifcloud_dns_zone" "basic" {
  name    = var.dns_zone_name
  comment = "tfacc-memo"
}

variable "dns_zone_name" {
    description = "test dns zone"
    type        = string
}
