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

resource "nifcloud_ssl_certificate" "example" {
  certificate = file("${path.module}/certs/certificate.pem")
  key         = file("${path.module}/certs/private_key.pem")
  ca          = file("${path.module}/certs/ca.pem")
  description = "memo"
}
