---
page_title: "NIFCLOUD: nifcloud_ess_domain_identity"
subcategory: "ESS"
description: |-
  Provides an ESS domain identity resource
---

# nifcloud_ess_domain_identity

Provides an ESS domain identity resource

## Example Usage

```hcl
terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_ess_domain_identity" "example" {
  domain = "example.com"
}
```

## Argument Reference

The following arguments are supported:


* `domain` - (Required) The domain name to assign to ESS


## Import

nifcloud_ess_domain_identity can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_ess_domain_identity.example foo
```
