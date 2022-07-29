---
page_title: "NIFCLOUD: nifcloud_storage_bucket"
subcategory: "Storage"
description: |-
  Provides a object storage service bucket resource.
---

# nifcloud_storage_bucket

Provides a object storage service bucket resource.

## Example Usage

```
export NIFCLOUD_STORAGE_ACCESS_KEY_ID=<your access key for storage service>
export NIFCLOUD_STORAGE_SECRET_ACCESS_KEY=<your secret key for storage service>
```

```hcl
terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

provider "nifcloud" {
  storage_region = "jp-east-1"
}

resource "nifcloud_storage_bucket" "example" {
  bucket = "example"
  policy = file("policy.json")

  versioning {
    enabled = false
  }
}

```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket.
* `policy` - (Optional) A bucket policy JSON document.
* `versioning` - (Optional) A configuration of the bucket versioning state. see [versioning](#versioning)

### versioning

#### Arguments

* `enabled` - (Optional) Enable versioning.

## Import

nifcloud_storage_bucket can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_storage_bucket.example foo
```
