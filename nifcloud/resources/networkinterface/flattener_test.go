package networkinterface

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_address":           "test_ip_address",
		"network_id":           "test_network_id",
		"network_interface_id": "test_network_interface_id",
		"description":          "test_description",
		"availability_zone":    "test_availability_zone",
	})
	rd.SetId("test_network_interface_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeNetworkInterfacesResponse
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
				res: &computing.DescribeNetworkInterfacesResponse{
					DescribeNetworkInterfacesOutput: &computing.DescribeNetworkInterfacesOutput{
						NetworkInterfaceSet: []computing.NetworkInterfaceSetOfDescribeNetworkInterfaces{
							{
								NiftyNetworkId:     nifcloud.String("test_network_id"),
								NetworkInterfaceId: nifcloud.String("test_network_interface_id"),
								PrivateIpAddress:   nifcloud.String("test_ip_address"),
								Description:        nifcloud.String("test_description"),
								AvailabilityZone:   nifcloud.String("test_availability_zone"),
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
				res: &computing.DescribeNetworkInterfacesResponse{
					DescribeNetworkInterfacesOutput: &computing.DescribeNetworkInterfacesOutput{
						NetworkInterfaceSet: []computing.NetworkInterfaceSetOfDescribeNetworkInterfaces{},
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
