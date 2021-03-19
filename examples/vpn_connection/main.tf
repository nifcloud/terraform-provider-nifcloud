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