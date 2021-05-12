package nassecuritygroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/stretchr/testify/assert"
)

func TestFlattenForCidrIP(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name":        "test_group_name",
		"description":       "test_description",
		"availability_zone": "test_zone",
		"rule": []interface{}{map[string]interface{}{
			"cidr_ip": "0.0.0.0/0",
		}},
	})
	rd.SetId("test_group_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *nas.DescribeNASSecurityGroupsResponse
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
				res: &nas.DescribeNASSecurityGroupsResponse{
					DescribeNASSecurityGroupsOutput: &nas.DescribeNASSecurityGroupsOutput{
						NASSecurityGroups: []nas.NASSecurityGroup{
							{
								NASSecurityGroupName:        nifcloud.String("test_group_name"),
								NASSecurityGroupDescription: nifcloud.String("test_description"),
								AvailabilityZone:            nifcloud.String("test_zone"),
								IPRanges: []nas.IPRange{
									{
										CIDRIP: nifcloud.String("0.0.0.0/0"),
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
				res: &nas.DescribeNASSecurityGroupsResponse{
					DescribeNASSecurityGroupsOutput: &nas.DescribeNASSecurityGroupsOutput{
						NASSecurityGroups: []nas.NASSecurityGroup{},
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

func TestFlattenForSecurityGroupName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name":        "test_group_name",
		"description":       "test_description",
		"availability_zone": "test_zone",
		"rule": []interface{}{map[string]interface{}{
			"security_group_name": "test_security_group_name",
		}},
	})
	rd.SetId("test_group_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *nas.DescribeNASSecurityGroupsResponse
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
				res: &nas.DescribeNASSecurityGroupsResponse{
					DescribeNASSecurityGroupsOutput: &nas.DescribeNASSecurityGroupsOutput{
						NASSecurityGroups: []nas.NASSecurityGroup{
							{
								NASSecurityGroupName:        nifcloud.String("test_group_name"),
								NASSecurityGroupDescription: nifcloud.String("test_description"),
								AvailabilityZone:            nifcloud.String("test_zone"),
								SecurityGroups: []nas.SecurityGroup{
									{
										SecurityGroupName: nifcloud.String("test_security_group_name"),
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
				res: &nas.DescribeNASSecurityGroupsResponse{
					DescribeNASSecurityGroupsOutput: &nas.DescribeNASSecurityGroupsOutput{
						NASSecurityGroups: []nas.NASSecurityGroup{},
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
