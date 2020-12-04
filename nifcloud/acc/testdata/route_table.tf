provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_route_table" "basic" {
  route {
    cidr_block = "192.168.1.0/24"
    network_id = nifcloud_private_lan.basic.id
  }

  route {
    cidr_block = "192.168.2.0/24"
    ip_address = "1.1.1.1"
  }
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name = "%s"
  cidr_block       = "192.168.1.0/24"
}
