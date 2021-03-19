package vpnconnection

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandCreateVpnConnectionInput(d *schema.ResourceData) *computing.CreateVpnConnectionInput {
	input := &computing.CreateVpnConnectionInput{
		Type: computing.TypeOfCreateVpnConnectionRequest(d.Get("type").(string)),
		NiftyIpsecConfiguration: &computing.RequestNiftyIpsecConfiguration{
			EncryptionAlgorithm:                  computing.EncryptionAlgorithmOfNiftyIpsecConfigurationForCreateVpnConnection(d.Get("ipsec_config_encryption_algorithm").(string)),
			HashAlgorithm:                        computing.HashAlgorithmOfNiftyIpsecConfigurationForCreateVpnConnection(d.Get("ipsec_config_hash_algorithm").(string)),
			PreSharedKey:                         nifcloud.String(d.Get("ipsec_config_pre_shared_key").(string)),
			InternetKeyExchange:                  computing.InternetKeyExchangeOfNiftyIpsecConfigurationForCreateVpnConnection(d.Get("ipsec_config_internet_key_exchange").(string)),
			InternetKeyExchangeLifetime:          nifcloud.Int64(int64(d.Get("ipsec_config_internet_key_exchange_lifetime").(int))),
			EncapsulatingSecurityPayloadLifetime: nifcloud.Int64(int64(d.Get("ipsec_config_encapsulating_security_payload_lifetime").(int))),
			DiffieHellmanGroup:                   nifcloud.Int64(int64(d.Get("ipsec_config_diffie_hellman_group").(int))),
		},
		NiftyVpnConnectionDescription: nifcloud.String(d.Get("description").(string)),
		Agreement:                     nifcloud.Bool(false),
	}

	if d.Get("type").(string) == "L2TPv3 / IPsec" {
		tunnel := computing.RequestNiftyTunnel{}

		tunnel.Type = computing.TypeOfNiftyTunnelForCreateVpnConnection(d.Get("tunnel_type").(string))

		if len(d.Get("tunnel_type").(string)) != 0 {
			tunnel.Mode = computing.ModeOfNiftyTunnelForCreateVpnConnection(d.Get("tunnel_mode").(string))
			tunnel.Encapsulation = computing.EncapsulationOfNiftyTunnelForCreateVpnConnection(d.Get("tunnel_encapsulation").(string))

			if d.Get("tunnel_mode").(string) == "Unmanaged" {
				tunnel.TunnelId = nifcloud.String(d.Get("tunnel_id").(string))
				tunnel.PeerTunnelId = nifcloud.String(d.Get("tunnel_peer_id").(string))
				tunnel.SessionId = nifcloud.String(d.Get("tunnel_session_id").(string))
				tunnel.PeerSessionId = nifcloud.String(d.Get("tunnel_peer_session_id").(string))

				if d.Get("tunnel_encapsulation").(string) == "UDP" {
					tunnel.SourcePort = nifcloud.String(d.Get("tunnel_source_port").(string))
					tunnel.DestinationPort = nifcloud.String(d.Get("tunnel_destination_port").(string))
				}
			}
		}
		input.NiftyTunnel = &tunnel
		input.NiftyVpnConnectionMtu = nifcloud.String(d.Get("mtu").(string))
	}

	if len(d.Get("vpn_gateway_id").(string)) != 0 {
		input.VpnGatewayId = nifcloud.String(d.Get("vpn_gateway_id").(string))
	}
	if len(d.Get("vpn_gateway_name").(string)) != 0 {
		input.NiftyVpnGatewayName = nifcloud.String(d.Get("vpn_gateway_name").(string))
	}
	if len(d.Get("customer_gateway_id").(string)) != 0 {
		input.CustomerGatewayId = nifcloud.String(d.Get("customer_gateway_id").(string))
	}
	if len(d.Get("customer_gateway_name").(string)) != 0 {
		input.NiftyCustomerGatewayName = nifcloud.String(d.Get("customer_gateway_name").(string))
	}

	return input
}

func expandDescribeVpnConnectionsInput(d *schema.ResourceData) *computing.DescribeVpnConnectionsInput {
	return &computing.DescribeVpnConnectionsInput{
		VpnConnectionId: []string{d.Id()},
	}
}

func expandDeleteVpnConnectionInput(d *schema.ResourceData) *computing.DeleteVpnConnectionInput {
	return &computing.DeleteVpnConnectionInput{
		VpnConnectionId: nifcloud.String(d.Id()),
		Agreement:       nifcloud.Bool(false),
	}
}
