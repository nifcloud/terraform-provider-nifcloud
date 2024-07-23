---
page_title: "NIFCLOUD: nifcloud_devops_backup_rule"
subcategory: "DevOps with GitLab"
description: |-
  Provides a DevOps backup rule resource.
---

# nifcloud_devops_backup_rule

Provides a DevOps backup rule resource.

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

resource "nifcloud_devops_backup_rule" "example" {
  name        = "example"
  instance_id = nifcloud_devops_instance.example.instance_id
  description = "memo"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the DevOps backup rule.
* `instance_id` - (Required) The name of the DevOps instance.
* `name` - (Required) The name of the DevOps backup rule.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `backup_time` - Cron expression for backup time.

## Import

nifcloud_devops_backup_rule can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_devops_backup_rule.example foo
```
