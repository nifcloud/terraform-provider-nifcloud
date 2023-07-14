package remoteaccessvpngateway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateRemoteAccessVpnGatewayInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":   1,
		"availability_zone": "test_availability_zone",
		"description":       "test_description",
		"network_interface": []interface{}{map[string]interface{}{
			"ip_address": "test_ip_address",
			"network_id": "test_network_id",
		}},
		"name":               "test_remote_access_vpn_gateway_name",
		"type":               "test_type",
		"cipher_suite":       []interface{}{"test_cipher_suite"},
		"pool_network_cidr":  "test_pool_network_cidr",
		"ssl_certificate_id": "test_ssl_certificate_id",
		"ca_certificate_id":  "test_ca_certificate_id",
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateRemoteAccessVpnGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateRemoteAccessVpnGatewayInput{
				AccountingType: nifcloud.Int32(1),
				Placement: &types.RequestPlacementOfCreateRemoteAccessVpnGateway{
					AvailabilityZone: nifcloud.String("test_availability_zone"),
				},
				Description: nifcloud.String("test_description"),
				NetworkInterface: []types.RequestNetworkInterfaceOfCreateRemoteAccessVpnGateway{
					{
						IpAddress: nifcloud.String("test_ip_address"),
						NetworkId: nifcloud.String("test_network_id"),
					},
				},
				RemoteAccessVpnGatewayName: nifcloud.String("test_remote_access_vpn_gateway_name"),
				RemoteAccessVpnGatewayType: types.RemoteAccessVpnGatewayTypeOfCreateRemoteAccessVpnGatewayRequest("test_type"),
				PoolNetworkCidr:            nifcloud.String("test_pool_network_cidr"),
				SSLCertificateId:           nifcloud.String("test_ssl_certificate_id"),
				CACertificateId:            nifcloud.String("test_ca_certificate_id"),
				CipherSuite:                []string{"test_cipher_suite"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateRemoteAccessVpnGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeRemoteAccessVpnGatewaysInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeRemoteAccessVpnGatewaysInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeRemoteAccessVpnGatewaysInput{
				RemoteAccessVpnGatewayId: []string{"test_remote_access_vpn_gateway_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeRemoteAccessVpnGatewaysInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeRemoteAccessVpnGatewayClientConfigInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeRemoteAccessVpnGatewayClientConfigInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeRemoteAccessVpnGatewayClientConfigInput{
				RemoteAccessVpnGatewayId: nifcloud.String("test_remote_access_vpn_gateway_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeRemoteAccessVpnGatewayClientConfigInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyRemoteAccessVpnGatewayAttributeInputForRemoteAccessVpnGatewayName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name": "test_remote_access_vpn_gateway_name",
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyRemoteAccessVpnGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyRemoteAccessVpnGatewayAttributeInput{
				RemoteAccessVpnGatewayId:   nifcloud.String("test_remote_access_vpn_gateway_id"),
				RemoteAccessVpnGatewayName: nifcloud.String("test_remote_access_vpn_gateway_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyRemoteAccessVpnGatewayAttributeInputForRemoteAccessVpnGatewayName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyRemoteAccessVpnGatewayAttributeInputForAccountingType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type": "test_accounting_type",
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyRemoteAccessVpnGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyRemoteAccessVpnGatewayAttributeInput{
				RemoteAccessVpnGatewayId: nifcloud.String("test_remote_access_vpn_gateway_id"),
				AccountingType:           types.AccountingTypeOfModifyRemoteAccessVpnGatewayAttributeRequest("test_accounting_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyRemoteAccessVpnGatewayAttributeInputForAccountingType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyRemoteAccessVpnGatewayAttributeInputForDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"description": "test_description",
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyRemoteAccessVpnGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyRemoteAccessVpnGatewayAttributeInput{
				RemoteAccessVpnGatewayId: nifcloud.String("test_remote_access_vpn_gateway_id"),
				Description:              nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyRemoteAccessVpnGatewayAttributeInputForDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyRemoteAccessVpnGatewayAttributeInputForType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"type": "test_type",
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyRemoteAccessVpnGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyRemoteAccessVpnGatewayAttributeInput{
				RemoteAccessVpnGatewayId:   nifcloud.String("test_remote_access_vpn_gateway_id"),
				RemoteAccessVpnGatewayType: types.RemoteAccessVpnGatewayTypeOfModifyRemoteAccessVpnGatewayAttributeRequest("test_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyRemoteAccessVpnGatewayAttributeInputForType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteRemoteAccessVpnGatewayInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteRemoteAccessVpnGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteRemoteAccessVpnGatewayInput{
				RemoteAccessVpnGatewayId: nifcloud.String("test_remote_access_vpn_gateway_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteRemoteAccessVpnGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandSetRemoteAccessVpnGatewayCACertificateInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ca_certificate_id": "test_ca_certificate_id",
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.SetRemoteAccessVpnGatewayCACertificateInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.SetRemoteAccessVpnGatewayCACertificateInput{
				RemoteAccessVpnGatewayId: nifcloud.String("test_remote_access_vpn_gateway_id"),
				CACertificateId:          nifcloud.String("test_ca_certificate_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandSetRemoteAccessVpnGatewayCACertificateInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUnsetRemoteAccessVpnGatewayCACertificateInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.UnsetRemoteAccessVpnGatewayCACertificateInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.UnsetRemoteAccessVpnGatewayCACertificateInput{
				RemoteAccessVpnGatewayId: nifcloud.String("test_remote_access_vpn_gateway_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUnsetRemoteAccessVpnGatewayCACertificateInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandSetRemoteAccessVpnGatewaySSLCertificateInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ssl_certificate_id": "test_ssl_certificate_id",
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.SetRemoteAccessVpnGatewaySSLCertificateInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.SetRemoteAccessVpnGatewaySSLCertificateInput{
				RemoteAccessVpnGatewayId: nifcloud.String("test_remote_access_vpn_gateway_id"),
				SSLCertificateId:         nifcloud.String("test_ssl_certificate_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandSetRemoteAccessVpnGatewaySSLCertificateInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandCreateRemoteAccessVpnGatewayUsersInput(t *testing.T) {
	user := map[string]interface{}{
		"name":        "test_user",
		"description": "test_description",
		"password":    "test_password",
	}

	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"user": []interface{}{user},
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateRemoteAccessVpnGatewayUsersInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateRemoteAccessVpnGatewayUsersInput{
				RemoteAccessVpnGatewayId: nifcloud.String("test_remote_access_vpn_gateway_id"),
				RemoteUser: []types.RequestRemoteUser{{
					UserName:    nifcloud.String("test_user"),
					Password:    nifcloud.String("test_password"),
					Description: nifcloud.String("test_description"),
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateRemoteAccessVpnGatewayUsersInput(tt.args, user)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteRemoteAccessVpnGatewayUsersInput(t *testing.T) {
	user := map[string]interface{}{
		"name":        "test_user",
		"description": "test_description",
		"password":    "test_password",
	}

	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"user": []interface{}{user},
	})
	rd.SetId("test_remote_access_vpn_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteRemoteAccessVpnGatewayUsersInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteRemoteAccessVpnGatewayUsersInput{
				RemoteAccessVpnGatewayId: nifcloud.String("test_remote_access_vpn_gateway_id"),
				RemoteUser: []types.RequestRemoteUserOfDeleteRemoteAccessVpnGatewayUsers{{
					UserName: nifcloud.String("test_user"),
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteRemoteAccessVpnGatewayUsersInput(tt.args, user)
			assert.Equal(t, tt.want, got)
		})
	}
}
