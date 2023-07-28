provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_remote_access_vpn_gateway" "basic" {
  name               = "%s"
  description        = "memo"
  availability_zone  = "east-21"
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
  private_lan_name  = "%s"
  availability_zone = "east-21"
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

resource "nifcloud_ssl_certificate" "upd" {
  certificate = tls_self_signed_cert.basic.cert_pem
  key         = tls_private_key.basic.private_key_pem
}
