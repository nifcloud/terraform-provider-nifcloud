package devopsfirewallgroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func updateFirewallGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	if d.IsNewResource() {
		if err := waitUntilFirewallGroupApplied(ctx, d, svc); err != nil {
			return diag.FromErr(fmt.Errorf("failed to wait until the DevOps firewall group became available: %s", err))
		}
	}

	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old rules and delete any obsolete ones
		obsoleteRuleIds := make([]string, len(ors.List()))
		for i, r := range ors.List() {
			rule := r.(map[string]interface{})
			obsoleteRuleIds[i] = rule["id"].(string)
		}

		if len(obsoleteRuleIds) != 0 {
			revokeInput := expandRevokeFirewallRulesInput(d, obsoleteRuleIds)

			_, err := svc.RevokeFirewallRules(ctx, revokeInput)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to revoke DevOps firewall rules: %s", err))
			}

			if err := waitUntilFirewallRulesRevoked(ctx, d, svc, obsoleteRuleIds); err != nil {
				return diag.FromErr(fmt.Errorf("failed to wait until the DevOps firewall group became available: %s", err))
			}
		}

		// Make sure we save the state of the currently configured rules
		rules := o.(*schema.Set).Intersection(n.(*schema.Set))
		if err := d.Set("rule", rules); err != nil {
			return diag.FromErr(err)
		}

		// Then loop through all the newly configured rules and create them
		newRules := make([]types.RequestRules, len(nrs.List()))
		for i, r := range nrs.List() {
			rule := r.(map[string]interface{})

			var port *int32
			if rule["port"] != 0 {
				v := int32(rule["port"].(int))
				port = &v
			}

			var description *string
			if rule["description"] != "" {
				v := rule["description"].(string)
				description = &v
			}

			newRules[i] = types.RequestRules{
				Protocol:    types.ProtocolOfrulesForAuthorizeFirewallRules(rule["protocol"].(string)),
				Port:        port,
				CidrIp:      types.CidrIpOfrulesForAuthorizeFirewallRules(rule["cidr_ip"].(string)),
				Description: description,
			}

			// Add each rule to state to save it after successful authorization
			rules.Add(rule)
		}

		if len(newRules) != 0 {
			authorizeInput := expandAuthorizeFirewallRulesInput(d, newRules)

			_, err := svc.AuthorizeFirewallRules(ctx, authorizeInput)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to authorize DevOps firewall rules: %s", err))
			}

			if err := waitUntilFirewallGroupApplied(ctx, d, svc); err != nil {
				return diag.FromErr(fmt.Errorf("failed to wait until the DevOps firewall group became available: %s", err))
			}
		}

		if err := d.Set("rule", rules); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("name", "description") {
		input := expandUpdateFirewallGroupInput(d)

		if _, err := svc.UpdateFirewallGroup(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update a FirewallGroup: %s", err))
		}
	}

	if d.HasChange("name") {
		d.SetId(d.Get("name").(string))
	}

	return readFirewallGroup(ctx, d, meta)
}
