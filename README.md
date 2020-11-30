# Terraform Provider for NIFCLOUD

![Build](https://github.com/nifcloud/terraform-provider-nifcloud/workflows/Build/badge.svg?branch=master)

The Terraform NIFCLOUD provider is a plugin for Terraform that allows for lifecycle management of NIFCLOUD resources.

---

## Using the provider

- [Terraform Website](https://terraform.io)
- [NIFCLOUD Provider Documentation](docs/index.md)
- [NIFCLOUD Provider Examples](examples/)

## Usage Example

```hcl
terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
      version = ">= 1.0.0"
    }
  }
}

provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_key_pair" "deployer" {
  key_name   = "deployerkey"
  public_key = "Base64-encoded public key string"
}

resource "nifcloud_instance" "web" {
  image_id = "221"
  key_name = nifcloud_key_pair.deployer.key_name

  network_interface {
    network_id = "net-COMMON_PRIVATE"
  }
}
```

Have a look at the [reference docs](docs/index.md) for more information on the supported resources and data sources.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+