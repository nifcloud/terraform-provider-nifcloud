package zone

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateHostedZoneInput(d)

	svc := meta.(*client.Client).DNS
	res, err := svc.CreateHostedZone(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating hosted zone: %s", err))
	}

	d.SetId(nifcloud.ToString(res.HostedZone.Name))

	return read(ctx, d, meta)
}
