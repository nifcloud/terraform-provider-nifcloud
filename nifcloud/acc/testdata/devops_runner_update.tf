provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_runner" "basic" {
  name              = "%s-upd"
  instance_type     = "e-small"
  availability_zone = "east-14"
  concurrent        = 50
  description       = "tfacc-memo-upd"
}
