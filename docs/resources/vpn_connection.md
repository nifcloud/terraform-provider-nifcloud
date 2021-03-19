---
page_title: "NIFCLOUD: nifcloud_vpn_connection"
subcategory: "Network"
description: |-
  Provides a vpn connection resource.
---

# nifcloud_vpn_connection

Provides a vpn connection resource.

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

resource "nifcloud_vpn_connection" "example" {
  type                                                 = "L2TPv3 / IPsec"
  vpn_gateway_name                                     = nifcloud_vpn_gateway.example.name
  customer_gateway_name                                = nifcloud_customer_gateway.example.name
  tunnel_type                                          = "L2TPv3"
  tunnel_mode                                          = "Unmanaged"
  tunnel_encapsulation                                 = "UDP"
  tunnel_id                                            = "1"
  tunnel_peer_id                                       = "2"
  tunnel_session_id                                    = "1"
  tunnel_peer_session_id                               = "2"
  tunnel_source_port                                   = "7777"
  tunnel_destination_port                              = "7777"
  mtu                                                  = "1000"
  ipsec_config_encryption_algorithm                    = "AES256"
  ipsec_config_hash_algorithm                          = "SHA256"
  ipsec_config_pre_shared_key                          = "test"
  ipsec_config_internet_key_exchange                   = "IKEv2"
  ipsec_config_internet_key_exchange_lifetime          = 300
  ipsec_config_encapsulating_security_payload_lifetime = 300
  ipsec_config_diffie_hellman_group                    = 5
  description                                          = "memo"
}

resource "nifcloud_customer_gateway" "example" {
  name                = "cgw001"
  ip_address          = "192.0.2.1"
  lan_side_ip_address = "192.168.100.10"
}

resource "nifcloud_vpn_gateway" "example" {
  name              = "vpngw001"
  type              = "small"
  availability_zone = "east-11"
  network_name      = nifcloud_private_lan.example.private_lan_name
  ip_address        = "192.168.1.1"
  security_group    = nifcloud_security_group.example.group_name
}

resource "nifcloud_private_lan" "example" {
  private_lan_name  = "example"
  availability_zone = "east-11"
  cidr_block        = "192.168.1.0/24"
}

resource "nifcloud_security_group" "example" {
  group_name        = "example"
  availability_zone = "east-11"
}

```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The type of vpn connection.
* `vpn_gateway_id` - (Optional) The id for the vpn gateway. Cannot be specified with `vpn_gateway_name` .
* `vpn_gateway_name` - (Optional) The name for the vpn gateway. Cannot be specified with `vpn_gateway_id` .
* `customer_gateway_id` - (Optional) The id for the customer gateway. Cannot be specified with `customer_gateway_name` .
* `customer_gateway_name` - (Optional) The name for the customer gateway. Cannot be specified with `customer_gateway_id` .
* `tunnel_type` - (Optional) The type of vpn connection tunnel.
* `tunnel_mode` - (Optional) The mode of vpn connection tunnel; `Unmanaged` or `Managed` .
* `tunnel_encapsulation` - (Optional) The encapsulation of vpn connection.
* `tunnel_id` - (Optional) The id for the vpn gateway tunnel.
* `tunnel_peer_id` - (Optional) The id for the customer gateway tunnel.
* `tunnel_session_id` - (Optional) The session id for the vpn gateway tunnel.
* `tunnel_peer_session_id` - (Optional) The session id for the customer gateway tunnel.
* `tunnel_source_port` - (Optional) The port for the vpn gateway tunnel.
* `tunnel_destination_port` - (Optional) The port for the customer gateway tunnel.
* `mtu` - (Optional) The MTU size for vpn connection.
* `ipsec_config_encryption_algorithm` - (Optional) The encryption algorithm for IPsec config.
  * Specifiable encryption algorithm: [AES128 / AES256 / 3DES]
* `ipsec_config_hash_algorithm` - (Optional) The hash algorithm for IPsec config.
  * Specifiable encryption algorithm: [SHA1 / MD5 / SHA256 / SHA384 / SHA512]
* `ipsec_config_pre_shared_key` - (Optional) The pre shared key for IPsec config.
* `ipsec_config_internet_key_exchange` - (Optional) The IKE protocol for IPsec config.
  * Specifiable encryption algorithm: [IKEv1 / IKEv2]
* `ipsec_config_internet_key_exchange_lifetime` - (Optional) The IKE SA expiration seconds for IPsec config.
* `ipsec_config_encapsulating_security_payload_lifetime` - (Optional) The ESP SA expiration seconds for IPsec config.
* `ipsec_config_diffie_hellman_group` - (Optional) The Diffie-Hellman Group for IKE and PFS.
* `description` - (Optional) The vpn connection description.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `vpn_connection_id` - The id of vpn connection.

## Import

nifcloud_vpn_connection can be imported using the `parameter corresponding to id`, e.g.

```
$ terraform import nifcloud_vpn_connection.example foo
````