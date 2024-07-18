---
page_title: "NIFCLOUD: nifcloud_devops_runner"
subcategory: "DevOps with GitLab"
description: |-
  Provides a DevOps Runner resource.
---

# nifcloud_devops_runner

Provides a DevOps Runner resource.

## Example Usage

```hcl
resource "nifcloud_devops_runner" "example" {
  name              = "example"
  instance_type     = "c-small"
  availability_zone = "east-11"
  concurrent        = 10
  description       = "memo"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The availability zone for the DevOps Runner.
* `concurrent` - (Optional) Limits how many jobs can run concurrently, across all registrations.
* `description` - (Optional) Description of the DevOps Runner.
* `instance_type` - (Required) The instance type of the DevOps Runner.
* `name` - (Required) The name of the DevOps Runner.
* `network_id` - (Optional, but required if `private_address` is provided) The ID of private lan.
* `private_address` - (Optional, but required if `network_id` is provided) Private IP address for the DevOps Runner.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `public_ip_address` - Public IP address for the DevOps Runner.
* `system_id` - GitLab Runner system ID.

## Import

nifcloud_devops_runner can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_devops_runner.example foo
```
