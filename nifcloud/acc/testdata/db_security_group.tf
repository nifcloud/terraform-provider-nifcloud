provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_db_security_group" "basic" {
  group_name        = "%s"
  description       = "memo"
  availability_zone = "east-21"

  rule {
    security_group_name = nifcloud_security_group.basic.group_name
  }

  rule {
    cidr_ip = "0.0.0.0/0"
  }
}

resource "nifcloud_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}