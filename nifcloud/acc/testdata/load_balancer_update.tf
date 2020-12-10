provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_load_balancer" "basic" {
  accounting_type = "2"
  availability_zones = ["east-21"]
  ip_version = "v4"
  load_balancer_name = "%s"
  network_volume = 20
  policy_type = "standard"
}
