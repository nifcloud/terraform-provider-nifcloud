package dbparametergroup

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
)

func expandCreateDBParameterGroupInput(d *schema.ResourceData) *rdb.CreateDBParameterGroupInput {
	return &rdb.CreateDBParameterGroupInput{
		DBParameterGroupName:   nifcloud.String(d.Get("name").(string)),
		DBParameterGroupFamily: nifcloud.String(d.Get("family").(string)),
		Description:            nifcloud.String(d.Get("description").(string)),
	}
}

func expandModifyDBParameterGroupInput(d *schema.ResourceData, parameters []rdb.RequestParameters) *rdb.ModifyDBParameterGroupInput {
	return &rdb.ModifyDBParameterGroupInput{
		DBParameterGroupName: nifcloud.String(d.Id()),
		Parameters:           parameters,
	}
}

func expandResetDBParameterGroupInput(d *schema.ResourceData, parameters []rdb.RequestParametersOfResetDBParameterGroup) *rdb.ResetDBParameterGroupInput {
	return &rdb.ResetDBParameterGroupInput{
		DBParameterGroupName: nifcloud.String(d.Id()),
		Parameters:           parameters,
	}
}

func expandDescribeDBParameterGroupsInput(d *schema.ResourceData) *rdb.DescribeDBParameterGroupsInput {
	return &rdb.DescribeDBParameterGroupsInput{
		DBParameterGroupName: nifcloud.String(d.Id()),
	}
}

func expandDescribeDBParametersInput(d *schema.ResourceData, marker string) *rdb.DescribeDBParametersInput {
	input := &rdb.DescribeDBParametersInput{
		DBParameterGroupName: nifcloud.String(d.Id()),
	}

	if marker != "" {
		input.Marker = nifcloud.String(marker)
	}

	input.Source = nifcloud.String("user")

	return input
}

func expandDeleteDBParameterGroupInput(d *schema.ResourceData) *rdb.DeleteDBParameterGroupInput {
	return &rdb.DeleteDBParameterGroupInput{
		DBParameterGroupName: nifcloud.String(d.Id()),
	}
}

func expandParameters(configured []interface{}) []rdb.Parameters {
	var parameters []rdb.Parameters

	for _, raw := range configured {
		rawParam := raw.(map[string]interface{})

		if rawParam["name"].(string) == "" {
			continue
		}

		param := rdb.Parameters{
			ParameterName:  nifcloud.String(rawParam["name"].(string)),
			ParameterValue: nifcloud.String(rawParam["value"].(string)),
			ApplyMethod:    nifcloud.String(rawParam["apply_method"].(string)),
		}
		parameters = append(parameters, param)
	}

	return parameters
}

func expandModifyDBParameterGroupParameters(configured []interface{}) []rdb.RequestParameters {
	var parameters []rdb.RequestParameters

	for _, raw := range configured {
		rawParam := raw.(map[string]interface{})

		if rawParam["name"].(string) == "" {
			continue
		}

		param := rdb.RequestParameters{
			ParameterName:  nifcloud.String(rawParam["name"].(string)),
			ParameterValue: nifcloud.String(rawParam["value"].(string)),
			ApplyMethod:    nifcloud.String(rawParam["apply_method"].(string)),
		}
		parameters = append(parameters, param)
	}

	return parameters
}

func expandResetDBParameterGroupParameters(configured []interface{}) []rdb.RequestParametersOfResetDBParameterGroup {
	var parameters []rdb.RequestParametersOfResetDBParameterGroup

	for _, raw := range configured {
		rawParam := raw.(map[string]interface{})

		if rawParam["name"].(string) == "" {
			continue
		}

		param := rdb.RequestParametersOfResetDBParameterGroup{
			ParameterName: nifcloud.String(rawParam["name"].(string)),
			ApplyMethod:   nifcloud.String(rawParam["apply_method"].(string)),
		}
		parameters = append(parameters, param)
	}

	return parameters
}
