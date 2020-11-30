provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "%s"
  description       = "tfacc-memo"
  availability_zone = "east-21"
  cidr_block        = "192.168.1.0/24"
  accounting_type   = "1"
}
