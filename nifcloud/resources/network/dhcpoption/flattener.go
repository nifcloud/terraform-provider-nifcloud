package dhcpoption

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeDhcpOptionsOutput) error {
	if res == nil || len(res.DhcpOptionsSet) == 0 || len(res.DhcpOptionsSet[0].DhcpConfigurationSet) == 0 {
		d.SetId("")
		return nil
	}

	dhcpOption := res.DhcpOptionsSet[0]

	if nifcloud.ToString(dhcpOption.DhcpOptionsId) != d.Id() {
		return fmt.Errorf("unable to find dhcp option within: %#v", res.DhcpOptionsSet)
	}

	for _, dcf := range dhcpOption.DhcpConfigurationSet {
		if nifcloud.ToString(dcf.Key) == "default-router" {
			if err := d.Set("default_router", dcf.ValueSet[0].Value); err != nil {
				return err
			}
		} else if nifcloud.ToString(dcf.Key) == "domain-name" {
			if err := d.Set("domain_name", dcf.ValueSet[0].Value); err != nil {
				return err
			}
		} else if nifcloud.ToString(dcf.Key) == "domain-name-servers" {
			domainNameServers := make([]string, len(dcf.ValueSet))
			for i, val := range dcf.ValueSet {
				domainNameServers[i] = nifcloud.ToString(val.Value)
			}
			if err := d.Set("domain_name_servers", domainNameServers); err != nil {
				return err
			}
		} else if nifcloud.ToString(dcf.Key) == "ntp-servers" {
			ntpServers := make([]string, len(dcf.ValueSet))
			for i, val := range dcf.ValueSet {
				ntpServers[i] = nifcloud.ToString(val.Value)
			}
			if err := d.Set("ntp_servers", ntpServers); err != nil {
				return err
			}
		} else if nifcloud.ToString(dcf.Key) == "netbios-name-servers" {
			netbionsServers := make([]string, len(dcf.ValueSet))
			for i, val := range dcf.ValueSet {
				netbionsServers[i] = nifcloud.ToString(val.Value)
			}
			if err := d.Set("netbios_name_servers", netbionsServers); err != nil {
				return err
			}
		} else if nifcloud.ToString(dcf.Key) == "netbios-node-type" {
			if err := d.Set("netbios_node_type", dcf.ValueSet[0].Value); err != nil {
				return err
			}
		} else if nifcloud.ToString(dcf.Key) == "lease-time" {
			if err := d.Set("lease_time", dcf.ValueSet[0].Value); err != nil {
				return err
			}
		}
	}

	if err := d.Set("dhcp_option_id", dhcpOption.DhcpOptionsId); err != nil {
		return err
	}
	return nil
}
