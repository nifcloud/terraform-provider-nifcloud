package privatelan

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyCreatePrivateLanInput(d)

	svc := meta.(*client.Client).Computing

	res, err := svc.NiftyCreatePrivateLan(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating private_lan: %s", err))
	}

	d.SetId(nifcloud.ToString(res.PrivateLan.NetworkId))
	return read(ctx, d, meta)
}
