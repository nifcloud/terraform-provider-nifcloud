provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_nat_table" "basic" {
  dnat {
    rule_number                    = "2"
    description                    = "dnat-memo-upd"
    protocol                       = "TCP"
    destination_port               = 80
    translation_address            = "192.168.1.2"
    translation_port               = 81
    inbound_interface_network_name = nifcloud_private_lan.basic.private_lan_name
  }
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name = "%s"
  cidr_block       = "192.168.1.0/24"
}
