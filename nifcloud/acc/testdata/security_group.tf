provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_security_group" "basic" {
  group_name             = "%s"
  description            = "memo"
  availability_zone      = "east-21"
  log_limit              = 1000
  revoke_rules_on_delete = false
}
