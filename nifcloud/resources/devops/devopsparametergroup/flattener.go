package devopsparametergroup

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
)

func flatten(d *schema.ResourceData, res *devops.GetParameterGroupOutput) error {
	if res == nil || res.ParameterGroup == nil {
		d.SetId("")
		return nil
	}

	group := res.ParameterGroup

	if nifcloud.ToString(group.ParameterGroupName) != d.Id() {
		return fmt.Errorf("unable to find parameter group within: %#v", group)
	}

	if err := d.Set("name", group.ParameterGroupName); err != nil {
		return err
	}

	if err := d.Set("description", group.Description); err != nil {
		return err
	}

	configParams := d.Get("parameter").(*schema.Set)
	confParams := expandParameters(configParams.List())
	var userParams []types.Parameters
	for _, param := range group.Parameters {
		for _, cp := range confParams {
			if cp.Name == nil {
				continue
			}
			if nifcloud.ToString(cp.Name) == nifcloud.ToString(param.Name) {
				// Skip for sensitive parameters.
				if param.IsSecret != nil && !nifcloud.ToBool(param.IsSecret) {
					userParams = append(userParams, param)
					break
				}
			}
		}
	}

	// Update states of 'parameter' only.
	// The values of 'sensitive_parameter' cannot be fetched from API,
	// therefore the terraform plugin keeps the local values configured in tf files.
	if err := d.Set("parameter", flattenParameters(userParams)); err != nil {
		return err
	}

	return nil
}

func flattenParameters(list []types.Parameters) []map[string]string {
	result := make([]map[string]string, 0, len(list))
	for _, i := range list {
		if i.Name != nil {
			r := make(map[string]string)
			r["name"] = *i.Name

			// Default empty string, guard against nil parameter values
			r["value"] = ""
			if i.Value != nil {
				r["value"] = *i.Value
			}

			result = append(result, r)
		}
	}
	return result
}
