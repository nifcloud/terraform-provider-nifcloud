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
	req := svc.AllocateAddressRequest(input)

	res, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating elastic ip: %s", err))
	}

	if d.Get("ip_type").(bool) {
		d.SetId(nifcloud.StringValue(res.PrivateIpAddress))
	} else {
		d.SetId(nifcloud.StringValue(res.PublicIp))
	}

	return update(ctx, d, meta)
}
