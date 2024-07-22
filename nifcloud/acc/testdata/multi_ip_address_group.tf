provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_multi_ip_address_group" "basic" {
  name              = "%s"
  description       = "tfacc-memo"
  availability_zone = "east-21"
  ip_address_count  = 1
}
