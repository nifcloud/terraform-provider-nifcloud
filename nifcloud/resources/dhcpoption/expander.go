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
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptions("default-router")
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if v := d.Get("domain_name").(string); v != "" {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptions("domain-name")
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if vs := d.Get("domain_name_servers").([]interface{}); len(vs) != 0 {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptions("domain-name-servers")
		for _, v := range vs {
			dc.ListOfRequestValue = append(dc.ListOfRequestValue, v.(string))
		}
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if vs := d.Get("ntp_servers").([]interface{}); len(vs) != 0 {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptions("ntp-servers")
		for _, v := range vs {
			dc.ListOfRequestValue = append(dc.ListOfRequestValue, v.(string))
		}
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if vs := d.Get("netbios_name_servers").([]interface{}); len(vs) != 0 {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptions("netbios-name-servers")
		for _, v := range vs {
			dc.ListOfRequestValue = append(dc.ListOfRequestValue, v.(string))
		}
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if v := d.Get("netbios_node_type").(string); v != "" {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptions("netbios-node-type")
		dc.ListOfRequestValue = append(dc.ListOfRequestValue, v)
		dhcpConfiguration = append(dhcpConfiguration, dc)
	}
	if v := d.Get("lease_time").(string); v != "" {
		dc := computing.RequestDhcpConfiguration{}
		dc.Key = computing.KeyOfDhcpConfigurationForCreateDhcpOptions("lease-time")
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
