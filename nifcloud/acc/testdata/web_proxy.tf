provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_router" "basic" {
  name              = "%s"
  availability_zone = "east-21"

  network_interface {
    network_id = "net-COMMON_GLOBAL"
    dhcp       = false
  }

  network_interface {
    network_id = nifcloud_private_lan.basic.id
  }

}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "%s"
  availability_zone = "east-21"
  cidr_block        = "192.168.1.0/24"
}

resource "nifcloud_web_proxy" "basic" {
  router_id                   = nifcloud_router.basic.id
  bypass_interface_network_id = nifcloud_private_lan.basic.id
  listen_interface_network_id = "net-COMMON_GLOBAL"
  listen_port                 = "8080"
  description                 = "memo"
  name_server                 = "1.1.1.1"
}
