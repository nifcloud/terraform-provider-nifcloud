package dbparametergroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateDBParameterGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"family":      "test_family",
		"description": "test_description",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.CreateDBParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.CreateDBParameterGroupInput{
				DBParameterGroupName:   nifcloud.String("test_name"),
				DBParameterGroupFamily: nifcloud.String("test_family"),
				Description:            nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateDBParameterGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyDBParameterGroupInput(t *testing.T) {
	parameters := []rdb.RequestParameters{
		{
			ParameterName:  nifcloud.String("test_name"),
			ParameterValue: nifcloud.String("test_value"),
			ApplyMethod:    nifcloud.String("test_apply_method"),
		},
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_db_parameter_group_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.ModifyDBParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.ModifyDBParameterGroupInput{
				DBParameterGroupName: nifcloud.String("test_db_parameter_group_id"),
				Parameters:           parameters,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyDBParameterGroupInput(tt.args, parameters)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandResetDBParameterGroupInput(t *testing.T) {
	parameters := []rdb.RequestParametersOfResetDBParameterGroup{
		{
			ParameterName: nifcloud.String("test_name"),
			ApplyMethod:   nifcloud.String("test_apply_method"),
		},
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_db_parameter_group_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.ResetDBParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.ResetDBParameterGroupInput{
				DBParameterGroupName: nifcloud.String("test_db_parameter_group_id"),
				Parameters:           parameters,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandResetDBParameterGroupInput(tt.args, parameters)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeDBParameterGroupsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_db_parameter_group_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.DescribeDBParameterGroupsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.DescribeDBParameterGroupsInput{
				DBParameterGroupName: nifcloud.String("test_db_parameter_group_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeDBParameterGroupsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeDBParametersInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_db_parameter_group_id")

	type args struct {
		data   *schema.ResourceData
		marker string
	}

	tests := []struct {
		name string
		args args
		want *rdb.DescribeDBParametersInput
	}{
		{
			name: "expands the resource data without marker",
			args: args{
				data:   rd,
				marker: "",
			},
			want: &rdb.DescribeDBParametersInput{
				DBParameterGroupName: nifcloud.String("test_db_parameter_group_id"),
			},
		},
		{
			name: "expands the resource data with marker",
			args: args{
				data:   rd,
				marker: "test_marker",
			},
			want: &rdb.DescribeDBParametersInput{
				DBParameterGroupName: nifcloud.String("test_db_parameter_group_id"),
				Marker:               nifcloud.String("test_marker"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeDBParametersInput(tt.args.data, tt.args.marker)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteDBParameterGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_db_parameter_group_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.DeleteDBParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.DeleteDBParameterGroupInput{
				DBParameterGroupName: nifcloud.String("test_db_parameter_group_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteDBParameterGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandParameters(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"parameter": []interface{}{
			map[string]interface{}{
				"name":         "test_name_01",
				"value":        "test_value_01",
				"apply_method": "test_apply_method_01",
			},
			map[string]interface{}{
				"name":         "test_name_02",
				"value":        "test_value_02",
				"apply_method": "test_apply_method_02",
			},
			map[string]interface{}{
				"value":        "test_value_03",
				"apply_method": "test_apply_method_03",
			},
		},
	})

	tests := []struct {
		name string
		args []interface{}
		want []rdb.Parameter
	}{
		{
			name: "expands the resource data",
			args: rd.Get("parameter").(*schema.Set).List(),
			want: []rdb.Parameter{
				{
					ParameterName:  nifcloud.String("test_name_02"),
					ParameterValue: nifcloud.String("test_value_02"),
					ApplyMethod:    nifcloud.String("test_apply_method_02"),
				},
				{
					ParameterName:  nifcloud.String("test_name_01"),
					ParameterValue: nifcloud.String("test_value_01"),
					ApplyMethod:    nifcloud.String("test_apply_method_01"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandParameters(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyDBParameterGroupParameters(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"parameter": []interface{}{
			map[string]interface{}{
				"name":         "test_name_01",
				"value":        "test_value_01",
				"apply_method": "test_apply_method_01",
			},
			map[string]interface{}{
				"name":         "test_name_02",
				"value":        "test_value_02",
				"apply_method": "test_apply_method_02",
			},
			map[string]interface{}{
				"value":        "test_value_03",
				"apply_method": "test_apply_method_03",
			},
		},
	})

	tests := []struct {
		name string
		args []interface{}
		want []rdb.RequestParameters
	}{
		{
			name: "expands the resource data",
			args: rd.Get("parameter").(*schema.Set).List(),
			want: []rdb.RequestParameters{
				{
					ParameterName:  nifcloud.String("test_name_02"),
					ParameterValue: nifcloud.String("test_value_02"),
					ApplyMethod:    nifcloud.String("test_apply_method_02"),
				},
				{
					ParameterName:  nifcloud.String("test_name_01"),
					ParameterValue: nifcloud.String("test_value_01"),
					ApplyMethod:    nifcloud.String("test_apply_method_01"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyDBParameterGroupParameters(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandResetDBParameterGroupParameters(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"parameter": []interface{}{
			map[string]interface{}{
				"name":         "test_name_01",
				"value":        "test_value_01",
				"apply_method": "test_apply_method_01",
			},
			map[string]interface{}{
				"name":         "test_name_02",
				"value":        "test_value_02",
				"apply_method": "test_apply_method_02",
			},
			map[string]interface{}{
				"value":        "test_value_03",
				"apply_method": "test_apply_method_03",
			},
		},
	})

	tests := []struct {
		name string
		args []interface{}
		want []rdb.RequestParametersOfResetDBParameterGroup
	}{
		{
			name: "expands the resource data",
			args: rd.Get("parameter").(*schema.Set).List(),
			want: []rdb.RequestParametersOfResetDBParameterGroup{
				{
					ParameterName: nifcloud.String("test_name_02"),
					ApplyMethod:   nifcloud.String("test_apply_method_02"),
				},
				{
					ParameterName: nifcloud.String("test_name_01"),
					ApplyMethod:   nifcloud.String("test_apply_method_01"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandResetDBParameterGroupParameters(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
