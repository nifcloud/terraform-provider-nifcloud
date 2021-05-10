provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_nas_security_group" "basic" {
  group_name        = "%s"
  description       = "memo"
  availability_zone = "east-21"

  rule {
    security_group_name = nifcloud_security_group.basic01.group_name
  }

  rule {
    security_group_name = nifcloud_security_group.basic02.group_name
  }

  rule {
    cidr_ip = "192.168.0.1/32"
  }
}

resource "nifcloud_security_group" "basic01" {
  group_name        = "%s01"
  availability_zone = "east-21"
}

resource "nifcloud_security_group" "basic02" {
  group_name        = "%s02"
  availability_zone = "east-21"
}
