---
page_title: "NIFCLOUD: nifcloud_devops_runner_parameter_group"
subcategory: "DevOps with GitLab"
description: |-
  Provides a DevOps Runner parameter group resource.
---

# nifcloud_devops_runner_parameter_group

Provides a DevOps Runner parameter group resource.

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

resource "nifcloud_devops_runner_parameter_group" "example" {
  name        = "example"
  description = "memo"

  docker_image      = "ruby"
  docker_privileged = true
  docker_shm_size   = 300000

  docker_extra_host {
    host_name  = "example.test"
    ip_address = "192.168.1.2"
  }

  docker_volume = ["/user_data:/cache"]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the DevOps Runner parameter group.
* `docker_disable_cache` - (Optional) The Docker executor has two levels of caching: a global one (like any other executor) and a local cache based on Docker volumes. This configuration flag acts only on the local one which disables the use of automatically created (not mapped to a host directory) cache volumes. In other words, it only prevents creating a container that holds temporary files of builds, it does not disable the cache if the runner is configured in distributed cache mode.
* `docker_disable_entrypoint_overwrite` - (Optional) Disable the image entrypoint overwriting.
* `docker_extra_host` - (Optional) Hosts that should be defined in container environment.
* `docker_image` - (Optional) The image to run jobs with.
* `docker_oom_kill_disable` - (Optional) If an out-of-memory (OOM) error occurs, do not kill processes in a container.
* `docker_privileged` - (Optional) Run all containers with the privileged flag enabled.
* `docker_shm_size` - (Optional) Shared memory size for images (in bytes).
* `docker_tls_verify` - (Optional) Enable or disable TLS verification of connections to Docker daemon. Disabled by default.
* `docker_volume` - (Optional) Additional volumes that should be mounted. Same syntax as the Docker -v flag.
* `name` - (Required) The name of the DevOps Runner parameter group.

## Import

nifcloud_devops_runner_parameter_group can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_devops_runner_parameter_group.example foo
```
