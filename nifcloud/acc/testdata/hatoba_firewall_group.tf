provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_hatoba_firewall_group" "basic" {
  name        = "%s"
  description = "memo"

  rule {
    protocol    = "TCP"
    direction   = "IN"
    from_port   = 80
    to_port     = 80
    cidr_ip     = "0.0.0.0/0"
    description = "rule memo"
  }

  rule {
    from_port = 443
    cidr_ip   = "0.0.0.0/0"
  }

  rule {
    protocol  = "TCP"
    direction = "OUT"
    from_port = 53
    to_port   = 53
    cidr_ip   = "8.8.8.8"
  }

  rule {
    protocol  = "UDP"
    direction = "OUT"
    from_port = 53
    to_port   = 53
    cidr_ip   = "8.8.8.8"
  }
}
