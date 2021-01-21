package dhcpconfig

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandNiftyCreateDhcpConfigStaticMappingInput(d *schema.ResourceData, staticmapping map[string]interface{}) *computing.NiftyCreateDhcpStaticMappingInput {
	return &computing.NiftyCreateDhcpStaticMappingInput{
		DhcpConfigId: nifcloud.String(d.Id()),
		IpAddress:    nifcloud.String(staticmapping["static_mapping_ipaddress"].(string)),
		MacAddress:   nifcloud.String(staticmapping["static_mapping_macaddress"].(string)),
		Description:  nifcloud.String(staticmapping["static_mapping_description"].(string)),
	}
}

func expandNiftyCreateDhcpConfigIPAddressPoolInput(d *schema.ResourceData, ipaddresspool map[string]interface{}) *computing.NiftyCreateDhcpIpAddressPoolInput {
	return &computing.NiftyCreateDhcpIpAddressPoolInput{
		DhcpConfigId:   nifcloud.String(d.Id()),
		StartIpAddress: nifcloud.String(ipaddresspool["ipaddress_pool_start"].(string)),
		StopIpAddress:  nifcloud.String(ipaddresspool["ipaddress_pool_stop"].(string)),
		Description:    nifcloud.String(ipaddresspool["ipaddress_pool_description"].(string)),
	}
}

func expandNiftyDescribeDhcpConfigsInput(d *schema.ResourceData) *computing.NiftyDescribeDhcpConfigsInput {
	return &computing.NiftyDescribeDhcpConfigsInput{
		DhcpConfigId: []string{d.Id()},
	}
}

func expandNiftyDeleteDhcpConfigInput(d *schema.ResourceData) *computing.NiftyDeleteDhcpConfigInput {
	return &computing.NiftyDeleteDhcpConfigInput{
		DhcpConfigId: nifcloud.String(d.Id()),
	}
}

func expandNiftyDeleteDhcpConfigStaticMappingInput(d *schema.ResourceData, staticmapping map[string]interface{}) *computing.NiftyDeleteDhcpStaticMappingInput {
	return &computing.NiftyDeleteDhcpStaticMappingInput{
		DhcpConfigId: nifcloud.String(d.Id()),
		IpAddress:    nifcloud.String(staticmapping["static_mapping_ipaddress"].(string)),
		MacAddress:   nifcloud.String(staticmapping["static_mapping_macaddress"].(string)),
	}
}

func expandNiftyDeleteDhcpConfigIPAddressPoolInput(d *schema.ResourceData, ipaddresspool map[string]interface{}) *computing.NiftyDeleteDhcpIpAddressPoolInput {
	return &computing.NiftyDeleteDhcpIpAddressPoolInput{
		DhcpConfigId:   nifcloud.String(d.Id()),
		StartIpAddress: nifcloud.String(ipaddresspool["ipaddress_pool_start"].(string)),
		StopIpAddress:  nifcloud.String(ipaddresspool["ipaddress_pool_stop"].(string)),
	}
}
