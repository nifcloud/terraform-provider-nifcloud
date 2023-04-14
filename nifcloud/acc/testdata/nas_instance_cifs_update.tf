provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_nas_instance" "basic" {
  identifier              = "%supd"
  allocated_storage       = 200
  availability_zone       = "east-21"
  description             = "memo-upd"
  protocol                = "cifs"
  type                    = 0
  master_username         = "tfacc"
  master_user_password    = "tfaccpass"
  authentication_type     = 0
  nas_security_group_name = nifcloud_nas_security_group.basic.group_name
}

resource "nifcloud_nas_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}
