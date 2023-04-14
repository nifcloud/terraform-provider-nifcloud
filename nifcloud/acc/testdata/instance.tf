provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_instance" "basic" {
  instance_id             = "%s"
  description             = "memo"
  availability_zone       = "east-21"
  accounting_type         = "2"
  disable_api_termination = true
  image_id                = "221"
  instance_type           = "small"
  key_name                = nifcloud_key_pair.basic.key_name
  security_group          = nifcloud_security_group.basic.group_name
  user_data               = "#!/bin/bash"

  depends_on = [nifcloud_key_pair.basic, nifcloud_security_group.basic, nifcloud_elastic_ip.private, nifcloud_elastic_ip.public]

  network_interface {
    network_id = "net-COMMON_PRIVATE"
    ip_address = nifcloud_elastic_ip.private.private_ip
  }

  network_interface {
    network_id = "net-COMMON_GLOBAL"
    ip_address = nifcloud_elastic_ip.public.public_ip
  }

}

resource "nifcloud_elastic_ip" "private" {
  ip_type           = true
  availability_zone = "east-21"
  description       = "tfacc-memo"
}

resource "nifcloud_elastic_ip" "public" {
  ip_type           = false
  availability_zone = "east-21"
  description       = "tfacc-memo"
}

resource "nifcloud_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}

resource "nifcloud_key_pair" "basic" {
  key_name   = "%s"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}
