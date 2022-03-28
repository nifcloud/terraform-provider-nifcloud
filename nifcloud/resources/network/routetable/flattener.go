package routetable

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeRouteTablesOutput) error {
	if res == nil || len(res.RouteTableSet) == 0 {
		d.SetId("")
		return nil
	}

	routeTable := res.RouteTableSet[0]

	if nifcloud.ToString(routeTable.RouteTableId) != d.Id() {
		return fmt.Errorf("unable to find route table within: %#v", res.RouteTableSet)
	}

	if err := d.Set("route_table_id", routeTable.RouteTableId); err != nil {
		return err
	}

	var routes []map[string]interface{}
	for _, r := range routeTable.RouteSet {
		// for vpn connection of IPsec VTI
		if nifcloud.ToString(r.Origin) == "EnableVgwRoutePropagation" {
			continue
		}

		route := map[string]interface{}{
			"ip_address": r.IpAddress,
			"cidr_block": r.DestinationCidrBlock,
		}

		var findElm map[string]interface{}
		for _, e := range d.Get("route").(*schema.Set).List() {
			elm := e.(map[string]interface{})

			if elm["cidr_block"] == nifcloud.ToString(r.DestinationCidrBlock) {
				findElm = elm
				break
			}
		}

		if findElm != nil {
			if findElm["network_id"] != "" {
				route["network_id"] = nifcloud.ToString(r.NetworkId)
			} else {
				route["network_name"] = nifcloud.ToString(r.NetworkName)
			}
		} else {
			route["network_id"] = nifcloud.ToString(r.NetworkId)
		}

		routes = append(routes, route)
	}

	if err := d.Set("route", routes); err != nil {
		return err
	}

	return nil
}
