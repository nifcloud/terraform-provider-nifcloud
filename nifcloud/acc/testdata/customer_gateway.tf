provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_customer_gateway" "basic" {
  type                               = "IPsec"
  ip_address                         = "192.168.0.1"
  bgp_asn                            = 65000
  nifty_customer_gateway_name        = "cgw001"
  nifty_lan_side_ip_address          = "192.168.0.1"
  nifty_lan_side_cidr_block          = "192.168.0.0/28"
  nifty_customer_gateway_description = "memo"
}
