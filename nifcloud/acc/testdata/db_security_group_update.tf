provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_db_security_group" "basic" {
  group_name        = "%supd"
  description       = "memo-upd"
  availability_zone = "east-21"

  rule {
    security_group_name = nifcloud_security_group.basic.group_name
  }

  rule {
    cidr_ip = "192.168.0.1/32"
  }
}

resource "nifcloud_security_group" "basic" {
  group_name        = "%supd"
  availability_zone = "east-21"
}