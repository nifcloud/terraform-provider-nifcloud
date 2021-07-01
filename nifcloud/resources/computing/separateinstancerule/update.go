package separateinstancerule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("description") {
		input := expandNiftyUpdateSeparateInstanceRuleInputForDescription(d)

		req := svc.NiftyUpdateSeparateInstanceRuleRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating separate instance rule description: %s", err))
		}
	}

	if d.HasChange("name") {
		input := expandNiftyUpdateSeparateInstanceRuleInputForName(d)

		req := svc.NiftyUpdateSeparateInstanceRuleRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating separate instance rule name %s", err))
		}

		d.SetId(d.Get("name").(string))

	}

	if d.HasChange("instance_id") {
		o, n := d.GetChange("instance_id")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		addInstances := ns.Difference(os).List()
		delInstances := os.Difference(ns).List()

		if len(delInstances) > 0 {
			input := expandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceIDInput(d, delInstances)

			req := svc.NiftyDeregisterInstancesFromSeparateInstanceRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed deregistering instances with separate instance rule: %s", err))
			}
		}

		if len(addInstances) > 0 {
			input := expandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceIDInput(d, addInstances)

			req := svc.NiftyRegisterInstancesWithSeparateInstanceRuleRequest(input)
			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed registering instances with separate instance rule: %s", err))
			}
		}
	}

	if d.HasChange("instance_unique_id") {
		o, n := d.GetChange("instance_unique_id")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		addInstances := ns.Difference(os).List()
		delInstances := os.Difference(ns).List()

		if len(delInstances) > 0 {
			input := expandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceUniqueIDInput(d, delInstances)

			req := svc.NiftyDeregisterInstancesFromSeparateInstanceRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed deregistering instances with separate instance rule: %s", err))
			}
		}

		if len(addInstances) > 0 {
			input := expandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceUniqueIDInput(d, addInstances)

			req := svc.NiftyRegisterInstancesWithSeparateInstanceRuleRequest(input)
			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed registering instances with separate instance rule: %s", err))
			}
		}
	}

	return read(ctx, d, meta)
}
