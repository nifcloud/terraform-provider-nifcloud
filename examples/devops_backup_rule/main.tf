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

resource "nifcloud_devops_backup_rule" "example" {
  name        = "example"
  instance_id = nifcloud_devops_instance.example.instance_id
  description = "memo"
}

resource "nifcloud_devops_instance" "example" {
  instance_id           = "example"
  instance_type         = "c-large"
  firewall_group_name   = nifcloud_devops_firewall_group.example.name
  parameter_group_name  = nifcloud_devops_parameter_group.example.name
  disk_size             = 100
  availability_zone     = "east-11"
  initial_root_password = "initialroo00ootpassword"
}

resource "nifcloud_devops_firewall_group" "example" {
  name              = "example"
  availability_zone = "east-11"
}

resource "nifcloud_devops_parameter_group" "example" {
  name = "example"
}
