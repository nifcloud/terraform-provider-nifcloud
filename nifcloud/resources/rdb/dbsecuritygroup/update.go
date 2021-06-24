package dbsecuritygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	if d.IsNewResource() {
		err := svc.WaitUntilDBSecurityGroupExists(ctx, expandDescribeDBSecurityGroupsInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until db security group available: %s", err))
		}
	}

	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old rules and delete any obsolete ones
		for _, r := range ors.List() {
			rule := r.(map[string]interface{})
			input := expandRevokeDBSecurityGroupIngressInput(d, rule)
			req := svc.RevokeDBSecurityGroupIngressRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating db security group to delete rule: %s", err))
			}

			if err := waitUntilDBSecurityGroupRuleRevoked(ctx, d, svc, rule); err != nil {
				return diag.FromErr(fmt.Errorf("failed wait db security group available: %s", err))
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
			input := expandAuthorizeDBSecurityGroupIngressInput(d, rule)
			req := svc.AuthorizeDBSecurityGroupIngressRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating db security group to create rule: %s", err))
			}

			if rule["cidr_ip"] != "" {
				err = svc.WaitUntilDBSecurityGroupIPRangesAuthorized(ctx, expandDescribeDBSecurityGroupsInput(d))
			} else if rule["security_group_name"] != "" {
				err = svc.WaitUntilDBSecurityGroupEC2SecurityGroupsAuthorized(ctx, expandDescribeDBSecurityGroupsInput(d))
			}

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed wait db security group available: %s", err))
			}

			rules.Add(rule)
			if err := d.Set("rule", rules); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return read(ctx, d, meta)
}
