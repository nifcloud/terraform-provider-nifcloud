package securitygrouprule

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"type":                       "IN",
		"cidr_ip":                    "0.0.0.0/0",
		"from_port":                  1,
		"security_group_names":       []interface{}{"test_security_group_name"},
		"protocol":                   "TCP",
		"source_security_group_name": "test_source_security_group_name",
		"to_port":                    65535,
		"description":                "test_description",
	})
	rd.SetId("some")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	wantNotFoundRuleRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"security_group_names": []interface{}{"test_security_group_name"},
	})

	type args struct {
		res *computing.DescribeSecurityGroupsOutput
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
				res: &computing.DescribeSecurityGroupsOutput{
					SecurityGroupInfo: []types.SecurityGroupInfo{
						{
							GroupName: nifcloud.String("test_security_group_name"),
							IpPermissions: []types.IpPermissions{{
								InOut:       nifcloud.String("IN"),
								IpRanges:    []types.IpRanges{{CidrIp: nifcloud.String("0.0.0.0/0")}},
								FromPort:    nifcloud.Int32(1),
								IpProtocol:  nifcloud.String("TCP"),
								Groups:      []types.Groups{{GroupName: nifcloud.String("test_source_security_group_name")}},
								ToPort:      nifcloud.Int32(65535),
								Description: nifcloud.String("test_description"),
							}},
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
				res: &computing.DescribeSecurityGroupsOutput{
					SecurityGroupInfo: []types.SecurityGroupInfo{},
				},
			},
			want: wantNotFoundRd,
		},
		{
			name: "flattens the response even when the resource has been removed externally in security group",
			args: args{
				d: wantNotFoundRuleRd,
				res: &computing.DescribeSecurityGroupsOutput{
					SecurityGroupInfo: []types.SecurityGroupInfo{{GroupName: nifcloud.String("test_security_group_name")}},
				},
			},
			want: wantNotFoundRuleRd,
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
