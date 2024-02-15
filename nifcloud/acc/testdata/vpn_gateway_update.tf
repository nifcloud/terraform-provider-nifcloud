provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_vpn_gateway" "basic" {
  description       = "memo-upd"
  name              = "%supd"
  type      = "medium"
  availability_zone = "east-21"
  accounting_type   = "1"
  network_name = nifcloud_private_lan.basic.private_lan_name
  ip_address = "192.168.3.2"
  security_group = nifcloud_security_group.basic.group_name
  depends_on     = ["nifcloud_route_table.rt"]
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
