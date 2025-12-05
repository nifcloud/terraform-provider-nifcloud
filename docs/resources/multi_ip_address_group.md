---
page_title: "NIFCLOUD: nifcloud_multi_ip_address_group"
subcategory: "Computing"
description: |-
  Provides a multi IP address group resource.
---

# nifcloud_multi_ip_address_group

Provides a multi IP address group resource.

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

resource "nifcloud_multi_ip_address_group" "example" {
  name              = "example"
  availability_zone = "east-11"
  ip_address_count  = 3
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) The availability zone.
* `description` - (Optional) The multi IP address group description.
* `ip_address_count` - (Required) The number of IP addresses to use.
* `name` - (Required) The name of the multi IP address group.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `default_gateway` - The default gateway of this multi IP address group network.
* `ip_addresses` - The list of IP addresses that can be used in this multi IP address group.
* `subnet_mask` - The subnet mask of this multi IP address group network.

## Import

nifcloud_multi_ip_address_group can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_multi_ip_address_group.example foo
```
