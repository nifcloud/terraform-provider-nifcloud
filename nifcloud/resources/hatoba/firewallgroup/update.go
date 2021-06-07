package firewallgroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Hatoba

	if d.HasChanges("name", "description") {
		input := expandUpdateFirewallGroupInput(d)
		req := svc.UpdateFirewallGroupRequest(input)
		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating Hatoba firewall group: %s", err))
		}

		d.SetId(d.Get("name").(string))
	}

	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old rules and delete any obsolete ones
		for _, r := range ors.List() {
			rule := r.(map[string]interface{})
			input := expandRevokeFirewallGroupInput(d, rule)
			req := svc.RevokeFirewallGroupRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed updating Hatoba firewall group to delete rule: %s", err))
			}

			if err := svc.WaitUntilFirewallRuleAuthorized(ctx, expandGetFirewallGroupInput(d)); err != nil {
				return diag.FromErr(fmt.Errorf("failed wait Hatoba firewall group available: %s", err))
			}
		}

		// Make sure we save the state of the currently configured rules
		rules := o.(*schema.Set).Intersection(n.(*schema.Set))
		if err := d.Set("rule", rules); err != nil {
			return diag.FromErr(err)
		}

		// Then loop through all the newly configured rules and create them
		for _, r := range nrs.List() {
			rule := r.(map[string]interface{})
			input := expandAuthorizeFirewallGroupInput(d, rule)
			req := svc.AuthorizeFirewallGroupRequest(input)
			res, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating Hatoba firewall group to create rule: %s", err))
			}

			if err := svc.WaitUntilFirewallRuleAuthorized(ctx, expandGetFirewallGroupInput(d)); err != nil {
				return diag.FromErr(fmt.Errorf("failed wait Hatoba firewall group available: %s", err))
			}

			for _, resRule := range res.FirewallGroup.Rules {
				if nifcloud.StringValue(resRule.Status) == "AUTHORIZING" {
					rule["id"] = resRule.Id
				}
			}

			rules.Add(rule)
			if err := d.Set("rule", rules); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return read(ctx, d, meta)
}
