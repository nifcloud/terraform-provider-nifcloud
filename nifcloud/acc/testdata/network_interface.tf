provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "%s"
  availability_zone = "east-21"
  cidr_block        = "192.168.100.0/24"
}

resource "nifcloud_router" "basic" {
  name              = "%s"
  availability_zone = "east-21"
  type              = "large"

  network_interface {
    network_id = nifcloud_private_lan.basic.id
    dhcp       = true
  }
}

resource "nifcloud_network_interface" "basic" {
  network_id        = nifcloud_private_lan.basic.id
  availability_zone = "east-21"
  description       = "%s"
  ip_address        = "static"

  depends_on = [nifcloud_router.basic, nifcloud_private_lan.basic]
}
