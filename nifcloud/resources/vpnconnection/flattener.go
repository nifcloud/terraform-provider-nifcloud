package vpnconnection

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeVpnConnectionsResponse) error {
	if res == nil || len(res.VpnConnectionSet) == 0 {
		d.SetId("")
		return nil
	}

	vpnConnection := res.VpnConnectionSet[0]

	if nifcloud.StringValue(vpnConnection.VpnConnectionId) != d.Id() {
		return fmt.Errorf("unable to find vpn connection within: %#v", res.VpnConnectionSet)
	}

	if err := d.Set("vpn_connection_id", vpnConnection.VpnConnectionId); err != nil {
		return err
	}

	if err := d.Set("type", vpnConnection.Type); err != nil {
		return err
	}

	if _, ok := d.GetOk("vpn_gateway_id"); ok {
		if err := d.Set("vpn_gateway_id", vpnConnection.VpnGatewayId); err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("vpn_gateway_name"); ok {
		if err := d.Set("vpn_gateway_name", vpnConnection.NiftyVpnGatewayName); err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("customer_gateway_id"); ok {
		if err := d.Set("customer_gateway_id", vpnConnection.CustomerGatewayId); err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("customer_gateway_name"); ok {
		if err := d.Set("customer_gateway_name", vpnConnection.NiftyCustomerGatewayName); err != nil {
			return err
		}
	}

	if vpnConnection.NiftyTunnel != nil {
		if err := d.Set("tunnel_type", vpnConnection.NiftyTunnel.Type); err != nil {
			return err
		}

		if err := d.Set("tunnel_mode", vpnConnection.NiftyTunnel.Mode); err != nil {
			return err
		}

		if err := d.Set("tunnel_encapsulation", vpnConnection.NiftyTunnel.Encapsulation); err != nil {
			return err
		}

		if err := d.Set("mtu", vpnConnection.NiftyIpsecConfiguration.Mtu); err != nil {
			return err
		}

		if _, ok := d.GetOk("tunnel_id"); ok {
			if err := d.Set("tunnel_id", vpnConnection.NiftyTunnel.TunnelId); err != nil {
				return err
			}
		}

		if _, ok := d.GetOk("tunnel_peer_id"); ok {
			if err := d.Set("tunnel_peer_id", vpnConnection.NiftyTunnel.PeerTunnelId); err != nil {
				return err
			}
		}

		if _, ok := d.GetOk("tunnel_session_id"); ok {
			if err := d.Set("tunnel_session_id", vpnConnection.NiftyTunnel.SessionId); err != nil {
				return err
			}
		}

		if _, ok := d.GetOk("tunnel_peer_session_id"); ok {
			if err := d.Set("tunnel_peer_session_id", vpnConnection.NiftyTunnel.PeerSessionId); err != nil {
				return err
			}
		}

		if _, ok := d.GetOk("tunnel_source_port"); ok {
			if err := d.Set("tunnel_source_port", vpnConnection.NiftyTunnel.SourcePort); err != nil {
				return err
			}
		}

		if _, ok := d.GetOk("tunnel_destination_port"); ok {
			if err := d.Set("tunnel_destination_port", vpnConnection.NiftyTunnel.DestinationPort); err != nil {
				return err
			}
		}
	}

	if err := d.Set("ipsec_config_encryption_algorithm", vpnConnection.NiftyIpsecConfiguration.EncryptionAlgorithm); err != nil {
		return err
	}

	if err := d.Set("ipsec_config_hash_algorithm", vpnConnection.NiftyIpsecConfiguration.HashingAlgorithm); err != nil {
		return err
	}

	if err := d.Set("ipsec_config_hash_algorithm", vpnConnection.NiftyIpsecConfiguration.HashingAlgorithm); err != nil {
		return err
	}

	if err := d.Set("ipsec_config_pre_shared_key", vpnConnection.NiftyIpsecConfiguration.PreSharedKey); err != nil {
		return err
	}

	if err := d.Set("ipsec_config_internet_key_exchange", vpnConnection.NiftyIpsecConfiguration.InternetKeyExchange); err != nil {
		return err
	}

	if err := d.Set("ipsec_config_internet_key_exchange_lifetime", vpnConnection.NiftyIpsecConfiguration.InternetKeyExchangeLifetime); err != nil {
		return err
	}

	if err := d.Set("ipsec_config_encapsulating_security_payload_lifetime", vpnConnection.NiftyIpsecConfiguration.EncapsulatingSecurityPayloadLifetime); err != nil {
		return err
	}

	if err := d.Set("ipsec_config_diffie_hellman_group", vpnConnection.NiftyIpsecConfiguration.DiffieHellmanGroup); err != nil {
		return err
	}

	if err := d.Set("description", vpnConnection.NiftyVpnConnectionDescription); err != nil {
		return err
	}

	return nil
}
