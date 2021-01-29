package nattable

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("snat") {
		o, n := d.GetChange("snat")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old snat rules and delete any obsolete ones
		for _, snat := range ors.List() {
			input := expandNiftyDeleteNatRuleInputForSnat(d, snat.(map[string]interface{}))
			req := svc.NiftyDeleteNatRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating nat table to delete snat rule: %s", err))
			}
		}

		// Make sure we save the state of the currently configured snat rules
		snats := o.(*schema.Set).Intersection(n.(*schema.Set))
		if err := d.Set("snat", snats); err != nil {
			return diag.FromErr(err)
		}

		// Then loop through all the newly configured snat rules and create them
		for _, snat := range nrs.List() {
			input := expandNiftyCreateNatRuleInputForSnat(d, snat.(map[string]interface{}))
			req := svc.NiftyCreateNatRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating nat table to create snat rule: %s", err))
			}

			snats.Add(snat)
			if err := d.Set("snat", snats); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("dnat") {
		o, n := d.GetChange("dnat")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old dnat rules and delete any obsolete ones
		for _, dnat := range ors.List() {
			input := expandNiftyDeleteNatRuleInputForDnat(d, dnat.(map[string]interface{}))
			req := svc.NiftyDeleteNatRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating nat table to delete dnat rule: %s", err))
			}
		}

		// Make sure we save the state of the currently configured dnat rules
		dnats := o.(*schema.Set).Intersection(n.(*schema.Set))
		if err := d.Set("dnat", dnats); err != nil {
			return diag.FromErr(err)
		}

		// Then loop through all the newly configured dnat rules and create them
		for _, dnat := range nrs.List() {
			input := expandNiftyCreateNatRuleInputForDnat(d, dnat.(map[string]interface{}))
			req := svc.NiftyCreateNatRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating nat table to create dnat rule: %s", err))
			}

			dnats.Add(dnat)
			if err := d.Set("dnat", dnats); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return read(ctx, d, meta)
}
