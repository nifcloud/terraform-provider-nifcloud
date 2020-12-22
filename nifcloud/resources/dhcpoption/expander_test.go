package dhcpoption

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateDhcpOptionsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"default_router":       "test_default_router",
		"domain_name":          "test_domain_name",
		"domain_name_servers":  []interface{}{"test_domain_name_servers1", "test_domain_name_servers2"},
		"ntp_servers":          []interface{}{"test_ntp_servers"},
		"netbios_name_servers": []interface{}{"test_netbios_name_servers1", "test_netbios_name_servers2"},
		"netbios_node_type":    "test_netbios_node_type",
		"lease_time":           "test_lease_time",
		"dhcp_option_id":       "test_dhcp_option_id",
	})
	rd.SetId("test_dhcp_option_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateDhcpOptionsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateDhcpOptionsInput{
				DhcpConfiguration: []computing.RequestDhcpConfiguration{
					{
						Key:                computing.KeyOfDhcpConfigurationForCreateDhcpOptions("default-router"),
						ListOfRequestValue: []string{"test_default_router"},
					},
					{
						Key:                computing.KeyOfDhcpConfigurationForCreateDhcpOptions("domain-name"),
						ListOfRequestValue: []string{"test_domain_name"},
					},
					{
						Key:                computing.KeyOfDhcpConfigurationForCreateDhcpOptions("domain-name-servers"),
						ListOfRequestValue: []string{"test_domain_name_servers1", "test_domain_name_servers2"},
					},
					{
						Key:                computing.KeyOfDhcpConfigurationForCreateDhcpOptions("ntp-servers"),
						ListOfRequestValue: []string{"test_ntp_servers"},
					},
					{
						Key:                computing.KeyOfDhcpConfigurationForCreateDhcpOptions("netbios-name-servers"),
						ListOfRequestValue: []string{"test_netbios_name_servers1", "test_netbios_name_servers2"},
					},
					{
						Key:                computing.KeyOfDhcpConfigurationForCreateDhcpOptions("netbios-node-type"),
						ListOfRequestValue: []string{"test_netbios_node_type"},
					},
					{
						Key:                computing.KeyOfDhcpConfigurationForCreateDhcpOptions("lease-time"),
						ListOfRequestValue: []string{"test_lease_time"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateDhcpOptionsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeDhcpOptionsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"dhce_option_id": "test_dhcp_option_id",
	})
	rd.SetId("test_dhcp_option_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeDhcpOptionsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeDhcpOptionsInput{
				DhcpOptionsId: []string{"test_dhcp_option_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeDhcpOptionsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteDhcpOptionsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"dhce_option_id": "test_dhcp_option_id",
	})
	rd.SetId("test_dhcp_option_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteDhcpOptionsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteDhcpOptionsInput{
				DhcpOptionsId: nifcloud.String("test_dhcp_option_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteDhcpOptionsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
