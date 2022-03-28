package vpngateway

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeVpnGatewaysOutput) error {

	if res == nil || len(res.VpnGatewaySet) == 0 {
		d.SetId("")
		return nil
	}
	vpnGateway := res.VpnGatewaySet[0]

	if nifcloud.ToString(vpnGateway.VpnGatewayId) != d.Id() {
		return fmt.Errorf("unable to find vpngateway within: %#v", res.VpnGatewaySet)
	}

	if err := d.Set("vpn_gateway_id", vpnGateway.VpnGatewayId); err != nil {
		return err
	}

	if err := d.Set("accounting_type", vpnGateway.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("description", vpnGateway.NiftyVpnGatewayDescription); err != nil {
		return err
	}

	if err := d.Set("name", vpnGateway.NiftyVpnGatewayName); err != nil {
		return err
	}

	if err := d.Set("type", vpnGateway.NiftyVpnGatewayType); err != nil {
		return err
	}

	if err := d.Set("availability_zone", vpnGateway.AvailabilityZone); err != nil {
		return err
	}

	for _, n := range vpnGateway.NetworkInterfaceSet {
		switch nifcloud.ToString(n.NetworkId) {
		case "net-COMMON_GLOBAL":
			if err := d.Set("public_ip_address", n.IpAddress); err != nil {
				return err
			}
		default:
			if _, ok := d.GetOk("network_name"); ok {
				if err := d.Set("network_name", n.NetworkName); err != nil {
					return err
				}
			} else {
				if err := d.Set("network_id", n.NetworkId); err != nil {
					return err
				}
			}
			if err := d.Set("ip_address", n.IpAddress); err != nil {
				return err
			}
		}
	}

	if len(vpnGateway.GroupSet) > 0 {
		if err := d.Set("security_group", vpnGateway.GroupSet[0].GroupId); err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("route_table_id"); ok {
		if err := d.Set("route_table_id", vpnGateway.RouteTableId); err != nil {
			return err
		}
	}

	if err := d.Set("route_table_association_id", vpnGateway.RouteTableAssociationId); err != nil {
		return err
	}

	return nil
}
