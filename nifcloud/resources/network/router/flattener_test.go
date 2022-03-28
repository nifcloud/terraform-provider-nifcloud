package router

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
		"accounting_type":          "test_accounting_type",
		"availability_zone":        "test_availability_zone",
		"description":              "test_description",
		"name":                     "test_router_name",
		"nat_table_id":             "test_nat_table_id",
		"nat_table_association_id": "test_nat_table_association_id",
		"network_interface": []interface{}{
			map[string]interface{}{
				"dhcp":            true,
				"dhcp_config_id":  "test_dhcp_config_id",
				"dhcp_options_id": "test_dhcp_options_id",
				"ip_address":      "test_ip_address",
				"network_id":      "test_network_id",
				"network_name":    "test_network_name",
			},
			map[string]interface{}{
				"ip_address": "test_global_ip_address",
				"network_id": "net-COMMON_GLOBAL",
			},
		},
		"router_id":                  "test_router_id",
		"route_table_id":             "test_route_table_id",
		"route_table_association_id": "test_route_table_association_id",
		"security_group":             "test_security_group",
		"type":                       "test_type",
	})
	rd.SetId("test_router_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.NiftyDescribeRoutersOutput
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
				res: &computing.NiftyDescribeRoutersOutput{
					RouterSet: []types.RouterSetOfNiftyDescribeRouters{
						{
							AvailabilityZone: nifcloud.String("test_availability_zone"),
							Description:      nifcloud.String("test_description"),
							GroupSet: []types.GroupSet{
								{
									GroupId: nifcloud.String("test_security_group"),
								},
							},
							NatTableAssociationId: nifcloud.String("test_nat_table_association_id"),
							NatTableId:            nifcloud.String("test_nat_table_id"),
							NetworkInterfaceSet: []types.NetworkInterfaceSetOfNiftyDescribeRouters{
								{
									Dhcp:          nifcloud.Bool(true),
									DhcpConfigId:  nifcloud.String("test_dhcp_config_id"),
									DhcpOptionsId: nifcloud.String("test_dhcp_options_id"),
									IpAddress:     nifcloud.String("test_ip_address"),
									NetworkId:     nifcloud.String("test_network_id"),
									NetworkName:   nifcloud.String("test_network_name"),
								},
								{
									IpAddress: nifcloud.String("test_global_ip_address"),
									NetworkId: nifcloud.String("net-COMMON_GLOBAL"),
								},
							},
							NextMonthAccountingType: nifcloud.String("test_accounting_type"),
							RouteTableAssociationId: nifcloud.String("test_route_table_association_id"),
							RouteTableId:            nifcloud.String("test_route_table_id"),
							RouterId:                nifcloud.String("test_router_id"),
							RouterName:              nifcloud.String("test_router_name"),
							Type:                    nifcloud.String("test_type"),
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
				res: &computing.NiftyDescribeRoutersOutput{
					RouterSet: []types.RouterSetOfNiftyDescribeRouters{},
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
