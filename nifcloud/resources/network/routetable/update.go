package routetable

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("route") {
		o, n := d.GetChange("route")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old routes and delete any obsolete ones
		for _, route := range ors.List() {
			input := expandDeleteRouteInput(d, route.(map[string]interface{}))

			_, err := svc.DeleteRoute(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating route table to delete route: %s", err))
			}
		}

		// Make sure we save the state of the currently configured rules
		routes := o.(*schema.Set).Intersection(n.(*schema.Set))
		if err := d.Set("route", routes); err != nil {
			return diag.FromErr(err)
		}

		// Then loop through all the newly configured routes and create them
		for _, route := range nrs.List() {
			input := expandCreateRouteInput(d, route.(map[string]interface{}))

			_, err := svc.CreateRoute(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating route table to create route: %s", err))
			}

			routes.Add(route)
			if err := d.Set("route", routes); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return read(ctx, d, meta)
}
