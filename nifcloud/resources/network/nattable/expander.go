package nattable

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandNiftyCreateNatRuleInputForSnat(d *schema.ResourceData, snat map[string]interface{}) *computing.NiftyCreateNatRuleInput {
	return &computing.NiftyCreateNatRuleInput{
		NatTableId:  nifcloud.String(d.Id()),
		NatType:     types.NatTypeOfNiftyCreateNatRuleRequestSnat,
		RuleNumber:  nifcloud.String(snat["rule_number"].(string)),
		Description: nifcloud.String(snat["description"].(string)),
		Protocol:    types.ProtocolOfNiftyCreateNatRuleRequest(snat["protocol"].(string)),
		Source: &types.RequestSource{
			Address: nifcloud.String(snat["source_address"].(string)),
			Port:    nifcloud.Int32(int32(snat["source_port"].(int))),
		},
		Translation: &types.RequestTranslation{
			Port: nifcloud.Int32(int32(snat["translation_port"].(int))),
		},
		OutboundInterface: &types.RequestOutboundInterface{
			NetworkId:   nifcloud.String(snat["outbound_interface_network_id"].(string)),
			NetworkName: nifcloud.String(snat["outbound_interface_network_name"].(string)),
		},
	}
}

func expandNiftyCreateNatRuleInputForDnat(d *schema.ResourceData, dnat map[string]interface{}) *computing.NiftyCreateNatRuleInput {
	return &computing.NiftyCreateNatRuleInput{
		NatTableId:  nifcloud.String(d.Id()),
		NatType:     types.NatTypeOfNiftyCreateNatRuleRequestDnat,
		RuleNumber:  nifcloud.String(dnat["rule_number"].(string)),
		Description: nifcloud.String(dnat["description"].(string)),
		Protocol:    types.ProtocolOfNiftyCreateNatRuleRequest(dnat["protocol"].(string)),
		Destination: &types.RequestDestination{
			Port: nifcloud.Int32(int32(dnat["destination_port"].(int))),
		},
		Translation: &types.RequestTranslation{
			Address: nifcloud.String(dnat["translation_address"].(string)),
			Port:    nifcloud.Int32(int32(dnat["translation_port"].(int))),
		},
		InboundInterface: &types.RequestInboundInterface{
			NetworkId:   nifcloud.String(dnat["inbound_interface_network_id"].(string)),
			NetworkName: nifcloud.String(dnat["inbound_interface_network_name"].(string)),
		},
	}
}

func expandNiftyDescribeNatTablesInput(d *schema.ResourceData) *computing.NiftyDescribeNatTablesInput {
	return &computing.NiftyDescribeNatTablesInput{
		NatTableId: []string{d.Id()},
	}
}

func expandNiftyDeleteNatRuleInputForSnat(d *schema.ResourceData, snat map[string]interface{}) *computing.NiftyDeleteNatRuleInput {
	return &computing.NiftyDeleteNatRuleInput{
		NatTableId: nifcloud.String(d.Id()),
		NatType:    types.NatTypeOfNiftyDeleteNatRuleRequestSnat,
		RuleNumber: nifcloud.String(snat["rule_number"].(string)),
	}
}

func expandNiftyDeleteNatRuleInputForDnat(d *schema.ResourceData, dnat map[string]interface{}) *computing.NiftyDeleteNatRuleInput {
	return &computing.NiftyDeleteNatRuleInput{
		NatTableId: nifcloud.String(d.Id()),
		NatType:    types.NatTypeOfNiftyDeleteNatRuleRequestDnat,
		RuleNumber: nifcloud.String(dnat["rule_number"].(string)),
	}
}

func expandNiftyDeleteNatTableInput(d *schema.ResourceData) *computing.NiftyDeleteNatTableInput {
	return &computing.NiftyDeleteNatTableInput{
		NatTableId: nifcloud.String(d.Id()),
	}
}
