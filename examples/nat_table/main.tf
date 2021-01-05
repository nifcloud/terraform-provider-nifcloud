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

resource "nifcloud_nat_table" "nat" {
  snat {
    rule_number                   = "1"
    description                   = "memo"
    protocol                      = "TCP"
    source_address                = "192.0.2.1"
    source_port                   = 80
    translation_port              = 81
    outbound_interface_network_id = "net-COMMON_PRIVATE"
  }

  dnat {
    rule_number                    = "1"
    description                    = "memo"
    protocol                       = "ALL"
    translation_address            = "192.168.1.1"
    inbound_interface_network_name = nifcloud_private_lan.pri.private_lan_name
  }
}

resource "nifcloud_private_lan" "pri" {
  private_lan_name  = "pri"
  description       = "pri lan 01"
  availability_zone = "east-11"
  cidr_block        = "192.168.1.0/24"
  accounting_type   = "2"
}
