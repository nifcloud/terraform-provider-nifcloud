package nattable

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlattenProtocolTcp(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"nat_table_id": "test_nat_table_id",
		"snat": []interface{}{map[string]interface{}{
			"rule_number":                     "1",
			"description":                     "test_description",
			"protocol":                        "TCP",
			"source_address":                  "192.0.2.1",
			"source_port":                     80,
			"translation_port":                81,
			"outbound_interface_network_id":   "test_network_id",
			"outbound_interface_network_name": "test_network_name",
		}},
		"dnat": []interface{}{map[string]interface{}{
			"rule_number":                    "1",
			"description":                    "test_description",
			"protocol":                       "TCP",
			"destination_port":               80,
			"translation_address":            "192.0.2.1",
			"translation_port":               81,
			"inbound_interface_network_id":   "test_network_id",
			"inbound_interface_network_name": "test_network_name",
		}},
	})
	rd.SetId("test_nat_table_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.NiftyDescribeNatTablesResponse
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
				res: &computing.NiftyDescribeNatTablesResponse{
					NiftyDescribeNatTablesOutput: &computing.NiftyDescribeNatTablesOutput{
						NatTableSet: []computing.NatTableSet{
							{
								NatTableId: nifcloud.String("test_nat_table_id"),
								NatRuleSet: []computing.NatRuleSet{
									{
										NatType:     nifcloud.String("snat"),
										RuleNumber:  nifcloud.String("1"),
										Description: nifcloud.String("test_description"),
										Protocol:    nifcloud.String("TCP"),
										Source: &computing.Source{
											Address: nifcloud.String("192.0.2.1"),
											Port:    nifcloud.Int64(80),
										},
										Translation: &computing.Translation{
											Port: nifcloud.Int64(81),
										},
										OutboundInterface: &computing.OutboundInterface{
											NetworkId:   nifcloud.String("test_network_id"),
											NetworkName: nifcloud.String("test_network_bame"),
										},
									},
									{
										NatType:     nifcloud.String("dnat"),
										RuleNumber:  nifcloud.String("1"),
										Description: nifcloud.String("test_description"),
										Protocol:    nifcloud.String("TCP"),
										Destination: &computing.Destination{
											Port: nifcloud.Int64(80),
										},
										Translation: &computing.Translation{
											Address: nifcloud.String("192.0.2.1"),
											Port:    nifcloud.Int64(81),
										},
										InboundInterface: &computing.InboundInterface{
											NetworkId:   nifcloud.String("test_network_id"),
											NetworkName: nifcloud.String("test_network_bame"),
										},
									},
								},
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
				res: &computing.NiftyDescribeNatTablesResponse{
					NiftyDescribeNatTablesOutput: &computing.NiftyDescribeNatTablesOutput{
						NatTableSet: []computing.NatTableSet{},
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

func TestFlattenProtocolAll(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"nat_table_id": "test_nat_table_id",
		"snat": []interface{}{map[string]interface{}{
			"rule_number":                     "1",
			"description":                     "test_description",
			"protocol":                        "ALL",
			"source_address":                  "192.0.2.1",
			"outbound_interface_network_id":   "test_network_id",
			"outbound_interface_network_name": "test_network_name",
		}},
		"dnat": []interface{}{map[string]interface{}{
			"rule_number":                    "1",
			"description":                    "test_description",
			"protocol":                       "ALL",
			"translation_address":            "192.0.2.1",
			"inbound_interface_network_id":   "test_network_id",
			"inbound_interface_network_name": "test_network_name",
		}},
	})
	rd.SetId("test_nat_table_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.NiftyDescribeNatTablesResponse
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
				res: &computing.NiftyDescribeNatTablesResponse{
					NiftyDescribeNatTablesOutput: &computing.NiftyDescribeNatTablesOutput{
						NatTableSet: []computing.NatTableSet{
							{
								NatTableId: nifcloud.String("test_nat_table_id"),
								NatRuleSet: []computing.NatRuleSet{
									{
										NatType:     nifcloud.String("snat"),
										RuleNumber:  nifcloud.String("1"),
										Description: nifcloud.String("test_description"),
										Protocol:    nifcloud.String("ALL"),
										Source: &computing.Source{
											Address: nifcloud.String("192.0.2.1"),
										},
										OutboundInterface: &computing.OutboundInterface{
											NetworkId:   nifcloud.String("test_network_id"),
											NetworkName: nifcloud.String("test_network_bame"),
										},
									},
									{
										NatType:     nifcloud.String("dnat"),
										RuleNumber:  nifcloud.String("1"),
										Description: nifcloud.String("test_description"),
										Protocol:    nifcloud.String("ALL"),
										Translation: &computing.Translation{
											Address: nifcloud.String("192.0.2.1"),
										},
										InboundInterface: &computing.InboundInterface{
											NetworkId:   nifcloud.String("test_network_id"),
											NetworkName: nifcloud.String("test_network_bame"),
										},
									},
								},
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
				res: &computing.NiftyDescribeNatTablesResponse{
					NiftyDescribeNatTablesOutput: &computing.NiftyDescribeNatTablesOutput{
						NatTableSet: []computing.NatTableSet{},
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
