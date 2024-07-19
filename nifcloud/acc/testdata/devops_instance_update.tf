provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_instance" "basic" {
  instance_id           = "%s"
  instance_type         = "e-large"
  firewall_group_name   = nifcloud_devops_firewall_group.upd.name
  parameter_group_name  = nifcloud_devops_parameter_group.basic.name
  disk_size             = 300
  availability_zone     = "east-14"
  description           = "tfacc-memo-upd"
  initial_root_password = "initialroo00ootpassword"
  to                    = "email-upd@example.com"
}

resource "nifcloud_devops_firewall_group" "basic" {
  name              = "%s"
  availability_zone = "east-14"
}

resource "nifcloud_devops_parameter_group" "basic" {
  name = "%s"
}

resource "nifcloud_devops_firewall_group" "upd" {
  name              = "%s-upd"
  availability_zone = "east-14"
}
