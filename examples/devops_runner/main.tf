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

resource "nifcloud_devops_runner" "example" {
  name              = "example"
  instance_type     = "c-small"
  availability_zone = "east-11"
  concurrent        = 10
  description       = "memo"
}
