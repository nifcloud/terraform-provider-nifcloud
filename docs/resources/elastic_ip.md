---
page_title: "NIFCLOUD: nifcloud_elastic_ip"
subcategory: "Computing"
description: |-
  Provides a elastic ip resource.
---

# nifcloud_elastic_ip

Provides a elastic ip resource.

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

resource "nifcloud_elastic_ip" "example" {
  ip_type           = false
  availability_zone = "east-11"
  description       = "memo"
}

```

## Argument Reference

The following arguments are supported:

* `ip_type` - (Required) Choice of the private ip address(true) or public ip address(false).
* `availability_zone` - (Required) The availability zone.
* `description` - (Optional) The key pair description.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `private_ip` - The private ip address.
* `public_ip` - The public ip address.

## Import

nifcloud_elastic_ip can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_elastic_ip.example foo
```
