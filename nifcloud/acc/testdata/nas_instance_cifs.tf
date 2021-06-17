provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_nas_instance" "basic" {
  identifier                     = "%s"
  allocated_storage              = 100
  availability_zone              = "east-21"
  description                    = "memo"
  protocol                       = "cifs"
  type                           = 0
  master_username                = "tfacc"
  master_user_password           = "tfaccpass"
  authentication_type            = 0
  nas_security_group_name        = nifcloud_nas_security_group.basic.group_name
}

resource "nifcloud_nas_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}

resource "nifcloud_instance" "ad" {
  instance_id       = "%s"
  description       = "memo"
  availability_zone = "east-21"
  accounting_type   = "2"
  image_id          = "221"
  instance_type     = "small"
  key_name          = nifcloud_key_pair.basic.key_name
  security_group    = nifcloud_security_group.basic.group_name
  user_data         = <<EOT
%s
EOT

  network_interface {
    network_id = "net-COMMON_PRIVATE"
  }

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }
}

resource "nifcloud_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}

resource "nifcloud_security_group_rule" "from_nas" {
  security_group_names = [nifcloud_security_group.basic.group_name]
  type                 = "IN"
  protocol             = "ANY"
  cidr_ip              = nifcloud_nas_instance.basic.private_ip_address
}

resource "nifcloud_key_pair" "basic" {
  key_name   = "%s"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}
