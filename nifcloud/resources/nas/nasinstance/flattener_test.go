package nasinstance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"allocated_storage":       100,
		"availability_zone":       "test_availability_zone",
		"description":             "test_description",
		"nas_security_group_name": "test_nas_security_group_name",
		"public_ip_address":       "test_public_ip_address",
		"private_ip_address":      "test_private_ip_address",
		"protocol":                "test_protocol",
		"master_username":         "test_master_username",
		"network_id":              "test_network_id",
		"authentication_type":     0,
		"type":                    0,
		"no_root_squash":          "false",
	})
	rd.SetId("test_identifier")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *nas.DescribeNASInstancesOutput
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
				res: &nas.DescribeNASInstancesOutput{
					NASInstances: []types.NASInstances{
						{
							NASInstanceIdentifier:  nifcloud.String("test_identifier"),
							AllocatedStorage:       nifcloud.Int32(100),
							AvailabilityZone:       nifcloud.String("test_availability_zone"),
							NASInstanceDescription: nifcloud.String("test_description"),
							NASSecurityGroups: []types.NASSecurityGroups{
								{
									NASSecurityGroupName: nifcloud.String("test_nas_security_group_name"),
								},
							},
							Endpoint: &types.Endpoint{
								Address:        nifcloud.String("test_public_ip_address"),
								PrivateAddress: nifcloud.String("test_private_ip_address"),
							},
							Protocol:           nifcloud.String("test_protocol"),
							MasterUsername:     nifcloud.String("test_master_username"),
							NetworkId:          nifcloud.String("test_network_id"),
							AuthenticationType: nifcloud.Int32(0),
							NASInstanceType:    nifcloud.Int32(0),
							NoRootSquash:       nifcloud.Bool(false),
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
				res: &nas.DescribeNASInstancesOutput{},
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
