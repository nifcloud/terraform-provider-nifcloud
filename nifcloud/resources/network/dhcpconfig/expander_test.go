package dhcpconfig

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandNiftyCreateDhcpConfigStaticMappingInput(t *testing.T) {
	staticmapping := map[string]interface{}{
		"static_mapping_ipaddress":   "192.168.1.10",
		"static_mapping_macaddress":  "00:00:5e:00:53:00",
		"static_mapping_description": "test_description",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"static_mapping": []interface{}{staticmapping},
	})
	rd.SetId("test_dhcp_config_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreateDhcpStaticMappingInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreateDhcpStaticMappingInput{
				DhcpConfigId: nifcloud.String("test_dhcp_config_id"),
				IpAddress:    nifcloud.String("192.168.1.10"),
				MacAddress:   nifcloud.String("00:00:5e:00:53:00"),
				Description:  nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreateDhcpConfigStaticMappingInput(tt.args, staticmapping)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyCreateDhcpConfigIpAddressPoolInput(t *testing.T) {
	ipaddresspool := map[string]interface{}{
		"ipaddress_pool_start":       "192.168.2.1",
		"ipaddress_pool_stop":        "192.168.2.100",
		"ipaddress_pool_description": "test_description",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ipaddress_pool": []interface{}{ipaddresspool},
	})
	rd.SetId("test_dhcp_config_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreateDhcpIpAddressPoolInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreateDhcpIpAddressPoolInput{
				DhcpConfigId:   nifcloud.String("test_dhcp_config_id"),
				StartIpAddress: nifcloud.String("192.168.2.1"),
				StopIpAddress:  nifcloud.String("192.168.2.100"),
				Description:    nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreateDhcpConfigIPAddressPoolInput(tt.args, ipaddresspool)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDescribeDhcpConfgsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_dhcp_config_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDescribeDhcpConfigsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDescribeDhcpConfigsInput{
				DhcpConfigId: []string{"test_dhcp_config_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDescribeDhcpConfigsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeleteDhcpConfigsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_dhcp_config_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteDhcpConfigInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteDhcpConfigInput{
				DhcpConfigId: nifcloud.String("test_dhcp_config_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteDhcpConfigInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeleteDhcpConfigStaticMappingInput(t *testing.T) {
	staticmapping := map[string]interface{}{
		"dhcp_config_id":            "test_dhcp_config_id",
		"static_mapping_ipaddress":  "192.168.1.10",
		"static_mapping_macaddress": "00:00:5e:00:53:00",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"static_mapping": []interface{}{staticmapping},
	})
	rd.SetId("test_dhcp_config_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteDhcpStaticMappingInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteDhcpStaticMappingInput{
				DhcpConfigId: nifcloud.String("test_dhcp_config_id"),
				IpAddress:    nifcloud.String("192.168.1.10"),
				MacAddress:   nifcloud.String("00:00:5e:00:53:00"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteDhcpConfigStaticMappingInput(tt.args, staticmapping)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeleteDhcpConfigIpAddressPoolInput(t *testing.T) {
	ipaddresspool := map[string]interface{}{
		"dhcp_config_id":       "test_dhcp_config_id",
		"ipaddress_pool_start": "192.168.2.1",
		"ipaddress_pool_stop":  "192.168.2.100",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ipaddress_pool": []interface{}{ipaddresspool},
	})
	rd.SetId("test_dhcp_config_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteDhcpIpAddressPoolInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteDhcpIpAddressPoolInput{
				DhcpConfigId:   nifcloud.String("test_dhcp_config_id"),
				StartIpAddress: nifcloud.String("192.168.2.1"),
				StopIpAddress:  nifcloud.String("192.168.2.100"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteDhcpConfigIPAddressPoolInput(tt.args, ipaddresspool)
			assert.Equal(t, tt.want, got)
		})
	}
}
