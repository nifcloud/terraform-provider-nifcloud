package dhcpoption

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
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

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeDhcpOptionsResponse
		d   *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want *schema.ResourceData
	}{
		{
			name: "flattens the response",
			args: args{
				d: rd,
				res: &computing.DescribeDhcpOptionsResponse{
					DescribeDhcpOptionsOutput: &computing.DescribeDhcpOptionsOutput{
						DhcpOptionsSet: []computing.DhcpOptionsSet{
							{
								DhcpConfigurationSet: []computing.DhcpConfigurationSet{
									{
										Key: nifcloud.String("default-router"),
										ValueSet: []computing.ValueSet{
											{
												Value: nifcloud.String("test_default_router"),
											},
										},
									},
									{
										Key: nifcloud.String("domain-name"),
										ValueSet: []computing.ValueSet{
											{
												Value: nifcloud.String("test_domain_name"),
											},
										},
									},
									{
										Key: nifcloud.String("domain-name-servers"),
										ValueSet: []computing.ValueSet{
											{
												Value: nifcloud.String("test_domain_name_servers1"),
											},
											{
												Value: nifcloud.String("test_domain_name_servers2"),
											},
										},
									},
									{
										Key: nifcloud.String("ntp-servers"),
										ValueSet: []computing.ValueSet{
											{
												Value: nifcloud.String("test_ntp_servers"),
											},
										},
									},
									{
										Key: nifcloud.String("netbios-name-servers"),
										ValueSet: []computing.ValueSet{
											{
												Value: nifcloud.String("test_netbios_name_servers1"),
											},
											{
												Value: nifcloud.String("test_netbios_name_servers2"),
											},
										},
									},
									{
										Key: nifcloud.String("netbios-node-type"),
										ValueSet: []computing.ValueSet{
											{
												Value: nifcloud.String("test_netbios_node_type"),
											},
										},
									},
									{
										Key: nifcloud.String("lease-time"),
										ValueSet: []computing.ValueSet{
											{
												Value: nifcloud.String("test_lease_time"),
											},
										},
									},
								},
								DhcpOptionsId: nifcloud.String("test_dhcp_option_id"),
							},
						},
					},
				},
			},
			want: rd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d: wantNotFoundRd,
				res: &computing.DescribeDhcpOptionsResponse{
					DescribeDhcpOptionsOutput: &computing.DescribeDhcpOptionsOutput{
						DhcpOptionsSet: []computing.DhcpOptionsSet{},
					},
				},
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			assert.NoError(t, err)

			if tt.args.res == nil {
				return
			}

			wantState := tt.want.State()
			if wantState == nil {
				tt.want.SetId("some")
				wantState = tt.want.State()
			}

			gotState := tt.args.d.State()
			if gotState == nil {
				tt.args.d.SetId("some")
				gotState = tt.args.d.State()
			}

			assert.Equal(t, wantState.Attributes, gotState.Attributes)
		})
	}
}
