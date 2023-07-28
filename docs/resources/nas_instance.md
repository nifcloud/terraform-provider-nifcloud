---
page_title: "NIFCLOUD: nifcloud_nas_instance"
subcategory: "NAS"
description: |-
  Provides a NAS instance resource.
---

# nifcloud_nas_instance

Provides a NAS instance resource.

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

resource "nifcloud_nas_instance" "example" {
  identifier              = "nas001"
  availability_zone       = "east-11"
  allocated_storage       = 100
  protocol                = "nfs"
  type                    = 0
  nas_security_group_name = nifcloud_nas_security_group.example.group_name
}

resource "nifcloud_nas_security_group" "example" {
  group_name        = "group001"
  availability_zone = "east-11"
}
```

## Argument Reference

The following arguments are supported:

* `allocated_storage` - (Required) The allocated storage in gibibytes.
* `availability_zone` - (Optional) The AZ for the NAS instance.
* `identifier` - (Required) The name of the NAS instance.
* `description` - (Optional) The NAS instance description.
* `nas_security_group_name` - (Optional) The security group name to associate with; which can be managed using the nifcloud_nas_security_group resource.
* `protocol` - (Required) The protocol of the NAS. `nfs` or `cifs`.
* `master_username` - (Require if protocol is CIFS) The master username.
* `master_user_password` - (Require if protocol is CIFS) The password for masater user.
* `authentication_type` - (Optional) The authentication type for CIFS. (0: local auth)
* `no_root_squash` - (Optional) Turn off root squashing.
* `network_id` - (Optional) The id of private lan.
* `private_ip_address` - (Optional) The private IP address of the NAS instance.
* `private_ip_address_subnet_mask` - (Required if `private_ip_address` is defined) The subnet mask of private IP address written in CIDR notation.
* `type` - (Optional) The type of NAS. (0: standard type, 1: high-speed type)

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `public_ip_address` - The public ip address of the NAS instance.

## Import

nifcloud_nas_instance can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_nas_instance.example foo
```
