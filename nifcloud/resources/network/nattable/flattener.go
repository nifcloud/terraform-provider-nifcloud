package nattable

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.NiftyDescribeNatTablesOutput) error {
	if res == nil || len(res.NatTableSet) == 0 {
		d.SetId("")
		return nil
	}

	natTable := res.NatTableSet[0]

	if nifcloud.ToString(natTable.NatTableId) != d.Id() {
		return fmt.Errorf("unable to find nat table within: %#v", res.NatTableSet)
	}

	if err := d.Set("nat_table_id", natTable.NatTableId); err != nil {
		return err
	}

	var snats []map[string]interface{}
	var dnats []map[string]interface{}

	for _, r := range natTable.NatRuleSet {

		if nifcloud.ToString(r.NatType) == "snat" {
			var findElm map[string]interface{}
			for _, e := range d.Get("snat").(*schema.Set).List() {
				elm := e.(map[string]interface{})

				if elm["rule_number"] == nifcloud.ToString(r.RuleNumber) {
					findElm = elm
					break
				}
			}

			snat := map[string]interface{}{
				"rule_number":    r.RuleNumber,
				"description":    r.Description,
				"protocol":       r.Protocol,
				"source_address": r.Source.Address,
			}

			if nifcloud.ToString(r.Protocol) != "ALL" && nifcloud.ToString(r.Protocol) != "ICMP" {
				snat["source_port"] = nifcloud.ToInt32(r.Source.Port)
				snat["translation_port"] = nifcloud.ToInt32(r.Translation.Port)
			}

			if findElm != nil {
				if findElm["outbound_interface_network_id"] != "" {
					snat["outbound_interface_network_id"] = nifcloud.ToString(r.OutboundInterface.NetworkId)
				}

				if findElm["outbound_interface_network_name"] != "" {
					snat["outbound_interface_network_name"] = nifcloud.ToString(r.OutboundInterface.NetworkName)
				}
			} else {
				snat["outbound_interface_network_id"] = nifcloud.ToString(r.OutboundInterface.NetworkId)
			}
			snats = append(snats, snat)

		} else if nifcloud.ToString(r.NatType) == "dnat" {
			var findElm map[string]interface{}
			for _, e := range d.Get("dnat").(*schema.Set).List() {
				elm := e.(map[string]interface{})

				if elm["rule_number"] == nifcloud.ToString(r.RuleNumber) {
					findElm = elm
					break
				}
			}

			dnat := map[string]interface{}{
				"rule_number":         r.RuleNumber,
				"description":         r.Description,
				"protocol":            r.Protocol,
				"translation_address": r.Translation.Address,
			}

			if nifcloud.ToString(r.Protocol) != "ALL" && nifcloud.ToString(r.Protocol) != "ICMP" {
				dnat["destination_port"] = nifcloud.ToInt32(r.Destination.Port)
				dnat["translation_port"] = nifcloud.ToInt32(r.Translation.Port)
			}

			if findElm != nil {
				if findElm["inbound_interface_network_id"] != "" {
					dnat["inbound_interface_network_id"] = nifcloud.ToString(r.InboundInterface.NetworkId)
				}

				if findElm["inbound_interface_network_name"] != "" {
					dnat["inbound_interface_network_name"] = nifcloud.ToString(r.InboundInterface.NetworkName)
				}
			} else {
				dnat["inbound_interface_network_id"] = nifcloud.ToString(r.InboundInterface.NetworkId)
			}
			dnats = append(dnats, dnat)
		}
	}

	if err := d.Set("snat", snats); err != nil {
		return err
	}

	if err := d.Set("dnat", dnats); err != nil {
		return err
	}

	return nil
}
