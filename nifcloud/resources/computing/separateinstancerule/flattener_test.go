package separateinstancerule

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":               "test_separate_name",
		"availability_zone":  "test_availability_zone",
		"description":        "test_description",
		"instance_id":        []interface{}{"test_instance_id1", "test_instance_id2"},
		"instance_unique_id": []interface{}{"test_instance_unique_id1", "test_instance_unique_id2"},
	})
	rd.SetId("test_separate_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.NiftyDescribeSeparateInstanceRulesResponse
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
				res: &computing.NiftyDescribeSeparateInstanceRulesResponse{
					NiftyDescribeSeparateInstanceRulesOutput: &computing.NiftyDescribeSeparateInstanceRulesOutput{
						SeparateInstanceRulesInfo: []computing.SeparateInstanceRulesInfo{
							{
								SeparateInstanceRuleName:        nifcloud.String("test_separate_name"),
								AvailabilityZone:                nifcloud.String("test_availability_zone"),
								SeparateInstanceRuleDescription: nifcloud.String("test_description"),
								InstancesSet: []computing.InstancesSetOfNiftyDescribeSeparateInstanceRules{
									{
										InstanceId: nifcloud.String("test_instance_id1"),
									},
									{
										InstanceId: nifcloud.String("test_instance_id2"),
									},
									{
										InstanceUniqueId: nifcloud.String("test_instance_unique_id1"),
									},
									{
										InstanceUniqueId: nifcloud.String("test_instance_unique_id2"),
									},
								},
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
				d: wantNotFoundRd,
				res: &computing.NiftyDescribeSeparateInstanceRulesResponse{
					NiftyDescribeSeparateInstanceRulesOutput: &computing.NiftyDescribeSeparateInstanceRulesOutput{
						SeparateInstanceRulesInfo: []computing.SeparateInstanceRulesInfo{},
					},
				},
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			assert.NoError(t, err)

			if tt.args.res == nil {
				return
			}

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
