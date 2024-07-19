provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_firewall_group" "basic" {
  name              = "%s-upd"
  availability_zone = "east-14"
  description       = "tfacc-memo-upd"

  rule {
    protocol = "TCP"
    port     = 443
    cidr_ip  = "192.168.1.0/24"
  }

  rule {
    protocol = "TCP"
    port     = 22
    cidr_ip  = "192.168.1.0/24"
  }

  rule {
    protocol    = "ICMP"
    cidr_ip     = "192.168.1.0/24"
    description = "pong"
  }
}
