terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_ess_domain_identity" "example" {
  domain = "example.com"
}
