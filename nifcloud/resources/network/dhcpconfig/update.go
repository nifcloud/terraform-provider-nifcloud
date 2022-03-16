package dhcpconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("static_mapping") {
		o, n := d.GetChange("static_mapping")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old static mapping ipaddresses and delete any obsolete ones
		for _, sm := range ors.List() {
			input := expandNiftyDeleteDhcpConfigStaticMappingInput(d, sm.(map[string]interface{}))
			_, err := svc.NiftyDeleteDhcpStaticMapping(ctx, input)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating dhcp config to delete static mapping : %s", err))
			}
		}

		// Make sure we save the state of the currently configured static mapping
		staticmapping := o.(*schema.Set).Intersection(n.(*schema.Set))
		if err := d.Set("static_mapping", staticmapping); err != nil {
			return diag.FromErr(err)
		}

		// Then loop through all the newly configured static mapping and create them
		for _, sm := range nrs.List() {
			input := expandNiftyCreateDhcpConfigStaticMappingInput(d, sm.(map[string]interface{}))
			_, err := svc.NiftyCreateDhcpStaticMapping(ctx, input)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating dhcp config to create static mapping: %s", err))
			}

			staticmapping.Add(sm)
			if err := d.Set("static_mapping", staticmapping); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("ipaddress_pool") {
		o, n := d.GetChange("ipaddress_pool")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old ipaddress pool and delete any obsolete ones
		for _, ip := range ors.List() {
			input := expandNiftyDeleteDhcpConfigIPAddressPoolInput(d, ip.(map[string]interface{}))
			_, err := svc.NiftyDeleteDhcpIpAddressPool(ctx, input)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating dhcp config to delete ipaddress pool: %s", err))
			}
		}

		// Make sure we save the state of the currently configured ipaddress pool
		ipaddresspool := o.(*schema.Set).Intersection(n.(*schema.Set))
		if err := d.Set("ipaddress_pool", ipaddresspool); err != nil {
			return diag.FromErr(err)
		}

		// Then loop through all the newly configured ipaddress pool and create them
		for _, ip := range nrs.List() {
			input := expandNiftyCreateDhcpConfigIPAddressPoolInput(d, ip.(map[string]interface{}))
			_, err := svc.NiftyCreateDhcpIpAddressPool(ctx, input)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating dhcp config to create ipaddress pool: %s", err))
			}

			ipaddresspool.Add(ip)
			if err := d.Set("ipaddress_pool", ipaddresspool); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return read(ctx, d, meta)
}
