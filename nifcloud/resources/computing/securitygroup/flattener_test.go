package securitygroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name":        "test_group_name",
		"description":       "test_description",
		"log_limit":         1000,
		"availability_zone": "test_availability_zone",
	})
	rd.SetId("test_group_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeSecurityGroupsResponse
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
				res: &computing.DescribeSecurityGroupsResponse{
					DescribeSecurityGroupsOutput: &computing.DescribeSecurityGroupsOutput{
						SecurityGroupInfo: []computing.SecurityGroupInfo{
							{
								GroupName:        nifcloud.String("test_group_name"),
								GroupDescription: nifcloud.String("test_description"),
								GroupLogLimit:    nifcloud.Int64(1000),
								AvailabilityZone: nifcloud.String("test_availability_zone"),
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
				res: &computing.DescribeSecurityGroupsResponse{
					DescribeSecurityGroupsOutput: &computing.DescribeSecurityGroupsOutput{
						SecurityGroupInfo: []computing.SecurityGroupInfo{},
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
