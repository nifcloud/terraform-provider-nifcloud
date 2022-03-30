---
page_title: "NIFCLOUD: nifcloud_ess_email_identity"
subcategory: "ESS"
description: |-
  Provides an ESS email identity resource
---

# nifcloud_ess_email_identity

Provides an ESS email identity resource

## Example Usage

```hcl
terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_ess_email_identity" "example" {
  email = "email@example.com"
}

```

## Argument Reference

The following arguments are supported:


* `email` - (Required) The email address to assign to ESS


## Import

nifcloud_ess_email_identity can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_ess_email_identity.example foo
```
