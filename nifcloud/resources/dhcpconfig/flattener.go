package dhcpconfig

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.NiftyDescribeDhcpConfigsResponse) error {
	if res == nil || len(res.DhcpConfigsSet) == 0 {
		d.SetId("")
		return nil
	}

	dhcpConfig := res.DhcpConfigsSet[0]

	if nifcloud.StringValue(dhcpConfig.DhcpConfigId) != d.Id() {
		return fmt.Errorf("unable to find dhcp config within: %#v", res.DhcpConfigsSet)
	}

	if err := d.Set("dhcp_config_id", dhcpConfig.DhcpConfigId); err != nil {
		return err
	}

	var staticmapping []map[string]interface{}
	var ipaddresspool []map[string]interface{}

	for _, s := range dhcpConfig.StaticMappingsSet {
		sm := map[string]interface{}{
			"static_mapping_ipaddress":   s.IpAddress,
			"static_mapping_macaddress":  s.MacAddress,
			"static_mapping_description": s.Description,
		}

		staticmapping = append(staticmapping, sm)
	}

	for _, i := range dhcpConfig.IpAddressPoolsSet {
		ip := map[string]interface{}{
			"ipaddress_pool_start":       i.StartIpAddress,
			"ipaddress_pool_stop":        i.StopIpAddress,
			"ipaddress_pool_description": i.Description,
		}

		ipaddresspool = append(ipaddresspool, ip)
	}

	if err := d.Set("static_mapping", staticmapping); err != nil {
		return err
	}

	if err := d.Set("ipaddress_pool", ipaddresspool); err != nil {
		return err
	}

	return nil
}
