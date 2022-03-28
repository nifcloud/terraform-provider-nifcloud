package vpnconnection

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateVpnConnectionInputForIdTunnel(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"type":                                        "L2TPv3 / IPsec",
		"vpn_gateway_id":                              "test_vpn_gateway_id",
		"customer_gateway_id":                         "test_customer_gateway_id",
		"tunnel_type":                                 "L2TPv3",
		"tunnel_mode":                                 "Unmanaged",
		"tunnel_encapsulation":                        "UDP",
		"tunnel_id":                                   "1",
		"tunnel_peer_id":                              "2",
		"tunnel_session_id":                           "1",
		"tunnel_peer_session_id":                      "2",
		"tunnel_source_port":                          "7777",
		"tunnel_destination_port":                     "7777",
		"ipsec_config_encryption_algorithm":           "AES256",
		"ipsec_config_hash_algorithm":                 "SHA256",
		"ipsec_config_pre_shared_key":                 "test_pre_shared_key",
		"ipsec_config_internet_key_exchange":          "IKEv2",
		"ipsec_config_internet_key_exchange_lifetime": 300,
		"ipsec_config_encapsulating_security_payload_lifetime": 300,
		"ipsec_config_diffie_hellman_group":                    5,
		"mtu":                                                  "1000",
		"description":                                          "test_description",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateVpnConnectionInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateVpnConnectionInput{
				Type:              types.TypeOfCreateVpnConnectionRequest("L2TPv3 / IPsec"),
				VpnGatewayId:      nifcloud.String("test_vpn_gateway_id"),
				CustomerGatewayId: nifcloud.String("test_customer_gateway_id"),
				NiftyTunnel: &types.RequestNiftyTunnel{
					Type:            types.TypeOfNiftyTunnelForCreateVpnConnection("L2TPv3"),
					Mode:            types.ModeOfNiftyTunnelForCreateVpnConnection("Unmanaged"),
					Encapsulation:   types.EncapsulationOfNiftyTunnelForCreateVpnConnection("UDP"),
					TunnelId:        nifcloud.String("1"),
					PeerTunnelId:    nifcloud.String("2"),
					SessionId:       nifcloud.String("1"),
					PeerSessionId:   nifcloud.String("2"),
					SourcePort:      nifcloud.String("7777"),
					DestinationPort: nifcloud.String("7777"),
				},
				NiftyIpsecConfiguration: &types.RequestNiftyIpsecConfiguration{
					EncryptionAlgorithm:                  types.EncryptionAlgorithmOfNiftyIpsecConfigurationForCreateVpnConnection("AES256"),
					HashAlgorithm:                        types.HashAlgorithmOfNiftyIpsecConfigurationForCreateVpnConnection("SHA256"),
					PreSharedKey:                         nifcloud.String("test_pre_shared_key"),
					InternetKeyExchange:                  types.InternetKeyExchangeOfNiftyIpsecConfigurationForCreateVpnConnection("IKEv2"),
					InternetKeyExchangeLifetime:          nifcloud.Int32(300),
					EncapsulatingSecurityPayloadLifetime: nifcloud.Int32(300),
					DiffieHellmanGroup:                   nifcloud.Int32(5),
				},
				NiftyVpnConnectionMtu:         nifcloud.String("1000"),
				NiftyVpnConnectionDescription: nifcloud.String("test_description"),
				Agreement:                     nifcloud.Bool(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateVpnConnectionInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandCreateVpnConnectionInputForNameNoTunnel(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"type":                                                 "IPsec",
		"vpn_gateway_name":                                     "test_vpn_gateway_name",
		"customer_gateway_name":                                "test_customer_gateway_name",
		"ipsec_config_encryption_algorithm":                    "AES256",
		"ipsec_config_hash_algorithm":                          "SHA256",
		"ipsec_config_pre_shared_key":                          "test_pre_shared_key",
		"ipsec_config_internet_key_exchange":                   "IKEv2",
		"ipsec_config_internet_key_exchange_lifetime":          300,
		"ipsec_config_encapsulating_security_payload_lifetime": 300,
		"ipsec_config_diffie_hellman_group":                    5,
		"mtu":                                                  "1000",
		"description":                                          "test_description",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateVpnConnectionInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateVpnConnectionInput{
				Type:                     types.TypeOfCreateVpnConnectionRequest("IPsec"),
				NiftyVpnGatewayName:      nifcloud.String("test_vpn_gateway_name"),
				NiftyCustomerGatewayName: nifcloud.String("test_customer_gateway_name"),
				NiftyIpsecConfiguration: &types.RequestNiftyIpsecConfiguration{
					EncryptionAlgorithm:                  types.EncryptionAlgorithmOfNiftyIpsecConfigurationForCreateVpnConnection("AES256"),
					HashAlgorithm:                        types.HashAlgorithmOfNiftyIpsecConfigurationForCreateVpnConnection("SHA256"),
					PreSharedKey:                         nifcloud.String("test_pre_shared_key"),
					InternetKeyExchange:                  types.InternetKeyExchangeOfNiftyIpsecConfigurationForCreateVpnConnection("IKEv2"),
					InternetKeyExchangeLifetime:          nifcloud.Int32(300),
					EncapsulatingSecurityPayloadLifetime: nifcloud.Int32(300),
					DiffieHellmanGroup:                   nifcloud.Int32(5),
				},
				NiftyVpnConnectionDescription: nifcloud.String("test_description"),
				Agreement:                     nifcloud.Bool(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateVpnConnectionInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeVpnConnectionsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeVpnConnectionsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeVpnConnectionsInput{
				VpnConnectionId: []string{"test_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeVpnConnectionsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteVpnConnectionInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteVpnConnectionInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteVpnConnectionInput{
				VpnConnectionId: nifcloud.String("test_id"),
				Agreement:       nifcloud.Bool(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteVpnConnectionInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
