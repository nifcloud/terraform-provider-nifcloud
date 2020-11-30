provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_security_group_rule" "basic_cidr" {
  security_group_names = [nifcloud_security_group.fw2.group_name]
  type                 = "OUT"
  protocol             = "ANY"
  cidr_ip              = "0.0.0.0/0"
  description          = "memo"
}

resource "nifcloud_security_group_rule" "basic_source" {
  security_group_names       = [nifcloud_security_group.fw4.group_name]
  type                       = "IN"
  from_port                  = 1
  to_port                    = 65535
  protocol                   = "TCP"
  source_security_group_name = nifcloud_security_group.fw5.group_name
  description                = "memo"

  depends_on = [nifcloud_security_group.fw5]
}

resource "nifcloud_security_group" "fw1" {
  group_name             = "%s"
  availability_zone      = "east-21"
  revoke_rules_on_delete = true
}

resource "nifcloud_security_group" "fw2" {
  group_name             = "%s"
  availability_zone      = "east-21"
  revoke_rules_on_delete = true
}

resource "nifcloud_security_group" "fw3" {
  group_name             = "%s"
  availability_zone      = "east-21"
  revoke_rules_on_delete = true
}

resource "nifcloud_security_group" "fw4" {
  group_name             = "%s"
  availability_zone      = "east-21"
  revoke_rules_on_delete = true

  depends_on = [nifcloud_security_group.fw5]
}

resource "nifcloud_security_group" "fw5" {
  group_name             = "%s"
  availability_zone      = "east-21"
  revoke_rules_on_delete = true
}
