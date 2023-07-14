---
page_title: "NIFCLOUD: nifcloud_remote_access_vpn_gateway"
subcategory: "Network"
description: |-
  Provides a remote access vpn gateway resource.
---

# nifcloud_remote_access_vpn_gateway

Provides a remote access vpn gateway resource.

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

resource "nifcloud_remote_access_vpn_gateway" "basic" {
  name               = "example"
  description        = "memo"
  availability_zone  = "east-11"
  accounting_type    = "2"
  type               = "small"
  pool_network_cidr  = "192.168.2.0/24"
  cipher_suite       = ["AES128-GCM-SHA256"]
  ssl_certificate_id = nifcloud_ssl_certificate.basic.id

  user {
    name        = "user1"
    password    = random_password.password.result
    description = "user1"
  }

  user {
    name        = "user2"
    password    = random_password.password.result
    description = "user2"
  }

  network_interface {
    network_id = nifcloud_private_lan.basic.id
    ip_address = "192.168.1.1"
  }

}

resource "nifcloud_private_lan" "basic" {
  availability_zone = "east-11"
  cidr_block        = "192.168.1.0/24"
}

resource "random_password" "password" {
  length = 8
}

resource "tls_private_key" "basic" {
  algorithm = "RSA"
}

resource "tls_self_signed_cert" "basic" {
  private_key_pem       = tls_private_key.basic.private_key_pem
  validity_period_hours = 3
  dns_names             = ["example.com"]
  allowed_uses          = ["client_auth"]

  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }
}

resource "nifcloud_ssl_certificate" "basic" {
  certificate = tls_self_signed_cert.basic.cert_pem
  key         = tls_private_key.basic.private_key_pem
}
```


## Argument Reference

The following arguments are supported:

* `accounting_type` - (Optional) Accounting type. (1: monthly, 2: pay per use).
* `availability_zone` - (Optional) The availability zone.
* `description` - (Optional) The remote access vpn gateway description.
* `ca_certificate_id` - (Optional) The ID of ca certificate.
* `cipher_suite` - (Required) he Cipher suite; can be specified one of `AES128-GCM-SHA256` `AES256-GCM-SHA384` `ECDHE-RSA-AES128-GCM-SHA256` `ECDHE-RSA-AES256-GCM-SHA384` .
* `name` - (Optional) The remote access vpn gateway name.
* `network_interface` - (Required) The network interface list. see [network interface](#network-interface).
* `pool_network_cidr` - (Required) The cidr of pool network; can be specified in the range of /16 to /27..
* `ssl_certificate_id` - (Required) The ID of ssl certificate.
* `type` - (Optional) The type of the remote access vpn gateway. Valid types are `small`, `medium`, `large`.
* `user` - (Optional) List of the remote access vpn gateway user. see [user](#user).

### network interface

* `ip_address` - (Required) The IP address of the network interface.
* `network_id` - (Required) The ID of the network to attach private lan network.

### user

* `description` - (Optional) The remote access vpn gateway user description.
* `name` - (Required) The name of remote access vpn gateway user.
* `password` - (Required) The password of remote access vpn gateway user.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `remote access vpn gateway_id` - The unique ID of the remote access vpn gateway.
* `client_config` - The base64 encoding remote access vpn gateway client config.

## Import

nifcloud_remote_access_vpn_gateway can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_remote access vpn gateway.example foo
```

## Notice

This provider does not support upgrading remote access vpn gateway.
If you want to upgrade the remote access vpn gateway, please upgrade from control panel and also you must delete backup of remote access vpn gateway after upgrade.
