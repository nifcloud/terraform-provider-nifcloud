---
page_title: "NIFCLOUD: nifcloud_elb_listener"
subcategory: "Computing"
description: |-
  Provides a multi load balancer listener resource.
---

# nifcloud_elb_listener

Provides a multi load balancer listener resource.

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

resource "nifcloud_elb" "front_end" {
  elb_name          = "frontend"
  availability_zone = "east-11"
  instance_port     = 8082
  protocol          = "HTTP"
  lb_port           = 80

  network_interface {
    network_id     = "net-COMMON_GLOBAL"
    is_vip_network = true
  }
}

resource "nifcloud_elb_listener" "front_end" {
  elb_id        = nifcloud_elb.front_end.id
  instance_port = 8082
  protocol      = "HTTP"
  lb_port       = 8080
  description   = "memo"
}

```

## Argument Reference

The following arguments are supported:


* `balancing_type` - (Optional) Balancing type. (1: Round-Robin, 2: Least-Connection).
* `description` - (Optional) The multi load balancer description.
* `elb_id` - (Required) The id of multi load balancer.
* `health_check_expectation_http_code` - (Optional) A list of the expected http code.
* `health_check_interval` - (Optional) The interval between health checks.
* `health_check_path` - (Optional) The path of the health check.
* `health_check_target` - (Optional) The target of the health check. Valid pattern is ${PROTOCOL}:${PORT} or ICMP.
* `instance_port` - (Required) The port on the instance to route to.
* `instances` - (Optional) A list of instance names to place in the multi load balancer pool.
* `lb_port` - (Required) The port to listen on for the multi load balancer.
* `protocol` - (Required) The protocol to listen on. Valid values are `HTTP` `HTTPS` `TCP` `UDP`.
* `session_stickiness_policy_enable` - (Optional) The flag of session stickiness policy.
* `session_stickiness_policy_expiration_period` - (Optional) The session stickiness policy expiration period.
* `session_stickiness_policy_method` - (Optional) The session stickiness policy method. (1: Source ip, 2: Cookie)
* `sorry_page_enable` - (Optional) The flag of sorry page.
* `sorry_page_redirect_url` - (Optional) The sorry page redirect url.
* `ssl_certificate_id` - (Optional) The id of the SSL certificate you have uploaded to NIFCLOUD.
* `unhealthy_threshold` - (Optional) The number of checks before the instance is declared unhealthy.

## Import

nifcloud_elb_listener can be imported using the `elb_id`, `protocol` , `lb_port` , `instance_port`.  
separated by underscores ( `_` ). All parts are required.

### Example

```
$ terraform import nifcloud_elb_listener.example foo_TCP_8080_8080
```
