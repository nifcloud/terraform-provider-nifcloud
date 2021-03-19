package dbparametergroup

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
)

func flatten(d *schema.ResourceData, groups *rdb.DescribeDBParameterGroupsResponse, parameters []rdb.Parameter) error {
	if groups == nil || len(groups.DBParameterGroups) == 0 {
		d.SetId("")
		return nil
	}

	if groups == nil {
		return fmt.Errorf("the response of DescribeDBParameters must not be nil")
	}

	group := groups.DBParameterGroups[0]

	if nifcloud.StringValue(group.DBParameterGroupName) != d.Id() {
		return fmt.Errorf("unable to find DB parameter group within: %#v", groups.DBParameterGroups)
	}

	if err := d.Set("name", group.DBParameterGroupName); err != nil {
		return err
	}

	if err := d.Set("family", group.DBParameterGroupFamily); err != nil {
		return err
	}

	if err := d.Set("description", group.Description); err != nil {
		return err
	}

	configParams := d.Get("parameter").(*schema.Set)
	var userParams []rdb.Parameter
	confParams := expandParameters(configParams.List())
	for _, param := range parameters {
		if param.ParameterName == nil {
			continue
		}

		for _, cp := range confParams {
			if cp.ParameterName == nil {
				continue
			}

			if nifcloud.StringValue(cp.ParameterName) == nifcloud.StringValue(param.ParameterName) {
				// override ApplyMethod with config value because RDB API does not return this field.
				param.ApplyMethod = cp.ApplyMethod
				userParams = append(userParams, param)
				break
			}
		}
	}

	if err := d.Set("parameter", flattenParameters(userParams)); err != nil {
		return err
	}

	return nil
}

func flattenParameters(list []rdb.Parameter) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		if i.ParameterName != nil {
			r := make(map[string]interface{})
			r["name"] = strings.ToLower(*i.ParameterName)

			r["value"] = ""
			if i.ParameterValue != nil {
				r["value"] = strings.ToLower(*i.ParameterValue)
			}

			if i.ApplyMethod != nil {
				r["apply_method"] = strings.ToLower(*i.ApplyMethod)
			}

			result = append(result, r)
		}
	}

	return result
}
