terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}


resource "nifcloud_separate_instance_rule" "example" {
  instance_id        = [nifcloud_instance.web1.instance_id, nifcloud_instance.web2.instance_id]
  availability_zone  = "east-11"
  description        = "test-upd"
  name               = "testacc001"
}
