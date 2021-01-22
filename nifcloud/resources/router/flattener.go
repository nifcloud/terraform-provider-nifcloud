package router

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.NiftyDescribeRoutersResponse) error {
	if res == nil || len(res.RouterSet) == 0 {
		d.SetId("")
		return nil
	}

	router := res.RouterSet[0]

	if nifcloud.StringValue(router.RouterId) != d.Id() {
		return fmt.Errorf("unable to find router within: %#v", res.RouterSet)
	}

	if err := d.Set("router_id", router.RouterId); err != nil {
		return err
	}

	if err := d.Set("name", router.RouterName); err != nil {
		return err
	}

	if err := d.Set("description", router.Description); err != nil {
		return err
	}

	if err := d.Set("availability_zone", router.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("route_table_id", router.RouteTableId); err != nil {
		return err
	}

	if err := d.Set("route_table_association_id", router.RouteTableAssociationId); err != nil {
		return err
	}

	if err := d.Set("nat_table_id", router.NatTableId); err != nil {
		return err
	}

	if err := d.Set("nat_table_association_id", router.NatTableAssociationId); err != nil {
		return err
	}

	if len(router.GroupSet) > 0 {
		if err := d.Set("security_group", router.GroupSet[0].GroupId); err != nil {
			return err
		}
	}

	if err := d.Set("accounting_type", router.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("type", router.Type); err != nil {
		return err
	}

	// sort network interfaces set because API returns unstable set.
	sort.Slice(router.NetworkInterfaceSet, func(i, j int) bool {
		return nifcloud.StringValue(router.NetworkInterfaceSet[i].NetworkId) < nifcloud.StringValue(router.NetworkInterfaceSet[j].NetworkId)
	})

	var networkInterfaces []map[string]interface{}
	for _, n := range router.NetworkInterfaceSet {
		ni := make(map[string]interface{})
		switch nifcloud.StringValue(n.NetworkId) {
		case "net-COMMON_GLOBAL", "net-COMMON_PRIVATE":
			ni["network_id"] = nifcloud.StringValue(n.NetworkId)
		default:
			var findElm map[string]interface{}
			for _, dn := range d.Get("network_interface").(*schema.Set).List() {
				elm := dn.(map[string]interface{})

				if elm["network_id"] == nifcloud.StringValue(n.NetworkId) {
					findElm = elm
					break
				}

				if elm["network_name"] == nifcloud.StringValue(n.NetworkName) {
					findElm = elm
					break
				}
			}

			if findElm != nil {
				if findElm["ip_address"] != nil && findElm["ip_address"] != "" {
					ni["ip_address"] = nifcloud.StringValue(n.IpAddress)
				}

				if findElm["network_id"] != nil && findElm["network_id"] != "" {
					ni["network_id"] = nifcloud.StringValue(n.NetworkId)
				}

				if findElm["network_name"] != nil && findElm["network_name"] != "" {
					ni["network_name"] = nifcloud.StringValue(n.NetworkName)
				}

				if findElm["dhcp"] != nil {
					ni["dhcp"] = nifcloud.BoolValue(n.Dhcp)
				}

				if findElm["dhcp_options_id"] != nil && findElm["dhcp_options_id"] != "" {
					ni["dhcp_options_id"] = nifcloud.StringValue(n.DhcpOptionsId)
				}

				if findElm["dhcp_config_id"] != nil && findElm["dhcp_config_id"] != "" {
					ni["dhcp_config_id"] = nifcloud.StringValue(n.DhcpConfigId)
				}
			} else {
				ni["network_id"] = nifcloud.StringValue(n.NetworkId)
				ni["network_name"] = nifcloud.StringValue(n.NetworkName)
				ni["ip_address"] = nifcloud.StringValue(n.IpAddress)
				ni["dhcp"] = nifcloud.BoolValue(n.Dhcp)
				ni["dhcp_options_id"] = nifcloud.StringValue(n.DhcpOptionsId)
				ni["dhcp_config_id"] = nifcloud.StringValue(n.DhcpConfigId)
			}
		}
		networkInterfaces = append(networkInterfaces, ni)
	}

	if err := d.Set("network_interface", networkInterfaces); err != nil {
		return err
	}

	return nil
}
