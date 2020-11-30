provider "nifcloud" {
  region = "jp-east-2"
}
resource "nifcloud_security_group" "basic" {
  group_name             = "%supd"
  description            = "memo-upd"
  availability_zone      = "east-21"
  log_limit              = 100000
  revoke_rules_on_delete = true
}
