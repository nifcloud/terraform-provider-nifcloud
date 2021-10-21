---
page_title: "NIFCLOUD: nifcloud_security_group_rule"
subcategory: "Computing"
description: |-
  Provides a security group rule resource. Represents a single in or out group rule, which can be added to external Security Groups.
---

# nifcloud_security_group_rule

Provides a security group rule resource. Represents a single in or out group rule, which can be added to external Security Groups.

## Example Usage

```hcl
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

resource "nifcloud_security_group" "example" {
  group_name        = "allowtcp"
  availability_zone = "east-11"
}

resource "nifcloud_security_group_rule" "example" {
  security_group_names = [nifcloud_security_group.example.group_name]
  type                 = "IN"
  from_port            = 0
  to_port              = 65535
  protocol             = "TCP"
  cidr_ip              = "0.0.0.0/0"
}

```

## Argument Reference

The following arguments are supported:


* `cidr_ip` - (Optional) The CIDR IP Address. Cannot be specified with `source_security_group_name` .
* `description` - (Optional) The security group rule description.
* `from_port` - (Optional) The start port
* `protocol` - (Optional) The protocol.
* `security_group_names` - (Required) The security group name list to apply this rule.
* `source_security_group_name` - (Optional) The security group name that allow access. Cannot be specified with `cidr_ip` .
* `to_port` - (Optional) The end port
* `type` - (Optional) The type of rule being created. Valid options are IN (Incoming) or OUT (Outgoing).

## Import

Security Group Rules can be imported using the `type` , `protocol` , `from_port` , `to_port` , source/destination (e.g. `cidr_ip` ) and `security_group_name(s)`
separated by underscores ( `_` ). All parts are required.

### Examples

Import an IN rule in security group `example` for TCP port 8000 with an IPv4 destination CIDR of `10.0.3.0/24` :

```
$ terraform import nifcloud_security_group_rule.example IN_TCP_8080_8080_10.0.3.0/24_example
```

Import an IN rule in security group `example1` and `example2` for TCP port 8000 with an IPv4 destination CIDR of `10.0.3.0/24` :

```
$ terraform import nifcloud_security_group_rule.example IN_TCP_8080_8080_10.0.3.0/24_example1_example2
```

Import a rule applicable to all protocols and ports with a security group source:

```
$ terraform import nifcloud_security_group_rule.example IN_ANY_-_-_sourcename_example
```
