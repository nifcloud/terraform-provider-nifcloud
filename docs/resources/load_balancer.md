---
page_title: "NIFCLOUD: nifcloud_load_balancer"
subcategory: "Network"
description: |-
  Provides a load_balancer resource.
---

# nifcloud_load_balancer

Provides a load_balancer resource.

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

resource "nifcloud_load_balancer" "l4lb" {
  load_balancer_name = "l4lb"
  instance_port = 80
  load_balancer_port = 80
  accounting_type = "1"
}

```

## Argument Reference

The following arguments are supported:


* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `balancing_type` - (Optional) Balancing type. (1: Round-Robin, 2: Least-Connection).
* `filter` - (Optional) A list of IP address filter for load balancer.
* `filter_type` - (Optional) The filter_type of filter (1: Allow, 2: Deny).
* `health_check_interval` - (Optional) The interval between health checks.
* `health_check_target` - (Optional) The target of the health check. Valid pattern is ${PROTOCOL}:${PORT} or ICMP.
* `healthy_threshold` - (Optional) The number of checks before the instance is declared healthy.
* `instance_port` - (Required) The port on the instance to route to.
* `instances` - (Optional) A list of instance names to place in the load balancer pool.
* `ip_version` - (Optional) The load balancer ip version(v4 or v6).
* `load_balancer_name` - (Required) The name for the load_balancer.
* `load_balancer_port` - (Required) The port to listen on for the load balancer.
* `network_volume` - (Optional) The load balancer max network volume.
* `policy_type` - (Optional) policy type (standard or ats).
* `session_stickiness_policy_enable` - (Optional) The flag of session stickiness policy.
* `session_stickiness_policy_expiration_period` - (Optional) The session stickiness policy expiration period.
* `sorry_page_enable` - (Optional) The flag of sorry page.
* `sorry_page_status_code` - (Optional) The HTTP status code for sorry page.
* `ssl_certificate_id` - (Optional) The id of the SSL certificate you have uploaded to NIFCLOUD.
* `ssl_policy_id` - (Optional) The id of the SSL policy.
* `ssl_policy_name` - (Optional) The name of the SSL policy.
* `unhealthy_threshold` - (Optional) The number of checks before the instance is declared unhealthy.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `dns_name` - dns_name


## Import

nifcloud_load_balancer can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_load_balancer.example foo
```
