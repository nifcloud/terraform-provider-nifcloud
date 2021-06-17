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

resource "nifcloud_nas_instance" "example" {
  identifier              = "nas001"
  availability_zone       = "east-11"
  allocated_storage       = 100
  protocol                = "nfs"
  type                    = 0
  nas_security_group_name = nifcloud_nas_security_group.example.group_name
}

resource "nifcloud_nas_security_group" "example" {
  group_name        = "group001"
  availability_zone = "east-11"
}
