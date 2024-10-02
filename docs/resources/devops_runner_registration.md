---
page_title: "NIFCLOUD: nifcloud_devops_runner_registration"
subcategory: "DevOps with GitLab"
description: |-
  Provides a DevOps Runner registration resource.
---

# nifcloud_devops_runner_registration

Provides a DevOps Runner registration resource.

## Example Usage

```hcl
resource "nifcloud_devops_runner_registration" "example" {
  runner_name          = nifcloud_devops_runner.example.name
  gitlab_url           = "https://gitlab.com/"
  parameter_group_name = nifcloud_devops_runner_parameter_group.example.name
  token                = "glrt-thegitlabrunnertoken"
}
```

## Argument Reference

The following arguments are supported:

* `gitlab_url` - (Required) GitLab URL.
* `parameter_group_name` - (Required) The name of the DevOps Runner parameter group to associate.
* `runner_name` - (Required) The name of the DevOps Runner.
* `token` - (Required) GitLab Runner token.

## Attribute Reference

* `id` - ID of the registration.

## Import

nifcloud_devops_runner_registration can be imported using the `runner_name` and `id`, separated by an underscore ( _ ). All parts are required.

```
$ terraform import nifcloud_devops_runner_registration.example foo_foo-abcde
```
