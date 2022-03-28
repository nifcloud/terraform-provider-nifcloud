package dhcpoption

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandCreateDhcpOptionsInput(d *schema.ResourceData) *computing.CreateDhcpOptionsInput {
	var dhcpConfiguration []types.RequestDhcpConfiguration
	if v := d.Get("default_router").(string); v != "" {
		dc := types.RequestDhcpConfiguration{}
		dc.Key = types.KeyOfDhcpConfigurationForCreateDhcpOptionsDefaultRouter
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if v := d.Get("domain_name").(string); v != "" {
		dc := types.RequestDhcpConfiguration{}
		dc.Key = types.KeyOfDhcpConfigurationForCreateDhcpOptionsDomainName
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if vs := d.Get("domain_name_servers").(*schema.Set).List(); len(vs) != 0 {
		dc := types.RequestDhcpConfiguration{}
		dc.Key = types.KeyOfDhcpConfigurationForCreateDhcpOptionsDomainNameServers
		for _, v := range vs {
			dc.ListOfRequestValue = append(dc.ListOfRequestValue, v.(string))
		}
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if vs := d.Get("ntp_servers").(*schema.Set).List(); len(vs) != 0 {
		dc := types.RequestDhcpConfiguration{}
		dc.Key = types.KeyOfDhcpConfigurationForCreateDhcpOptionsNtpServers
		for _, v := range vs {
			dc.ListOfRequestValue = append(dc.ListOfRequestValue, v.(string))
		}
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if vs := d.Get("netbios_name_servers").(*schema.Set).List(); len(vs) != 0 {
		dc := types.RequestDhcpConfiguration{}
		dc.Key = types.KeyOfDhcpConfigurationForCreateDhcpOptionsNetbiosNameServers
		for _, v := range vs {
			dc.ListOfRequestValue = append(dc.ListOfRequestValue, v.(string))
		}
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if v := d.Get("netbios_node_type").(string); v != "" {
		dc := types.RequestDhcpConfiguration{}
		dc.Key = types.KeyOfDhcpConfigurationForCreateDhcpOptionsNetbiosNodeType
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if v := d.Get("lease_time").(string); v != "" {
		dc := types.RequestDhcpConfiguration{}
		dc.Key = types.KeyOfDhcpConfigurationForCreateDhcpOptionsLeaseTime
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	return &computing.CreateDhcpOptionsInput{
		DhcpConfiguration: dhcpConfiguration,
	}
}

func expandDescribeDhcpOptionsInput(d *schema.ResourceData) *computing.DescribeDhcpOptionsInput {
	return &computing.DescribeDhcpOptionsInput{
		DhcpOptionsId: []string{d.Id()},
	}
}

func expandDeleteDhcpOptionsInput(d *schema.ResourceData) *computing.DeleteDhcpOptionsInput {
	return &computing.DeleteDhcpOptionsInput{
		DhcpOptionsId: nifcloud.String(d.Id()),
	}
}
