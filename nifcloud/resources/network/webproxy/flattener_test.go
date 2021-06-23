package webproxy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"description":                   "test_description",
		"router_id":                     "test_router_id",
		"router_name":                   "test_router_name",
		"listen_port":                   "test_listen_port",
		"name_server":                   "test_name_server",
		"bypass_interface_network_id":   "test_bypass_interface_network_id",
		"bypass_interface_network_name": "test_bypass_interface_network_name",
		"listen_interface_network_id":   "test_listen_interface_network_id",
		"listen_interface_network_name": "test_listen_interface_network_name",
	})
	rd.SetId("test_router_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.NiftyDescribeWebProxiesResponse
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
				res: &computing.NiftyDescribeWebProxiesResponse{
					NiftyDescribeWebProxiesOutput: &computing.NiftyDescribeWebProxiesOutput{
						WebProxy: []computing.WebProxyOfNiftyDescribeWebProxies{
							{
								RouterId:    nifcloud.String("test_router_id"),
								RouterName:  nifcloud.String("test_router_name"),
								Description: nifcloud.String("test_description"),
								ListenPort:  nifcloud.String("test_listen_port"),
								ListenInterface: &computing.ListenInterface{
									NetworkName: nifcloud.String("test_listen_interface_network_name"),
									NetworkId:   nifcloud.String("test_listen_interface_network_id"),
								},
								BypassInterface: &computing.BypassInterface{
									NetworkName: nifcloud.String("test_bypass_interface_network_name"),
									NetworkId:   nifcloud.String("test_bypass_interface_network_id"),
								},
								Option: &computing.OptionOfNiftyDescribeWebProxies{
									NameServer: nifcloud.String("test_name_server"),
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
				res: &computing.NiftyDescribeWebProxiesResponse{
					NiftyDescribeWebProxiesOutput: &computing.NiftyDescribeWebProxiesOutput{
						WebProxy: []computing.WebProxyOfNiftyDescribeWebProxies{},
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
