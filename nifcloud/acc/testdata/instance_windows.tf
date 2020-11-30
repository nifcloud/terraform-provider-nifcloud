provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_instance" "basic" {
  instance_id  = "%s"
  image_id     = "189"
  admin        = "testadmin"
  password     = "testpassword"
  license_name = "RDS"
  license_num  = 1

  network_interface {
    network_id = "net-COMMON_PRIVATE"
  }
}
