provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_vpn_gateway" "basic" {
  description       = "memo"
  name              = "%s"
  type      = "small"
  availability_zone = "east-21"
  accounting_type   = "2"
  network_name = nifcloud_private_lan.basic.private_lan_name
  ip_address = "192.168.3.1"
  security_group = nifcloud_security_group.basic.group_name
  route_table_id = nifcloud_route_table.rt.route_table_id
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

resource "nifcloud_route_table" "rt" {
  route {
    cidr_block = "10.0.1.0/24"
    ip_address = "192.168.3.0"
  }

  route {
    cidr_block = "10.0.2.0/24"
    network_id = "net-COMMON_GLOBAL"
  }
}
