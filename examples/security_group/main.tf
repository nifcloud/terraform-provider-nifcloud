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

resource "nifcloud_security_group" "allow_tls" {
  group_name        = "allowtls"
  description       = "Allow TLS inbound traffic"
  availability_zone = "east-11"
}
