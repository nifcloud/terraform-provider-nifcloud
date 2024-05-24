---
page_title: "NIFCLOUD: nifcloud_devops_parameter_group"
subcategory: "DevOps with GitLab"
description: |-
  Provides a DevOps parameter group resource.
---

# nifcloud_devops_parameter_group

Provides a DevOps parameter group resource.

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

resource "nifcloud_devops_parameter_group" "example" {
  name        = "example"
  description = "memo"

  sensitive_parameter {
    name  = "smtp_password"
    value = "mystrongpassword"
  }

  parameter {
    name  = "smtp_user_name"
    value = "user1"
  }

  parameter {
    name  = "gitlab_email_from"
    value = "from@mail.com"
  }

  parameter {
    name  = "gitlab_email_reply_to"
    value = "reply-to@mail.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the DevOps parameter group.
* `name` - (Required) The name of the DevOps parameter group.
* `parameter` - (Optional) List of the DevOps parameters. see [parameter](#parameter)
* `sensitive_parameter` - (Optional) List of the DevOps parameters whose value should be masked. see [sensitive_parameter](#sensitive_parameter)

### parameter

* `name` - (Required) The name of the parameter.
* `value` - (Required) The value of the parameter.

### sensitive_parameter

* `name` - (Required) The name of the parameter. Valid value is `smtp_password`.
* `value` - (Required) The value of the parameter.

Note: The calculated difference in the values of this parameter may be inaccurate due to API restrictions.

## Import

nifcloud_devops_parameter_group can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_devops_parameter_group.example foo
```
