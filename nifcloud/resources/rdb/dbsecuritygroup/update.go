package dbsecuritygroup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	deadline, _ := ctx.Deadline()
	if d.IsNewResource() {
		err := rdb.NewDBSecurityGroupExistsWaiter(svc).Wait(ctx, expandDescribeDBSecurityGroupsInput(d), time.Until(deadline))
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

			_, err := svc.RevokeDBSecurityGroupIngress(ctx, input)
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

			_, err := svc.AuthorizeDBSecurityGroupIngress(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating db security group to create rule: %s", err))
			}

			if rule["cidr_ip"] != "" {
				err = rdb.NewDBSecurityGroupIPRangesAuthorizedWaiter(svc).Wait(ctx, expandDescribeDBSecurityGroupsInput(d), time.Until(deadline))
			} else if rule["security_group_name"] != "" {
				err = rdb.NewDBSecurityGroupEC2SecurityGroupsAuthorizedWaiter(svc).Wait(ctx, expandDescribeDBSecurityGroupsInput(d), time.Until(deadline))
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
