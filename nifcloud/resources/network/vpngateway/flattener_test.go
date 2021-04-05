package vpngateway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"vpn_gateway_id":             "test_vpn_gateway_id",
		"name":                       "test_name",
		"type":                       "test_type",
		"availability_zone":          "test_availability_zone",
		"accounting_type":            "test_accounting_type",
		"description":                "test_description",
		"network_id":                 "test_network_id",
		"network_name":               "test_network_name",
		"ip_address":                 "test_ip_address",
		"public_ip_address":          "test_public_ip_address",
		"security_group":             "test_security_group",
		"route_table_id":             "test_route_table_id",
		"route_table_association_id": "test_route_table_association_id",
	})
	rd.SetId("test_vpngateway_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeVpnGatewaysResponse
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
				res: &computing.DescribeVpnGatewaysResponse{
					DescribeVpnGatewaysOutput: &computing.DescribeVpnGatewaysOutput{
						VpnGatewaySet: []computing.VpnGatewaySetOfDescribeVpnGateways{
							{
								VpnGatewayId:               nifcloud.String("test_vpngateway_id"),
								NiftyVpnGatewayName:        nifcloud.String("test_name"),
								NiftyVpnGatewayType:        nifcloud.String("test_type"),
								AvailabilityZone:           nifcloud.String("test_availability_zone"),
								NextMonthAccountingType:    nifcloud.String("test_accounting_type"),
								NiftyVpnGatewayDescription: nifcloud.String("test_description"),
								NetworkInterfaceSet: []computing.NetworkInterfaceSetOfDescribeVpnGateways{
									{
										NetworkId:   nifcloud.String("test_network_id"),
										NetworkName: nifcloud.String("test_network_name"),
										IpAddress:   nifcloud.String("test_ip_address"),
									},
									{
										NetworkId: nifcloud.String("net-COMMON_GLOBAL"),
										IpAddress: nifcloud.String("test_public_ip_address"),
									},
								},
								GroupSet: []computing.GroupSetOfDescribeVpnGateways{
									{
										GroupId: nifcloud.String("test_security_group"),
									},
								},
								RouteTableId:            nifcloud.String("test_route_table_id"),
								RouteTableAssociationId: nifcloud.String("test_route_table_association_id"),
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
				res: &computing.DescribeVpnGatewaysResponse{
					DescribeVpnGatewaysOutput: &computing.DescribeVpnGatewaysOutput{
						VpnGatewaySet: []computing.VpnGatewaySetOfDescribeVpnGateways{},
					},
				},
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			fmt.Printf("t = %v\n", t)
			fmt.Printf("err = %v\n", err)
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
