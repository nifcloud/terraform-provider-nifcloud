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

resource "nifcloud_elb" "front_end" {
  elb_name          = "frontend"
  availability_zone = "east-11"
  instance_port     = 8082
  protocol          = "HTTP"
  lb_port           = 80

  network_interface {
    network_id     = "net-COMMON_GLOBAL"
    is_vip_network = true
  }
}

resource "nifcloud_elb_listener" "front_end" {
  elb_id        = nifcloud_elb.front_end.id
  instance_port = 8082
  protocol      = "HTTP"
  lb_port       = 8080
  description   = "memo"
}
