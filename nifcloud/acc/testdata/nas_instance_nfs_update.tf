provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_nas_instance" "basic" {
  identifier                     = "%supd"
  allocated_storage              = 200
  availability_zone              = "east-21"
  description                    = "memo-upd"
  protocol                       = "nfs"
  type                           = 0
  no_root_squash                 = true
  network_id                     = nifcloud_private_lan.basic.id
  private_ip_address             = "192.168.1.101"
  private_ip_address_subnet_mask = "/24"
  nas_security_group_name        = nifcloud_nas_security_group.basic.group_name
}

resource "nifcloud_nas_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "%s"
  availability_zone = "east-21"
  cidr_block        = "192.168.1.0/24"
}
