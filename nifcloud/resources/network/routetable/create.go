package routetable

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

	res, err := svc.CreateRouteTable(ctx, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating routeTable: %s", err))
	}

	routeTableID := res.RouteTable.RouteTableId
	d.SetId(nifcloud.ToString(routeTableID))

	return update(ctx, d, meta)
}
