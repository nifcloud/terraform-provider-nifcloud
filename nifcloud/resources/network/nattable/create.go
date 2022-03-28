package nattable

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

	res, err := svc.NiftyCreateNatTable(ctx, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating natTable: %s", err))
	}

	natTableID := res.NatTable.NatTableId
	d.SetId(nifcloud.ToString(natTableID))

	return update(ctx, d, meta)
}
