provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_vpn_connection" "basic" {
  type                                                 = "L2TPv3 / IPsec"
  vpn_gateway_name                                     = nifcloud_vpn_gateway.basic.name
  customer_gateway_name                                = nifcloud_customer_gateway.basic.name
  tunnel_type                                          = "L2TPv3"
  tunnel_mode                                          = "Unmanaged"
  tunnel_encapsulation                                 = "UDP"
  tunnel_id                                            = "1"
  tunnel_peer_id                                       = "2"
  tunnel_session_id                                    = "1"
  tunnel_peer_session_id                               = "2"
  tunnel_source_port                                   = "7777"
  tunnel_destination_port                              = "7778"
  mtu                                                  = "1000"
  ipsec_config_encryption_algorithm                    = "AES256"
  ipsec_config_hash_algorithm                          = "SHA256"
  ipsec_config_pre_shared_key                          = "test"
  ipsec_config_internet_key_exchange                   = "IKEv2"
  ipsec_config_internet_key_exchange_lifetime          = 300
  ipsec_config_encapsulating_security_payload_lifetime = 301
  ipsec_config_diffie_hellman_group                    = 5
  description                                          = "tfacc-memo"
}

resource "nifcloud_customer_gateway" "basic" {
  name                = "%s"
  ip_address          = "192.0.0.1"
  lan_side_ip_address = "192.168.100.10"
}

resource "nifcloud_vpn_gateway" "basic" {
  name              = "%s"
  type              = "small"
  availability_zone = "east-21"
  network_name      = nifcloud_private_lan.basic.private_lan_name
  ip_address        = "192.168.3.1"
  security_group     = nifcloud_security_group.basic.group_name
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "%s"
  availability_zone = "east-21"
  cidr_block        = "192.168.3.0/24"
}

resource "nifcloud_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}