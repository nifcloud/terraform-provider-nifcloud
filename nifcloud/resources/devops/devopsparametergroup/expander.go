package devopsparametergroup

import (
	"reflect"

	"github.com/ettle/strcase"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
)

func expandCreateParameterGroupInput(d *schema.ResourceData) *devops.CreateParameterGroupInput {
	return &devops.CreateParameterGroupInput{
		ParameterGroupName: nifcloud.String(d.Get("name").(string)),
		Description:        nifcloud.String(d.Get("description").(string)),
	}
}

func expandUpdateParameterGroupInput(d *schema.ResourceData, parameters *types.RequestParameters) *devops.UpdateParameterGroupInput {
	return &devops.UpdateParameterGroupInput{
		ParameterGroupName:        nifcloud.String(d.Id()),
		ChangedParameterGroupName: nifcloud.String(d.Get("name").(string)),
		Description:               nifcloud.String(d.Get("description").(string)),
		Parameters:                parameters,
	}
}

func expandGetParameterGroupInput(d *schema.ResourceData) *devops.GetParameterGroupInput {
	return &devops.GetParameterGroupInput{
		ParameterGroupName: nifcloud.String(d.Id()),
	}
}

func expandDeleteParameterGroupInput(d *schema.ResourceData) *devops.DeleteParameterGroupInput {
	return &devops.DeleteParameterGroupInput{
		ParameterGroupName: nifcloud.String(d.Id()),
	}
}

func expandUpdateParameterGroupParameters(configured []map[string]string) *types.RequestParameters {
	parameters := &types.RequestParameters{}

	structParams := reflect.Indirect(reflect.ValueOf(parameters))
	for _, rawParam := range configured {
		if rawParam["name"] == "" {
			continue
		}

		name := strcase.ToPascal(rawParam["name"])
		value := rawParam["value"]
		structParams.FieldByName(name).Set(reflect.ValueOf(nifcloud.String(value)))
	}

	return parameters
}

func expandParameters(configured []interface{}) []types.Parameters {
	var parameters []types.Parameters

	for _, raw := range configured {
		rawParam := raw.(map[string]interface{})

		if rawParam["name"].(string) == "" {
			continue
		}

		param := types.Parameters{
			Name:  nifcloud.String(rawParam["name"].(string)),
			Value: nifcloud.String(rawParam["value"].(string)),
		}
		parameters = append(parameters, param)
	}

	return parameters
}
