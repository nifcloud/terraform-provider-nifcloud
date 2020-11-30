provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "%supd"
  description       = "tfacc-memo-upd"
  availability_zone = "east-21"
  cidr_block        = "192.168.2.0/24"
  accounting_type   = "2"
}
