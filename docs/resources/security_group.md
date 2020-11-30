---
page_title: "NIFCLOUD: nifcloud_security_group"
subcategory: "Computing"
description: |-
  Provides a security group resource.
---

# nifcloud_security_group

Provides a security group resource.

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

resource "nifcloud_security_group" "allow_tls" {
  group_name        = "allowtls"
  description       = "Allow TLS inbound traffic"
  availability_zone = "east-11"
}

```

## Argument Reference

The following arguments are supported:


* `availability_zone` - (Required) The availability zone.
* `description` - (Optional) The security group description.
* `group_name` - (Required) The name for the security group.
* `log_limit` - (Optional) The number of log data for security group.
* `revoke_rules_on_delete` - (Optional) Instruct Terraform to revoke all of the Security Groups attached In and Out rules before deleting the rule itself.

## Import

nifcloud_security_group can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_security_group.example foo
```
