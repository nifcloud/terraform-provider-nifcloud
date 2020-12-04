provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_route_table" "basic" {
  route {
    cidr_block   = "192.168.3.0/24"
    network_name = nifcloud_private_lan.basic.private_lan_name
  }
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name = "%s"
  cidr_block       = "192.168.1.0/24"
}
