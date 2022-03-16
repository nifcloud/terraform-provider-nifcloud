package nattable

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandNiftyCreateNatRuleInputForSnat(t *testing.T) {
	snat := map[string]interface{}{
		"rule_number":                     "1",
		"description":                     "test_description",
		"protocol":                        "TCP",
		"source_address":                  "192.0.2.1",
		"source_port":                     80,
		"translation_port":                81,
		"outbound_interface_network_id":   "test_network_id",
		"outbound_interface_network_name": "test_network_name",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"snat": []interface{}{snat},
	})
	rd.SetId("test_nat_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreateNatRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreateNatRuleInput{
				NatTableId:  nifcloud.String("test_nat_table_id"),
				NatType:     types.NatTypeOfNiftyCreateNatRuleRequestSnat,
				RuleNumber:  nifcloud.String("1"),
				Description: nifcloud.String("test_description"),
				Protocol:    types.ProtocolOfNiftyCreateNatRuleRequestTcp,
				Source: &types.RequestSource{
					Address: nifcloud.String("192.0.2.1"),
					Port:    nifcloud.Int32(80),
				},
				Translation: &types.RequestTranslation{
					Port: nifcloud.Int32(81),
				},
				OutboundInterface: &types.RequestOutboundInterface{
					NetworkId:   nifcloud.String("test_network_id"),
					NetworkName: nifcloud.String("test_network_name"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreateNatRuleInputForSnat(tt.args, snat)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyCreateNatRuleInputForDnat(t *testing.T) {
	dnat := map[string]interface{}{
		"rule_number":                    "1",
		"description":                    "test_description",
		"protocol":                       "TCP",
		"destination_port":               80,
		"translation_address":            "192.0.2.1",
		"translation_port":               81,
		"inbound_interface_network_id":   "test_network_id",
		"inbound_interface_network_name": "test_network_name",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"dnat": []interface{}{dnat},
	})
	rd.SetId("test_nat_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreateNatRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreateNatRuleInput{
				NatTableId:  nifcloud.String("test_nat_table_id"),
				NatType:     types.NatTypeOfNiftyCreateNatRuleRequestDnat,
				RuleNumber:  nifcloud.String("1"),
				Description: nifcloud.String("test_description"),
				Protocol:    types.ProtocolOfNiftyCreateNatRuleRequestTcp,
				Destination: &types.RequestDestination{
					Port: nifcloud.Int32(80),
				},
				Translation: &types.RequestTranslation{
					Address: nifcloud.String("192.0.2.1"),
					Port:    nifcloud.Int32(81),
				},
				InboundInterface: &types.RequestInboundInterface{
					NetworkId:   nifcloud.String("test_network_id"),
					NetworkName: nifcloud.String("test_network_name"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreateNatRuleInputForDnat(tt.args, dnat)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDescribeNatTablesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_nat_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDescribeNatTablesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDescribeNatTablesInput{
				NatTableId: []string{"test_nat_table_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDescribeNatTablesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeleteNatRuleInputForSnat(t *testing.T) {
	snat := map[string]interface{}{
		"rule_number": "1",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"snat": []interface{}{snat},
	})
	rd.SetId("test_nat_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteNatRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteNatRuleInput{
				NatTableId: nifcloud.String("test_nat_table_id"),
				NatType:    types.NatTypeOfNiftyDeleteNatRuleRequestSnat,
				RuleNumber: nifcloud.String("1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteNatRuleInputForSnat(tt.args, snat)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeleteNatRuleInputForDnat(t *testing.T) {
	dnat := map[string]interface{}{
		"rule_number": "1",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"dnat": []interface{}{dnat},
	})
	rd.SetId("test_nat_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteNatRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteNatRuleInput{
				NatTableId: nifcloud.String("test_nat_table_id"),
				NatType:    types.NatTypeOfNiftyDeleteNatRuleRequestDnat,
				RuleNumber: nifcloud.String("1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteNatRuleInputForDnat(tt.args, dnat)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeleteNatTableInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_nat_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteNatTableInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteNatTableInput{
				NatTableId: nifcloud.String("test_nat_table_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteNatTableInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
