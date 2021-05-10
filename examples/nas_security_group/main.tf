terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_nas_security_group" "example" {
  group_name        = "group001"
  description       = "memo"
  availability_zone = "east-11"

  rule {
    cidr_ip = "0.0.0.0/0"
  }

  rule {
    security_group_name = nifcloud_security_group.example.group_name
  }
}

resource "nifcloud_security_group" "example" {
  group_name        = "group001"
  availability_zone = "east-11"
}
