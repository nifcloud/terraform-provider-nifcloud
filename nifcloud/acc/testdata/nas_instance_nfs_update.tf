provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_nas_instance" "basic" {
  identifier              = "%supd"
  allocated_storage       = 200
  availability_zone       = "east-21"
  description             = "memo-upd"
  protocol                = "nfs"
  type                    = 0
  no_root_squash          = true
  nas_security_group_name = nifcloud_nas_security_group.basic.group_name
}

resource "nifcloud_nas_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}
