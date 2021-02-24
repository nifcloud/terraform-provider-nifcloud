provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_vpn_gateway" "basic" {
  nifty_vpn_gateway_description       = "memo-upd"
  nifty_vpn_gateway_name              = "%supd"
  nifty_vpn_gateway_type      = "medium"
  availability_zone = "east-21"
  accounting_type   = "1"
  network_name = nifcloud_private_lan.basic.private_lan_name
  ip_address = "192.168.3.2"
  security_group = nifcloud_security_group.basic.group_name
  route_table_id = ""
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "%s"
  availability_zone = "east-21"
  cidr_block        = "192.168.3.0/24"
}

resource "nifcloud_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}
