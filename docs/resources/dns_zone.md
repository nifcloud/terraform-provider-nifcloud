---
page_title: "NIFCLOUD: nifcloud_dns_zone"
subcategory: "DNS"
description: |-
  Provides a DNS zone resource.
---

# nifcloud_dns_zone

Provides a DNS zone resource.

## Example Usage

```hcl
terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_dns_zone" "example" {
  name    = "example.test"
  comment = "memo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the hosted zone.
* `comment` - (Optional) The comment of the hosted zone.

## Attributes Reference

* `name_servers` - The list of name servers.

## Import

nifcloud_dns_zone can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_dns_zone.example foo
```
