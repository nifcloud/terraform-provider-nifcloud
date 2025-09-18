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

resource "nifcloud_multi_ip_address_group" "basic" {
  name              = "basic"
  description       = "memo"
  availability_zone = "east-12"
  ip_address_count  = 1
}

resource "nifcloud_key_pair" "web" {
  key_name   = "webkey"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}

resource "nifcloud_security_group" "web" {
  group_name        = "webfw"
  availability_zone = "east-12"
}

data "nifcloud_image" "ubuntu" {
  image_name = "Ubuntu Server 24.04 LTS"
}

resource "nifcloud_instance" "web" {
  instance_id       = "web001"
  availability_zone = "east-12"
  image_id          = data.nifcloud_image.ubuntu.id
  key_name          = nifcloud_key_pair.web.key_name
  security_group    = nifcloud_security_group.web.group_name
  instance_type     = "c2-small"
  accounting_type   = "2"

  network_interface {
    network_id = "net-MULTI_IP_ADDRESS"
    multi_ip_address_group_id = nifcloud_multi_ip_address_group.basic.id
  }

  network_interface {
    network_id = "net-COMMON_PRIVATE"
  }

  multi_ip_address_configuration_user_data = templatefile("scripts/userdata.sh", {
    ip_address = nifcloud_multi_ip_address_group.basic.ip_addresses[0]
    default_gateway = nifcloud_multi_ip_address_group.basic.default_gateway
    subnet_mask = nifcloud_multi_ip_address_group.basic.subnet_mask
  })
}
