---
layout: "nifcloud"
page_title: "NIFCLOUD Provider"
sidebar_current: "docs-nifcloud-index"
description: |-
  The Terraform NIFCLOUD provider is a plugin for Terraform that allows for lifecycle management of NIFCLOUD resources.
---

# NIFCLOUD Provider

The NIFCLOUD provider is used to interact with the resources supported by
NIFCLOUD. The provider needs to be configured with the NIFCLOUD credentials before
it can be used.

You can set environment variable `NIFCLLOUD_ACCESS_KEY_ID` and `NIFCLOUD_SECRET_ACCESS_KEY`

Use the navigation to the left to read about the available resources.

## Example Usage

Example [provider configuration](https://www.terraform.io/docs/configuration/providers.html) in `main.tf` file:

```hcl
provider nifcloud {
  region     = "jp-east-1"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

Example provider configuration using `environment variables`:

```sh
export NIFCLOUD_DEFAULT_REGION=jp-east-1
export NIFCLOUD_ACCESS_KEY_ID=my-access-key
export NIFCLOUD_SECRET_ACCESS_KEY=my-secret-key
```

## Argument Reference

The NIFCLOUD provider requires a few basic parameters:

- `region` - (Required) This is the NIFCLOUD region. It must be provided, but it can also be sourced from the `NIFCLOUD_DEFAULT_REGION` environment variable.
- `access_key` - (Required) This is the NIFCLOUD access key. It must be provided, but it can also be sourced from the `NIFCLOUD_ACCESS_KEY_ID` environment variable.
- `secret_key` - (Required) This is the NIFCLOUD secret key. It must be provided, but it can also be sourced from the `NIFCLOUD_SECRET_ACCESS_KEY` environment variable.
