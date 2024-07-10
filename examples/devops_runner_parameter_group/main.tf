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

resource "nifcloud_devops_runner_parameter_group" "example" {
  name        = "example"
  description = "memo"

  docker_image      = "ruby"
  docker_privileged = true
  docker_shm_size   = 300000

  docker_extra_host {
    host_name  = "example.test"
    ip_address = "192.168.1.2"
  }

  docker_volume = ["/user_data:/cache"]
}
