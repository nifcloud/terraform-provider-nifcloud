package devopsparametergroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"description": "test_description",
		"parameter": []interface{}{
			map[string]interface{}{
				"name":  "test_name_01",
				"value": "test_value_01",
			},
		},
	})
	rd.SetId("test_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *devops.GetParameterGroupOutput
		d   *schema.ResourceData
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
				res: &devops.GetParameterGroupOutput{
					ParameterGroup: &types.ParameterGroup{
						ParameterGroupName: nifcloud.String("test_name"),
						Description:        nifcloud.String("test_description"),
						Parameters: []types.Parameters{
							{
								Name:     nifcloud.String("test_name_01"),
								Value:    nifcloud.String("test_value_01"),
								IsSecret: nifcloud.Bool(false),
							},
							{
								Name:     nifcloud.String("test_name_02"),
								Value:    nifcloud.String("test_value_02"),
								IsSecret: nifcloud.Bool(true),
							},
							{
								Name:  nifcloud.String("test_name_03"),
								Value: nifcloud.String("test_value_03"),
							},
							{
								Name:  nil,
								Value: nil,
							},
						},
					},
				},
			},
			want: rd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d:   wantNotFoundRd,
				res: nil,
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
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
			assert.True(t, tt.want.Get("parameter").(*schema.Set).Equal(tt.args.d.Get("parameter")))
		})
	}
}
