---
page_title: "NIFCLOUD: nifcloud_key_pair"
subcategory: "Computing"
description: |-
  Upload and register the specified SSH public key.
---

# nifcloud_key_pair

Upload and register the specified SSH public key.

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

resource "nifcloud_key_pair" "deployer" {
  key_name    = "deployerkey"
  public_key  = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
  description = "memo"
}

```

## Argument Reference

The following arguments are supported:


* `description` - (Optional) The key pair description.
* `key_name` - (Required) The name for the key pair.
* `public_key` - (Required) The public key material.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `fingerprint` - The MD5 public key fingerprint.


## Import

nifcloud_key_pair can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_key_pair.example foo
```
