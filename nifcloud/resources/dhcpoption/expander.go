package dhcpoption

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandCreateDhcpOptionsInput(d *schema.ResourceData) *computing.CreateDhcpOptionsInput {
	var dhcpConfiguration []computing.RequestDhcpConfiguration
	if v := d.Get("default_router").(string); v != "" {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptionsDefaultRouter
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if v := d.Get("domain_name").(string); v != "" {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptionsDomainName
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if vs := d.Get("domain_name_servers").([]interface{}); len(vs) != 0 {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptionsDomainNameServers
		for _, v := range vs {
			dc.ListOfRequestValue = append(dc.ListOfRequestValue, v.(string))
		}
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if vs := d.Get("ntp_servers").([]interface{}); len(vs) != 0 {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptionsNtpServers
		for _, v := range vs {
			dc.ListOfRequestValue = append(dc.ListOfRequestValue, v.(string))
		}
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if vs := d.Get("netbios_name_servers").([]interface{}); len(vs) != 0 {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptionsNetbiosNameServers
		for _, v := range vs {
			dc.ListOfRequestValue = append(dc.ListOfRequestValue, v.(string))
		}
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if v := d.Get("netbios_node_type").(string); v != "" {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptionsNetbiosNodeType
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if v := d.Get("lease_time").(string); v != "" {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptionsLeaseTime
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
