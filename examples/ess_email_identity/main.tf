terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_ess_email_identity" "example" {
  email = "email@example.com"
}
