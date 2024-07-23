provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_backup_rule" "basic" {
  name        = "%s"
  instance_id = nifcloud_devops_instance.basic.instance_id
  description = "tfacc-memo"
}

resource "nifcloud_devops_instance" "basic" {
  instance_id           = "%s"
  instance_type         = "c-large"
  firewall_group_name   = nifcloud_devops_firewall_group.basic.name
  parameter_group_name  = nifcloud_devops_parameter_group.basic.name
  disk_size             = 100
  availability_zone     = "east-14"
  initial_root_password = "initialroo00ootpassword"
}

resource "nifcloud_devops_firewall_group" "basic" {
  name              = "%s"
  availability_zone = "east-14"
}

resource "nifcloud_devops_parameter_group" "basic" {
  name = "%s"
}