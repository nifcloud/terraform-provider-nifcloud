provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_nas_instance" "basic" {
  identifier              = "%s"
  allocated_storage       = 100
  availability_zone       = "east-21"
  description             = "memo"
  protocol                = "nfs"
  type                    = 0
  no_root_squash          = false
  nas_security_group_name = nifcloud_nas_security_group.basic.group_name
}

resource "nifcloud_nas_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}
