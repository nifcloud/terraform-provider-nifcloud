provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_multi_ip_address_group" "basic" {
  name              = "%supd"
  description       = "tfacc-memo-upd"
  availability_zone = "east-21"
  ip_address_count  = 3
}
