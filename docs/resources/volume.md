---
page_title: "NIFCLOUD: nifcloud_volume"
subcategory: "Computing"
description: |-
  Provides a volume resource.
---

# nifcloud_volume

Provides a volume resource.

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

resource "nifcloud_volume" "web" {
  size            = 100
  volume_id       = "volume001"
  disk_type       = "High-Speed Storage A"
  instance_id     = nifcloud_instance.web.instance_id
  reboot          = "true"
  accounting_type = "2"
  description     = "memo"
}

resource "nifcloud_instance" "web" {
  instance_id       = "web001"
  availability_zone = "east-12"
  image_id          = data.nifcloud_image.ubuntu.id
  key_name          = nifcloud_key_pair.web.key_name
  security_group    = nifcloud_security_group.web.group_name
  instance_type     = "small"
  accounting_type   = "2"

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }

  network_interface {
    network_id = "net-COMMON_PRIVATE"
  }
}

resource "nifcloud_key_pair" "web" {
  key_name   = "webkey"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}

resource "nifcloud_security_group" "web" {
  group_name        = "webfw"
  availability_zone = "east-12"
}

data "nifcloud_image" "ubuntu" {
  image_name = "Ubuntu Server 22.04 LTS"
}
```

## Argument Reference

The following arguments are supported:

* `size` - (Required) The disk size.
  * Specifiable size: [100/200/300/.../4000]
  * `disk_type` `Flash Storage` cannot specify more than 1100 size.
* `volume_id` - (Optional) The volume name.
* `disk_type` - (Optional) The disk type. See [disk_type](#disk_type).
* `reboot` - (Optional) The reboot type. See [reboot](#reboot).
* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `description` - (Optional) The volume description.
* `instance_id` - (Optional) The instance name. Cannot be specified with `instance_unique_id`. If you want to change the attached volume, please use this argument.
* `instance_unique_id` - (Optional) The unique ID of instance. Cannot be specified with `instance_id`. This argument is deprecated.

## disk_type

Selectable type:

* Standard Storage
* High-Speed Storage A
* High-Speed Storage B
* Flash Storage
* Standard Flash Storage A
* Standard Flash Storage B
* High-Speed Flash Storage A
* High-Speed Flash Storage B

## reboot

When you want to increase the size of an existing volume, the argument that specifies server restart options.

* `force` - Force restart the instance.
* `true` - (Default) Restart the instance.
* `false` - Do not restart the instance.

## Import

nifcloud_volume can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_volume.example foo
```
