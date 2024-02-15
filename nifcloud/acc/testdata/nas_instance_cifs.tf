provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_nas_instance" "basic" {
  identifier              = "%s"
  allocated_storage       = 100
  availability_zone       = "east-21"
  description             = "memo"
  protocol                = "cifs"
  type                    = 0
  master_username         = "tfacc"
  master_user_password    = "tfaccpass"
  nas_security_group_name = nifcloud_nas_security_group.basic.group_name
}

resource "nifcloud_nas_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}
