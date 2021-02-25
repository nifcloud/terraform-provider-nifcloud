package vpngateway

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeVpnGatewaysResponse) error {

	if res == nil || len(res.VpnGatewaySet) == 0 {
		d.SetId("")
		return nil
	}
	vpnGateway := res.VpnGatewaySet[0]

	if nifcloud.StringValue(vpnGateway.VpnGatewayId) != d.Id() {
		return fmt.Errorf("unable to find vpngateway within: %#v", res.VpnGatewaySet)
	}

	if err := d.Set("vpn_gateway_id", vpnGateway.VpnGatewayId); err != nil {
		return err
	}

	if err := d.Set("accounting_type", vpnGateway.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("nifty_vpn_gateway_description", vpnGateway.NiftyVpnGatewayDescription); err != nil {
		return err
	}

	if err := d.Set("nifty_vpn_gateway_name", vpnGateway.NiftyVpnGatewayName); err != nil {
		return err
	}

	if err := d.Set("nifty_vpn_gateway_type", vpnGateway.NiftyVpnGatewayType); err != nil {
		return err
	}

	if err := d.Set("availability_zone", vpnGateway.AvailabilityZone); err != nil {
		return err
	}

	for _, n := range vpnGateway.NetworkInterfaceSet {
		switch nifcloud.StringValue(n.NetworkId) {
		case "net-COMMON_GLOBAL", "net-COMMON_PRIVATE":
			// Can not use "net-COMMON_GLOBAL", "net-COMMON_PRIVATE"
		default:
			if row, ok := d.GetOk("network_id"); ok {
				if err := d.Set("network_id", row); err != nil {
					return err
				}
			}
			if row, ok := d.GetOk("network_name"); ok {
				if err := d.Set("network_name", row); err != nil {
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

	if err := d.Set("route_table_id", vpnGateway.RouteTableId); err != nil {
		return err
	}

	if err := d.Set("route_table_association_id", vpnGateway.RouteTableAssociationId); err != nil {
		return err
	}

	return nil
}
