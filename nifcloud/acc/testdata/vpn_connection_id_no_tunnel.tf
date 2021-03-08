provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_vpn_connection" "basic" {
  type                                                 = "IPsec"
  vpn_gateway_id                                       = nifcloud_vpn_gateway.basic.id
  customer_gateway_id                                  = nifcloud_customer_gateway.basic.id
  ipsec_config_encryption_algorithm                    = "AES128"
  ipsec_config_hash_algorithm                          = "SHA1"
  ipsec_config_pre_shared_key                          = "test"
  ipsec_config_internet_key_exchange                   = "IKEv2"
  ipsec_config_internet_key_exchange_lifetime          = 300
  ipsec_config_encapsulating_security_payload_lifetime = 301
  ipsec_config_diffie_hellman_group                    = 2
  description                                          = "tfacc-memo"
}

resource "nifcloud_customer_gateway" "basic" {
  name                = "%s"
  ip_address          = "192.0.2.1"
  lan_side_ip_address = "192.168.100.10"
  lan_side_cidr_block = "192.168.100.0/28"
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
