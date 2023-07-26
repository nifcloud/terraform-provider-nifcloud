package remoteaccessvpngateway

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeRemoteAccessVpnGatewaysOutput) error {
	if res == nil || len(res.RemoteAccessVpnGatewaySet) == 0 {
		d.SetId("")
		return nil
	}

	remoteAccessVpnGateway := res.RemoteAccessVpnGatewaySet[0]

	if nifcloud.ToString(remoteAccessVpnGateway.RemoteAccessVpnGatewayId) != d.Id() {
		return fmt.Errorf("unable to find remote access vpn gateway within: %#v", res.RemoteAccessVpnGatewaySet)
	}

	if err := d.Set("remote_access_vpn_gateway_id", remoteAccessVpnGateway.RemoteAccessVpnGatewayId); err != nil {
		return err
	}

	if err := d.Set("name", remoteAccessVpnGateway.RemoteAccessVpnGatewayName); err != nil {
		return err
	}

	if err := d.Set("description", remoteAccessVpnGateway.Description); err != nil {
		return err
	}

	if err := d.Set("availability_zone", remoteAccessVpnGateway.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("accounting_type", remoteAccessVpnGateway.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("type", remoteAccessVpnGateway.RemoteAccessVpnGatewayType); err != nil {
		return err
	}

	if err := d.Set("pool_network_cidr", remoteAccessVpnGateway.PoolNetworkCidr); err != nil {
		return err
	}

	if err := d.Set("ca_certificate_id", remoteAccessVpnGateway.CaCertificateId); err != nil {
		return err
	}

	if err := d.Set("ssl_certificate_id", remoteAccessVpnGateway.SslCertificateId); err != nil {
		return err
	}

	var cipherSuite []string
	for _, c := range remoteAccessVpnGateway.CipherSuiteSet {
		cipherSuite = append(cipherSuite, nifcloud.ToString(c.CipherSuite))
	}

	if err := d.Set("cipher_suite", cipherSuite); err != nil {
		return err
	}

	var networkInterfaces []map[string]interface{}
	for _, n := range remoteAccessVpnGateway.NetworkInterfaceSet {
		switch nifcloud.ToString(n.NiftyNetworkId) {
		case "net-COMMON_GLOBAL":
			continue
		default:
			networkInterfaces = append(networkInterfaces, map[string]interface{}{
				"network_id": nifcloud.ToString(n.NiftyNetworkId),
				"ip_address": nifcloud.ToString(n.PrivateIpAddress),
			})
		}
	}

	if err := d.Set("network_interface", networkInterfaces); err != nil {
		return err
	}

	var users []map[string]interface{}
	for _, u := range remoteAccessVpnGateway.RemoteUserSet {
		var findElm map[string]interface{}
		for _, du := range d.Get("user").(*schema.Set).List() {
			elm := du.(map[string]interface{})
			if elm["name"] == nifcloud.ToString(u.UserName) {
				findElm = elm
				break
			}
		}

		user := make(map[string]interface{})
		user["name"] = nifcloud.ToString(u.UserName)
		user["description"] = nifcloud.ToString(u.Description)

		if findElm != nil {
			if findElm["password"] != nil && findElm["password"] != "" {
				user["password"] = findElm["password"]
			}
		}

		users = append(users, user)
	}

	if err := d.Set("user", users); err != nil {
		return err
	}

	return nil
}

func flattenClientConfig(d *schema.ResourceData, res *computing.DescribeRemoteAccessVpnGatewayClientConfigOutput) error {
	if err := d.Set("client_config", res.FileData); err != nil {
		return err
	}
	return nil
}
