package elasticip

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandAllocateAddressInput(d)

	svc := meta.(*client.Client).Computing
	res, err := svc.AllocateAddress(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating elastic ip: %s", err))
	}

	if d.Get("ip_type").(bool) {
		d.SetId(nifcloud.ToString(res.PrivateIpAddress))
	} else {
		d.SetId(nifcloud.ToString(res.PublicIp))
	}

	return update(ctx, d, meta)
}
