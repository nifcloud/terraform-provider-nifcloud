---
page_title: "NIFCLOUD: nifcloud_load_balancer"
subcategory: "Computing"
description: |-
  Provides a load balancer resource.
---

## nifcloud_load_balancer

Provides a load balancer resource.

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
  accounting_type = "1"
  load_balancer_name = "nl4lb"
  load_balancer_port = 80
  instance_port = 80
}
```

## Argument Reference

The following arguments are supported:

* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `balancing_type` - (Optional) Balancing type. (1: Round-Robin, 2: Least-Connection).
* `filter` - (Optional) A list of IP address filter for load balancer.
* `filter_type` - (Optional) The filter_type of filter (1: Allow, 2: Deny).
* healthy_threshold
* `health_check_interval` - (Optional) The interval between health checks.
* `health_check_target` - (Optional) The target of the health check. Valid pattern is ${PROTOCOL}:${PORT} or ICMP.
* `healthy_threshold1` - (Optional) The number of checks before the instance is declared healthy.
* `instances` - (Optional) A list of instance names to place in the multi load balancer pool.
* `instance_port` - 
* `ip_version` - (Optional) The IP version. (v4 or v6).
* `load_balancer_name` - (Required) The load balancer name.
* `network_volume` - (Optional) The network volume.
* `session_stickiness_policy_enable` - (Optional) The flag of session stickiness policy.
* `session_stickiness_policy_expiration_period` - (Optional) The session stickiness policy expiration period.
* `sorry_page_enable` - (Optional) The flag of sorry page.
* `sorry_page_status_code` - (Optional)  The HTTP status code for sorry page.
* `ssl_certificate_id` - (Optional) The id of the SSL certificate you have uploaded to NIFCLOUD.
* `ssl_policy_id` - (Optional) The id of the SSL policy.
* `ssl_policy_name` - (Optional) The name of the SSL policy.
* `policy_type` - (Optional) The policy type. (standard or ats).
* `unhealthy_threshold` - (Optional) The number of checks before the instance is declared unhealthy.

## Import

nifcloud_load_balancer can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_load_balancer.example foo
```
