package nassecuritygroup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).NAS
	deadline, _ := ctx.Deadline()

	if d.IsNewResource() {
		err := nas.NewNASSecurityGroupExistsWaiter(svc).Wait(ctx, expandDescribeNASSecurityGroupsInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until NAS security group available: %s", err))
		}
	}

	if d.HasChanges("group_name", "description") {
		input := expandModifyNASSecurityGroupInput(d)
		_, err := svc.ModifyNASSecurityGroup(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating NAS security group: %s", err))
		}

		d.SetId(d.Get("group_name").(string))
	}

	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old rules and delete any obsolete ones
		for _, r := range ors.List() {
			rule := r.(map[string]interface{})
			input := expandRevokeNASSecurityGroupIngressInput(d, rule)
			_, err := svc.RevokeNASSecurityGroupIngress(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating NAS security group to delete rule: %s", err))
			}

			if err := waitUntilNASSecurityGroupRuleRevoked(ctx, d, svc, rule); err != nil {
				return diag.FromErr(fmt.Errorf("failed wait NAS security group available: %s", err))
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
			input := expandAuthorizeNASSecurityGroupIngressInput(d, rule)
			_, err := svc.AuthorizeNASSecurityGroupIngress(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating NAS security group to create rule: %s", err))
			}

			if rule["cidr_ip"] != "" {
				err = nas.NewNASSecurityGroupIPRangesAuthorizedWaiter(svc).Wait(ctx, expandDescribeNASSecurityGroupsInput(d), time.Until(deadline))
			} else if rule["security_group_name"] != "" {
				err = nas.NewNASSecurityGroupSecurityGroupsAuthorizedWaiter(svc).Wait(ctx, expandDescribeNASSecurityGroupsInput(d), time.Until(deadline))
			}
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed wait NAS security group available: %s", err))
			}

			rules.Add(rule)
			if err := d.Set("rule", rules); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return read(ctx, d, meta)
}
