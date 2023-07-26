package remoteaccessvpngateway

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
		"accounting_type":   "test_accounting_type",
		"availability_zone": "test_availability_zone",
		"description":       "test_description",
		"name":              "test_remote_access_vpn_gateway_name",
		"network_interface": []interface{}{
			map[string]interface{}{
				"ip_address": "test_ip_address",
				"network_id": "test_network_id",
			},
		},
		"user": []interface{}{
			map[string]interface{}{
				"name":        "test_user",
				"description": "test_description",
			},
		},
		"remote_access_vpn_gateway_id": "test_remote_access_vpn_gateway_id",
		"type":                         "test_type",
		"cipher_suite":                 []interface{}{"test_cipher_suite"},
		"ca_certificate_id":            "test_ca_certificate_id",
		"ssl_certificate_id":           "test_ssl_certificate_id",
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeRemoteAccessVpnGatewaysOutput
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
				res: &computing.DescribeRemoteAccessVpnGatewaysOutput{
					RemoteAccessVpnGatewaySet: []types.RemoteAccessVpnGatewaySetOfDescribeRemoteAccessVpnGateways{
						{
							AvailabilityZone: nifcloud.String("test_availability_zone"),
							Description:      nifcloud.String("test_description"),
							NetworkInterfaceSet: []types.NetworkInterfaceSetOfDescribeRemoteAccessVpnGateways{
								{
									PrivateIpAddress: nifcloud.String("test_ip_address"),
									NiftyNetworkId:   nifcloud.String("test_network_id"),
								},
							},
							NextMonthAccountingType:    nifcloud.String("test_accounting_type"),
							RemoteAccessVpnGatewayId:   nifcloud.String("test_remote_access_vpn_gateway_id"),
							RemoteAccessVpnGatewayName: nifcloud.String("test_remote_access_vpn_gateway_name"),
							RemoteAccessVpnGatewayType: nifcloud.String("test_type"),
							PoolNetworkCidr:            nifcloud.String("test_pool_network_cidr"),
							RemoteUserSet: []types.RemoteUserSet{
								{
									UserName:    nifcloud.String("test_user"),
									Description: nifcloud.String("test_description"),
								},
							},
							SslCertificateId: nifcloud.String("test_ssl_certificate_id"),
							CaCertificateId:  nifcloud.String("test_ca_certificate_id"),
							CipherSuiteSet:   []types.CipherSuiteSet{{CipherSuite: nifcloud.String("test_cipher_suite")}},
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
				res: &computing.DescribeRemoteAccessVpnGatewaysOutput{
					RemoteAccessVpnGatewaySet: []types.RemoteAccessVpnGatewaySetOfDescribeRemoteAccessVpnGateways{},
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

func TestFlattenClientConfig(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"client_config": "test_client_config",
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeRemoteAccessVpnGatewayClientConfigOutput
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
				res: &computing.DescribeRemoteAccessVpnGatewayClientConfigOutput{
					FileData: nifcloud.String("test_client_config"),
				},
			},
			want: rd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d:   wantNotFoundRd,
				res: &computing.DescribeRemoteAccessVpnGatewayClientConfigOutput{},
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flattenClientConfig(tt.args.d, tt.args.res)
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
