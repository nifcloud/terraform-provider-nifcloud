terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

provider "nifcloud" {
  storage_region = "jp-east-1"
}

resource "nifcloud_storage_bucket" "example" {
  bucket = "example"

  versioning {
    enabled = true
  }

  policy = file("policy.json")
}
