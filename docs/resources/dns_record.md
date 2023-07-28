---
page_title: "NIFCLOUD: nifcloud_dns_record"
subcategory: "DNS"
description: |-
  Provides a DNS record resource.
---

# nifcloud_dns_record

Provides a DNS record resource.

## Example Usage

```hcl
terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_dns_record" "example1" {
  zone_id = nifcloud_dns_zone.example.id
  name    = "test1"
  type    = "A"
  ttl     = 300
  record  = "192.168.0.1"
  comment = "memo"
}

resource "nifcloud_dns_record" "example2" {
  zone_id = nifcloud_dns_zone.example.id
  name    = "test2.example.test"
  type    = "A"
  ttl     = 300
  record  = "192.168.0.2"
  comment = "memo"
}

resource "nifcloud_dns_record" "example3" {
  zone_id = nifcloud_dns_zone.example.id
  name    = "@"
  type    = "A"
  ttl     = 300
  record  = "192.168.0.3"
  comment = "memo"
}

resource "nifcloud_dns_zone" "example" {
  name    = "example.test"
  comment = "memo"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) The ID of the hosted zone to contain this record.
* `name` - (Required) The name of the record.
* `type` - (Required) The type of the record.
* `record` - (Required) The value of the record.
* `ttl` - (Optional) The TTL of the record.
* `weighted_routing_policy` - (Optional) The configs for weighted routing policy. Conflicts with failover_routing_policy. see [weighted_routing_policy](#weighted_routing_policy)
* `failover_routing_policy` - (Optional) The configs for failover routing policy. Conflicts with weighted_routing_policy. see [failover_routing_policy](#failover_routing_policy)
* `comment` - (Optional) The comment of the record.

### weighted_routing_policy

#### Arguments

* `weight` - (Optional) The record weighted value.

### failover_routing_policy

#### Arguments

* `type` - (Optional) The record failover type.
* `health_check` - (Optional) The configs for health check if using failover. see [health_check](#health_check)

### health_check

#### Arguments

* `protocol` - (Optional) The health check protocol.
* `ip_address` - (Optional) The health check IP address.
* `port` - (Optional) The health check port.
* `resource_path` - (Optional) The health check resource path if using HTTP or HTTPS protocol.
* `resource_domain` - (Optional) The health check resource domain if using HTTP or HTTPS protocol.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `set_identifier` - (Optional) The unique identifier to differentiate records with routing policies from one another.

## Import

nifcloud_dns_record can be imported using the `set_identifier`, `zone_id` and `name`.
separated by underscores ( _ ). All parts are required. `name` should be matched with the value of `name` in tf files.

```
$ terraform import nifcloud_dns_record.example XXXXXXXXX_example.test_test.example.test
```
