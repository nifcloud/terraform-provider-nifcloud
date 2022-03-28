package dhcpconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	svc := meta.(*client.Client).Computing
	res, err := svc.NiftyCreateDhcpConfig(ctx, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating DhcpConfig: %s", err))
	}

	dhcpConfigID := res.DhcpConfig.DhcpConfigId
	d.SetId(nifcloud.ToString(dhcpConfigID))

	return update(ctx, d, meta)
}
