provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_elastic_ip" "basic" {
  ip_type           = false
  availability_zone = "east-21"
  description       = "tfacc-memo-upd"
}
