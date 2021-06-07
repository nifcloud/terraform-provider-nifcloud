provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_hatoba_firewall_group" "basic" {
  name        = "%supd"
  description = "memo-upd"

  rule {
    from_port   = 443
    cidr_ip     = "0.0.0.0/0"
    description = "HTTPS incomming"
  }

  rule {
    protocol  = "ANY"
    direction = "IN"
    cidr_ip   = "192.168.0.0/24"
  }

  rule {
    protocol  = "ANY"
    direction = "OUT"
    cidr_ip   = "0.0.0.0/0"
  }
}
