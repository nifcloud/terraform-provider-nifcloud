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

resource "nifcloud_customer_gateway" "example" {
  type                         = "IPsec"
  ip_address                   = "192.168.0.1"
  customer_gateway_name        = "cgw002"
  lan_side_ip_address          = "192.168.0.1"
  lan_side_cidr_block          = "192.168.0.0/28"
  customer_gateway_description = "memo"
}
