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
		before, after := d.GetChange("instance_id")
		removeList, addList := instanceChangeList(before, after)

		if len(removeList) > 0 {
			input := expandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceIDInput(d, removeList)
			req := svc.NiftyDeregisterInstancesFromSeparateInstanceRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed deregistering instances with separate instance rule: %s", err))
			}
		}

		if len(addList) > 0 {
			input := expandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceIDInput(d, addList)
			req := svc.NiftyRegisterInstancesWithSeparateInstanceRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed registering instances with separate instance rule: %s", err))
			}
		}
	}

	if d.HasChange("instance_unique_id") {
		before, after := d.GetChange("instance_unique_id")
		removeList, addList := instanceChangeList(before, after)

		if len(removeList) > 0 {
			input := expandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceUniqueIDInput(d, removeList)
			req := svc.NiftyDeregisterInstancesFromSeparateInstanceRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed deregistering instances with separate instance rule: %s", err))
			}
		}

		if len(addList) > 0 {
			input := expandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceUniqueIDInput(d, addList)
			req := svc.NiftyRegisterInstancesWithSeparateInstanceRuleRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed registering instances with separate instance rule: %s", err))
			}
		}
	}

	return read(ctx, d, meta)
}

func instanceChangeList(before interface{}, after interface{}) ([]string, []string) {
	var addList, removeList []string
	for _, b := range before.([]interface{}) {
		found := false
		for _, a := range after.([]interface{}) {
			if a.(string) == b.(string) {
				found = true
				break
			}
		}
		if !found {
			removeList = append(removeList, b.(string))
		}
	}

	for _, a := range after.([]interface{}) {
		found := false
		for _, b := range before.([]interface{}) {
			if a.(string) == b.(string) {
				found = true
				break
			}
		}
		if !found {
			addList = append(addList, a.(string))
		}
	}

	return removeList, addList
}
