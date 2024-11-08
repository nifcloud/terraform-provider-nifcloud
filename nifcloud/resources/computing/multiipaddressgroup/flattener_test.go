package multiipaddressgroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_id")

	wantRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_name",
		"description":       "test_description",
		"availability_zone": "test_availability_zone",
		"default_gateway":   "192.168.0.1",
		"subnet_mask":       "255.255.255.0",
		"ip_address_count":  "3",
		"ip_addresses": []interface{}{
			"192.168.0.11",
			"192.168.0.12",
			"192.168.0.13",
		},
		"ip_addresses.#": "3",
	})
	wantRd.SetId("test_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeMultiIpAddressGroupsOutput
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
				res: &computing.DescribeMultiIpAddressGroupsOutput{
					MultiIpAddressGroupsSet: []types.MultiIpAddressGroupsSet{
						{
							MultiIpAddressGroupId:   nifcloud.String("test_id"),
							MultiIpAddressGroupName: nifcloud.String("test_name"),
							Description:             nifcloud.String("test_description"),
							AvailabilityZone:        nifcloud.String("test_availability_zone"),
							MultiIpAddressNetwork: &types.MultiIpAddressNetworkOfDescribeMultiIpAddressGroups{
								DefaultGateway: nifcloud.String("192.168.0.1"),
								SubnetMask:     nifcloud.String("255.255.255.0"),
								IpAddressesSet: []types.IpAddressesSet{
									{
										IpAddress: nifcloud.String("192.168.0.11"),
									},
									{
										IpAddress: nifcloud.String("192.168.0.12"),
									},
									{
										IpAddress: nifcloud.String("192.168.0.13"),
									},
								},
							},
						},
					},
				},
			},
			want: wantRd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d:   schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{}),
				res: &computing.DescribeMultiIpAddressGroupsOutput{},
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
