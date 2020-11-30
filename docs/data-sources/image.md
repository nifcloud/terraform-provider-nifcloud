---
page_title: "NIFCLOUD: nifcloud_image"
subcategory: "Computing"
description: |-
  Use this data source to get the ID of a image for use in nifcloud_instance resources.
---

# data.nifcloud_image

Use this data source to get the ID of a image for use in nifcloud_instance resources.

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

data "nifcloud_image" "ubuntu" {
  image_name = "Ubuntu Server 20.04 LTS"
}
```

## Argument Reference

The following arguments are supported:


* `image_name` - (Required) The name of image.
* `owner` - (Optional) The image owner; valid values: `niftycloud` (standard image) `self` (current account) `other` (other user).

## Attributes Reference

id is set to the ID of the found image.In addition, the following attributes are exported:

* `image_id` - The id of image.
