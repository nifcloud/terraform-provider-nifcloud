---
page_title: "NIFCLOUD: nifcloud_devops_instance"
subcategory: "DevOps with GitLab"
description: |-
  Provides a DevOps instance resource.
---

# nifcloud_devops_instance

Provides a DevOps instance resource.

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

resource "nifcloud_devops_instance" "example" {
  instance_id           = "example"
  instance_type         = "c-large"
  firewall_group_name   = nifcloud_devops_firewall_group.example.name
  parameter_group_name  = nifcloud_devops_parameter_group.example.name
  disk_size             = 100
  availability_zone     = "east-11"
  initial_root_password = "initialroo00ootpassword"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The availability zone for the DevOps instance.
* `container_registry_bucket_name` - (Optional) The name of the bucket to put container registry objects.
* `description` - (Optional) Description of the DevOps instance.
* `disk_size` - (Required) The allocated storage in gigabytes.
* `firewall_group_name` - (Required) The name of the DevOps firewall group to associate.
* `initial_root_password` - (Required) Initial password for the root user.
* `instance_id` - (Required) The name of the DevOps instance.
* `instance_type` - (Required) The instance type of the DevOps instance.
* `lfs_bucket_name` - (Optional) The name of the bucket to put LFS objects.
* `network_id` - (Optional, but required if `private_address` is provided) The ID of private lan.
* `object_storage_account` - (Optional, but required if `object_storage_region` is provided) The account name of the object storage service.
* `object_storage_region` - (Optional, but required if `object_storage_account` is provided) The region where the bucket exists.
* `packages_bucket_name` - (Optional) The name of the bucket to put packages.
* `parameter_group_name` - (Required) The name of the DevOps parameter group to associate.
* `private_address` - (Optional, but required if `network_id` is provided) Private IP address for the DevOps instance.
* `to` - (Optional) Mail address where alerts are sent.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `gitlab_url` - URL for GitLab.
* `registry_url` - URL for GitLab container registry.

## Import

nifcloud_devops_instance can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_devops_instance.example foo
```
