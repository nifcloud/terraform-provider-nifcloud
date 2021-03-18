package nattable

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandNiftyCreateNatRuleInputForSnat(d *schema.ResourceData, snat map[string]interface{}) *computing.NiftyCreateNatRuleInput {
	return &computing.NiftyCreateNatRuleInput{
		NatTableId:  nifcloud.String(d.Id()),
		NatType:     computing.NatTypeOfNiftyCreateNatRuleRequestSnat,
		RuleNumber:  nifcloud.String(snat["rule_number"].(string)),
		Description: nifcloud.String(snat["description"].(string)),
		Protocol:    computing.ProtocolOfNiftyCreateNatRuleRequest(snat["protocol"].(string)),
		Source: &computing.RequestSource{
			Address: nifcloud.String(snat["source_address"].(string)),
			Port:    nifcloud.Int64(int64(snat["source_port"].(int))),
		},
		Translation: &computing.RequestTranslation{
			Port: nifcloud.Int64(int64(snat["translation_port"].(int))),
		},
		OutboundInterface: &computing.RequestOutboundInterface{
			NetworkId:   nifcloud.String(snat["outbound_interface_network_id"].(string)),
			NetworkName: nifcloud.String(snat["outbound_interface_network_name"].(string)),
		},
	}
}

func expandNiftyCreateNatRuleInputForDnat(d *schema.ResourceData, dnat map[string]interface{}) *computing.NiftyCreateNatRuleInput {
	return &computing.NiftyCreateNatRuleInput{
		NatTableId:  nifcloud.String(d.Id()),
		NatType:     computing.NatTypeOfNiftyCreateNatRuleRequestDnat,
		RuleNumber:  nifcloud.String(dnat["rule_number"].(string)),
		Description: nifcloud.String(dnat["description"].(string)),
		Protocol:    computing.ProtocolOfNiftyCreateNatRuleRequest(dnat["protocol"].(string)),
		Destination: &computing.RequestDestination{
			Port: nifcloud.Int64(int64(dnat["destination_port"].(int))),
		},
		Translation: &computing.RequestTranslation{
			Address: nifcloud.String(dnat["translation_address"].(string)),
			Port:    nifcloud.Int64(int64(dnat["translation_port"].(int))),
		},
		InboundInterface: &computing.RequestInboundInterface{
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
		NatType:    computing.NatTypeOfNiftyDeleteNatRuleRequestSnat,
		RuleNumber: nifcloud.String(snat["rule_number"].(string)),
	}
}

func expandNiftyDeleteNatRuleInputForDnat(d *schema.ResourceData, dnat map[string]interface{}) *computing.NiftyDeleteNatRuleInput {
	return &computing.NiftyDeleteNatRuleInput{
		NatTableId: nifcloud.String(d.Id()),
		NatType:    computing.NatTypeOfNiftyDeleteNatRuleRequestDnat,
		RuleNumber: nifcloud.String(dnat["rule_number"].(string)),
	}
}

func expandNiftyDeleteNatTableInput(d *schema.ResourceData) *computing.NiftyDeleteNatTableInput {
	return &computing.NiftyDeleteNatTableInput{
		NatTableId: nifcloud.String(d.Id()),
	}
}
