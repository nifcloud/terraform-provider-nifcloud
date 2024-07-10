provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_runner_parameter_group" "basic" {
  name        = "%s-upd"
  description = "tfacc-memo-upd"

  docker_image      = "ruby:3"
  docker_privileged = false
  docker_shm_size   = 600000

  docker_extra_host {
    host_name  = "example.test"
    ip_address = "192.168.1.3"
  }

  docker_volume = ["/user_data"]
}
