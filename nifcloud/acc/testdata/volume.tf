provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_volume" "basic" {
  size            = 100
  volume_id       = "%s"
  disk_type       = "High-Speed Storage A"
  instance_id     = nifcloud_instance.basic.instance_id
  reboot          = "true"
  accounting_type = "1"
  description     = "memo"
}

resource "nifcloud_instance" "basic" {
  instance_id             = "%s"
  description             = "memo"
  availability_zone       = "east-21"
  accounting_type         = "2"
  image_id                = data.nifcloud_image.ubuntu.id
  instance_type           = "mini"
  key_name                = nifcloud_key_pair.basic.key_name
  security_group          = nifcloud_security_group.basic.group_name

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

resource "nifcloud_key_pair" "basic" {
  key_name   = "%s"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}

data "nifcloud_image" "ubuntu" {
  image_name = "Ubuntu Server 22.04 LTS"
}
