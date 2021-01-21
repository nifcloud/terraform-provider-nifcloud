provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_customer_gateway" "basic" {
  type                         = "IPsec"
  ip_address                   = "192.168.0.1"
  customer_gateway_name        = "%supd"
  lan_side_ip_address          = "192.168.0.1"
  lan_side_cidr_block          = "192.168.0.0/28"
  customer_gateway_description = "memoupdated"
}
