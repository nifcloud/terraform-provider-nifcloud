terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_network_interface" "example" {
  network_id        = nifcloud_private_lan.dmz.id
  availability_zone = "east-12"

  depends_on = [nifcloud_router.example]
}

resource "nifcloud_instance" "example" {
  image_id = data.nifcloud_image.ubuntu.id
  key_name = nifcloud_key_pair.example.key_name

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }

  network_interface {
    network_id = nifcloud_private_lan.private.id
  }

  network_interface {
    network_interface_id = nifcloud_network_interface.example.id
  }

  depends_on = [nifcloud_router.example]
}

resource "nifcloud_key_pair" "example" {
  key_name   = "examplekey"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}

resource "nifcloud_private_lan" "private" {
  private_lan_name  = "private"
  availability_zone = "east-12"
  cidr_block        = "192.168.1.0/24"
}

resource "nifcloud_private_lan" "dmz" {
  private_lan_name  = "dmz"
  availability_zone = "east-12"
  cidr_block        = "192.168.2.0/24"
}

resource "nifcloud_router" "example" {
  name              = "example"
  availability_zone = "east-12"

  network_interface {
    network_name = nifcloud_private_lan.private.private_lan_name
    dhcp         = true
  }

  network_interface {
    network_name = nifcloud_private_lan.dmz.private_lan_name
    dhcp         = true
  }
}

data "nifcloud_image" "ubuntu" {
  image_name = "Ubuntu Server 22.04 LTS"
}
