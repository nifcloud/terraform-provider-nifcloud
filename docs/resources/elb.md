---
page_title: "NIFCLOUD: nifcloud_elb"
subcategory: "Network"
description: |-
  Provides a multi load balancer resource.
---

# nifcloud_elb

Provides a multi load balancer resource.

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

resource "nifcloud_elb" "bar" {
  elb_name          = "foobar"
  availability_zone = "east-11"
  instance_port     = 80
  protocol          = "HTTP"
  lb_port           = 80

  network_interface {
    network_id     = "net-COMMON_GLOBAL"
    is_vip_network = true
  }

}

```

## Argument Reference

The following arguments are supported:


* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `availability_zone` - (Required) The availability zone.
* `balancing_type` - (Optional) Balancing type. (1: Round-Robin, 2: Least-Connection).
* `description` - (Optional) The multi load balancer description.
* `elb_name` - (Optional) The name for the multi load balancer.
* `health_check_expectation_http_code` - (Optional) A list of the expected http code.
* `health_check_interval` - (Optional) The interval between health checks.
* `health_check_path` - (Optional) The path of the health check.
* `health_check_target` - (Optional) The target of the health check. Valid pattern is ${PROTOCOL}:${PORT} or ICMP.
* `instance_port` - (Required) The port on the instance to route to.
* `instances` - (Optional) A list of instance names to place in the multi load balancer pool.
* `lb_port` - (Required) The port to listen on for the multi load balancer.
* `network_interface` - (Required) The network interface list. see [network interface](#network-interface).
* `network_volume` - (Optional) Maximum network volume for the multi load balancer.
* `protocol` - (Required) The protocol to listen on. Valid values are `HTTP` `HTTPS` `TCP` `UDP`.
* `route_table_id` - (Optional) The id of route table to attach.
* `session_stickiness_policy_enable` - (Optional) The flag of session stickiness policy.
* `session_stickiness_policy_expiration_period` - (Optional) The session stickiness policy expiration period.
* `session_stickiness_policy_method` - (Optional) The session stickiness policy method. (1: Source ip, 2: Cookie)
* `sorry_page_enable` - (Optional) The flag of sorry page.
* `sorry_page_redirect_url` - (Optional) The sorry page redirect url.
* `ssl_certificate_id` - (Optional) The id of the SSL certificate you have uploaded to NIFCLOUD.
* `unhealthy_threshold` - (Optional) The number of checks before the instance is declared unhealthy.

### network_interface

#### Arguments

* `ip_address` - (Optional) The IP address of multi load balancer.
* `is_vip_network` - (Optional) The flag of vip network.
* `network_id` - (Optional) The ID of the network to attach; `net-COMMON_GLOBAL` or `net-COMMON_PRIVATE` or private lan network id.
* `network_name` - (Optional) The private lan name of the network to attach.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `dns_name` - The ip address of multi load balancer vip network.
* `elb_id` - The id of multi load balancer.
* `route_table_association_id` - The id of route table association.
* `version` - The version of multi load balancer.

## Import

nifcloud_elb can be imported using the `elb_id`, `protocol` , `lb_port` , `instance_port`.  
separated by underscores ( `_` ). All parts are required.

### Example

```
$ terraform import nifcloud_elb.example foo_TCP_8080_8080
```
