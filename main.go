package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: nifcloud.Provider,
	})
}
