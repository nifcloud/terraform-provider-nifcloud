provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_nat_table" "basic" {
  snat {
    rule_number                   = "1"
    description                   = "snat-memo"
    protocol                      = "TCP"
    source_address                = "192.0.2.1"
    source_port                   = 80
    translation_port              = 81
    outbound_interface_network_id = "net-COMMON_PRIVATE"
  }

  dnat {
    rule_number                  = "1"
    description                  = "dnat-memo"
    protocol                     = "ALL"
    translation_address          = "192.168.1.1"
    inbound_interface_network_id = nifcloud_private_lan.basic.id
  }
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name = "%s"
  cidr_block       = "192.168.1.0/24"
}
