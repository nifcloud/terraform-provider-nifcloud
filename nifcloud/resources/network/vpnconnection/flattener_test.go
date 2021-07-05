package vpnconnection

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"type":                                        "L2TPv3 / IPsec",
		"vpn_gateway_id":                              "test_vpn_gateway_id",
		"vpn_gateway_name":                            "test_vpn_gateway_name",
		"customer_gateway_id":                         "test_customer_gateway_id",
		"customer_gateway_name":                       "test_customer_gateway_name",
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

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeVpnConnectionsResponse
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
				res: &computing.DescribeVpnConnectionsResponse{
					DescribeVpnConnectionsOutput: &computing.DescribeVpnConnectionsOutput{
						VpnConnectionSet: []computing.VpnConnectionSet{
							{
								Type:                     nifcloud.String("L2TPv3 / IPsec"),
								VpnGatewayId:             nifcloud.String("test_vpn_gateway_id"),
								NiftyVpnGatewayName:      nifcloud.String("test_vpn_gateway_name"),
								CustomerGatewayId:        nifcloud.String("test_customer_gateway_id"),
								NiftyCustomerGatewayName: nifcloud.String("test_customer_gateway_name"),
								NiftyTunnel: &computing.NiftyTunnel{
									Type:            nifcloud.String("L2TPv3"),
									Mode:            nifcloud.String("Unmanaged"),
									Encapsulation:   nifcloud.String("UDP"),
									TunnelId:        nifcloud.String("1"),
									PeerTunnelId:    nifcloud.String("2"),
									SessionId:       nifcloud.String("1"),
									PeerSessionId:   nifcloud.String("2"),
									SourcePort:      nifcloud.String("7777"),
									DestinationPort: nifcloud.String("7777"),
								},
								NiftyIpsecConfiguration: &computing.NiftyIpsecConfiguration{
									EncryptionAlgorithm:                  nifcloud.String("AES256"),
									HashingAlgorithm:                     nifcloud.String("SHA256"),
									PreSharedKey:                         nifcloud.String("test_pre_shared_key"),
									InternetKeyExchange:                  nifcloud.String("IKEv2"),
									InternetKeyExchangeLifetime:          nifcloud.Int64(300),
									EncapsulatingSecurityPayloadLifetime: nifcloud.Int64(300),
									DiffieHellmanGroup:                   nifcloud.Int64(5),
									Mtu:                                  nifcloud.String("1000"),
								},
								NiftyVpnConnectionDescription: nifcloud.String("test_description"),
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
				res: &computing.DescribeVpnConnectionsResponse{
					DescribeVpnConnectionsOutput: &computing.DescribeVpnConnectionsOutput{
						VpnConnectionSet: []computing.VpnConnectionSet{},
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
