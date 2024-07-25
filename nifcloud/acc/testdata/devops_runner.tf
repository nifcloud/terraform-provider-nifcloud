provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_runner" "basic" {
  name              = "%s"
  instance_type     = "c-small"
  availability_zone = "east-14"
  concurrent        = 10
  description       = "tfacc-memo"
}
