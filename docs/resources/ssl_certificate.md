---
page_title: "NIFCLOUD: nifcloud_ssl_certificate"
subcategory: "Computing"
description: |-
  Provides a SSL certificate resource.
---

# nifcloud_security_group

Provides a SSL certificate resource.

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

resource "nifcloud_ssl_certificate" "example" {
  certificate = file("${path.module}/certificate.pem")
  key         = file("${path.module}/private_key.pem")
  ca          = file("${path.module}/ca.pem")
  description = "memo"
}

```

## Argument Reference

The following arguments are supported:

* `certificate` - (Required) The PEM encoded certificate.
* `key` - (Required) The PEM encoded private key.
* `ca` - (Optional) The PEM encoded certificate authority.
* `description` - (Optional) The SSL certificate description.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `fqdn_id` - The identifier of a certificate.
* `fqdn` - The FQDN of a certificate (same as the Common Name).

## Import

nifcloud_ssl_certificate can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_ssl_certificate.example foo
```
