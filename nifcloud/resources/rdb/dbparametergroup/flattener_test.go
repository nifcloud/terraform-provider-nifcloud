package dbparametergroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"family":      "test_family",
		"description": "test_description",
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
				"name":         "test_name_03",
				"value":        "test_value_03",
				"apply_method": "test_apply_method_03",
			},
		},
	})
	rd.SetId("test_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		groups     *rdb.DescribeDBParameterGroupsOutput
		parameters []types.Parameters
		d          *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want *schema.ResourceData
	}{
		{
			name: "flattens the response",
			args: args{
				d: rd,
				groups: &rdb.DescribeDBParameterGroupsOutput{
					DBParameterGroups: []types.DBParameterGroupsOfDescribeDBParameterGroups{
						{
							DBParameterGroupName:   nifcloud.String("test_name"),
							DBParameterGroupFamily: nifcloud.String("test_family"),
							Description:            nifcloud.String("test_description"),
						},
					},
				},
				parameters: []types.Parameters{
					{
						ParameterName:  nifcloud.String("test_name_01"),
						ParameterValue: nifcloud.String("test_value_01"),
					},
					{
						ParameterName:  nifcloud.String("test_name_02"),
						ParameterValue: nifcloud.String("test_value_02"),
					},
					{
						ParameterName:  nifcloud.String("test_name_03"),
						ParameterValue: nifcloud.String("test_value_03"),
						ApplyMethod:    nifcloud.String("test_apply_method_03"),
					},
					{
						ParameterName:  nifcloud.String("test_name_04"),
						ParameterValue: nifcloud.String("test_value_04"),
					},
					{
						ParameterName:  nil,
						ParameterValue: nil,
					},
				},
			},
			want: rd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d: wantNotFoundRd,
				groups: &rdb.DescribeDBParameterGroupsOutput{
					DBParameterGroups: []types.DBParameterGroupsOfDescribeDBParameterGroups{},
				},
				parameters: []types.Parameters{},
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.groups, tt.args.parameters)
			assert.NoError(t, err)

			wantState := tt.want.State()
			if wantState == nil {
				tt.want.SetId("some")
				wantState = tt.want.State()
			}

			gotState := tt.args.d.State()
			if gotState == nil {
				tt.args.d.SetId("some")
				gotState = tt.args.d.State()
			}

			assert.Equal(t, wantState.Attributes, gotState.Attributes)
		})
	}
}
