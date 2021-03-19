package dhcpconfig

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"static_mappiing": []interface{}{map[string]interface{}{
			"static_mappiing_ipaddress":  "192.168.1.10",
			"static_mappiing_macaddress": "00:00:5e:00:53:00",
			"static_mapping_description": "test_description",
		}},
		"ipaddress_pool": []interface{}{map[string]interface{}{
			"ipaddress_pool_start": "192.168.2.1",
			"ipaddress_pool_stop":  "192.168.2.100",
			"description":          "test_description",
		}},
		"dhcp_config_id": "test_dhcp_config_id",
	})
	rd.SetId("test_dhcp_config_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.NiftyDescribeDhcpConfigsResponse
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
				res: &computing.NiftyDescribeDhcpConfigsResponse{
					NiftyDescribeDhcpConfigsOutput: &computing.NiftyDescribeDhcpConfigsOutput{
						DhcpConfigsSet: []computing.DhcpConfigsSet{
							{
								DhcpConfigId: nifcloud.String("test_dhcp_config_id"),
								StaticMappingsSet: []computing.StaticMappingsSet{
									{
										IpAddress:   nifcloud.String("192.168.1.10"),
										MacAddress:  nifcloud.String("00:00:5e:00:53:00"),
										Description: nifcloud.String("test_description"),
									},
								},
								IpAddressPoolsSet: []computing.IpAddressPoolsSet{
									{
										StartIpAddress: nifcloud.String("192.168.2.1"),
										StopIpAddress:  nifcloud.String("192.168.2.100"),
										Description:    nifcloud.String("test_description"),
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
				res: &computing.NiftyDescribeDhcpConfigsResponse{
					NiftyDescribeDhcpConfigsOutput: &computing.NiftyDescribeDhcpConfigsOutput{
						DhcpConfigsSet: []computing.DhcpConfigsSet{},
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
