---
page_title: "NIFCLOUD: nifcloud_instance"
subcategory: "Computing"
description: |-
  Provides a instance resource.
---

# nifcloud_instance

Provides a instance resource.

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
  image_name = "Ubuntu Server 20.04 LTS"
}
```

## Argument Reference

The following arguments are supported:


* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `admin` - (Optional) Admin user for windows os.
* `availability_zone` - (Optional) The availability zone.
* `description` - (Optional) The instance description.
* `disable_api_termination` - (Optional) If true, enables instance termination protection.
* `image_id` - (Required) The os image identifier to use for the instance.
* `instance_id` - (Optional) The instance name.
* `instance_type` - (Optional) The type of instance to start. Updates to this field will trigger a stop/start of the instance.
* `key_name` - (Optional) The key name of the Key Pair to use for the instance; which can be managed using the nifcloud_key_pair resource.
* `license_name` - (Optional) The license name.
* `license_num` - (Optional) The license count.
* `password` - (Optional) Admin password for windows os.
* `security_group` - (Optional) The security group name to associate with; which can be managed using the nifcloud_security_group resource.
* `user_data` - (Optional) The user data to provide when launching the instance.
* `network_interface` - (Required) The network interface list. see [network interface](#network-interface).

### network interface

* `ip_address` - (Optional) The IP address to select from `static` or `elastic IP address` or `static IP address`; Default(null) is DHCP.
* `network_id` - (Optional) The ID of the network to attach; `net-COMMON_GLOBAL` or `net-COMMON_PRIVATE` or `private lan network id` .
* `network_name` - (Optional) The private lan name of the network to attach.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `instance_state` - The state of the instance.
* `private_ip` - The private ip address of instance.
* `public_ip` - The public ip address of instance.
* `unique_id` - The unique ID of instance.


## Import

nifcloud_instance can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_instance.example foo
```
