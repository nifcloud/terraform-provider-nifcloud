package dbparametergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

const maxModifyParams = 20

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	if d.HasChange("parameter") {
		o, n := d.GetChange("parameter")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		parametersToModify := getParametersToModify(os, ns)
		if len(parametersToModify) != 0 {
			for parametersToModify != nil {
				var targetParams []types.RequestParameters
				if len(parametersToModify) <= maxModifyParams {
					targetParams, parametersToModify = parametersToModify[:], nil
				} else {
					targetParams, parametersToModify = parametersToModify[:maxModifyParams], parametersToModify[maxModifyParams:]
				}

				input := expandModifyDBParameterGroupInput(d, targetParams)

				if _, err := svc.ModifyDBParameterGroup(ctx, input); err != nil {
					return diag.FromErr(fmt.Errorf("failed modifying DBParameterGroup: %s", err))
				}
			}
		}

		parametersToReset := getParametersToReset(os, ns)
		if len(parametersToReset) != 0 {
			for parametersToReset != nil {
				var targetParams []types.RequestParametersOfResetDBParameterGroup
				if len(parametersToReset) <= maxModifyParams {
					targetParams, parametersToReset = parametersToReset[:], nil
				} else {
					targetParams, parametersToReset = parametersToReset[:maxModifyParams], parametersToReset[maxModifyParams:]
				}

				input := expandResetDBParameterGroupInput(d, targetParams)

				if _, err := svc.ResetDBParameterGroup(ctx, input); err != nil {
					return diag.FromErr(fmt.Errorf("failed resetting DBParameterGroup: %s", err))
				}
			}
		}
	}

	return read(ctx, d, meta)
}

func getParametersToModify(old, new *schema.Set) []types.RequestParameters {
	return expandModifyDBParameterGroupParameters(new.Difference(old).List())
}

func getParametersToReset(old, new *schema.Set) []types.RequestParametersOfResetDBParameterGroup {
	toReset := map[string]types.RequestParametersOfResetDBParameterGroup{}
	for _, p := range expandResetDBParameterGroupParameters(old.List()) {
		if p.ParameterName != nil {
			toReset[*p.ParameterName] = p
		}
	}
	for _, p := range expandResetDBParameterGroupParameters(new.List()) {
		if p.ParameterName != nil {
			delete(toReset, *p.ParameterName)
		}
	}

	var toResetParameters []types.RequestParametersOfResetDBParameterGroup
	for _, v := range toReset {
		toResetParameters = append(toResetParameters, v)
	}

	return toResetParameters
}
