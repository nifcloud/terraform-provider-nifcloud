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

resource "nifcloud_load_balancer" "basic" {
  load_balancer_name = "l4lb"
  instance_port = 80
  load_balancer_port = 80
  accounting_type = "1"
}
