package devopsparametergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func updateParameterGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	if d.HasChanges("name", "description", "parameter", "sensitive_parameter") {
		// Merge parameter and sensitive_parameter sets
		op, np := d.GetChange("parameter")
		if op == nil {
			op = new(schema.Set)
		}
		if np == nil {
			np = new(schema.Set)
		}
		ops := op.(*schema.Set)
		nps := np.(*schema.Set)

		osp, nsp := d.GetChange("sensitive_parameter")
		if osp == nil {
			osp = new(schema.Set)
		}
		if nsp == nil {
			nsp = new(schema.Set)
		}
		osps := osp.(*schema.Set)
		nsps := nsp.(*schema.Set)

		os := ops.Union(osps)
		ns := nps.Union(nsps)

		// Fetch parameters to list all parameters
		// since ResourceData contains only configured parameters in tf file.
		getParameterGroupRes, err := svc.GetParameterGroup(ctx, expandGetParameterGroupInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to read a DevOps parameter group: %s", err))
		}
		currentParams := flattenParameters(getParameterGroupRes.ParameterGroup.Parameters)

		parametersToUpdate := getParametersToUpdate(currentParams, os, ns)

		input := expandUpdateParameterGroupInput(d, parametersToUpdate)

		if _, err := svc.UpdateParameterGroup(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update a DevOps parameter group: %s", err))
		}
	}

	if d.HasChange("name") {
		d.SetId(d.Get("name").(string))
	}

	return readParameterGroup(ctx, d, meta)
}

func getParametersToUpdate(params []map[string]string, os, ns *schema.Set) *types.RequestParameters {
	// Create a set of parameter names which has been removed from config.
	toRemove := map[string]struct{}{}
	for _, o := range os.List() {
		oldParam := o.(map[string]interface{})
		toRemove[oldParam["name"].(string)] = struct{}{}
	}
	for _, n := range ns.List() {
		newParam := n.(map[string]interface{})
		delete(toRemove, newParam["name"].(string))
	}

	diffs := ns.Difference(os).List()

	for _, param := range params {
		// If the parameter is marked as removed,
		// set an empty value since it is always allowed.
		if _, ok := toRemove[param["name"]]; ok {
			param["value"] = ""
			continue
		}

		// API returns the masked value when the parameter is sensistive.
		isSensitive := param["value"] == "********"
		found := false

		for _, diff := range diffs {
			diffParam := diff.(map[string]interface{})

			if diffParam["name"] == param["name"] {
				param["value"] = diffParam["value"].(string)
				found = true
				break
			}
		}

		// To avoid updating the sensitive parameter with a masked value,
		// set a previously configued value or an empty value.
		if isSensitive && !found {
			param["value"] = ""
			for _, o := range os.List() {
				oldParam := o.(map[string]interface{})
				if oldParam["name"].(string) == param["name"] {
					param["value"] = oldParam["value"].(string)
					break
				}
			}
		}
	}

	return expandUpdateParameterGroupParameters(params)
}
