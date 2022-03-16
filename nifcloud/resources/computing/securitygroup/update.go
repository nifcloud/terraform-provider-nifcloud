package securitygroup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if d.IsNewResource() {
		err := computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{d.Id()}}, time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
		}
	}

	if d.HasChange("group_name") {
		input := expandUpdateSecurityGroupInputForName(d)

		_, err := svc.UpdateSecurityGroup(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating securityGroup name: %s", err))
		}

		groupName := d.Get("group_name").(string)
		d.SetId(groupName)

		err = computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{d.Id()}}, time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
		}
	}

	if d.HasChange("description") {
		input := expandUpdateSecurityGroupInputForDescription(d)

		_, err := svc.UpdateSecurityGroup(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating securityGroup description: %s", err))
		}

		err = computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{d.Id()}}, time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
		}
	}

	if d.HasChange("log_limit") {
		input := expandUpdateSecurityGroupInputForLogLimit(d)

		_, err := svc.UpdateSecurityGroup(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating securityGroup log_limit: %s", err))
		}

		err = computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{d.Id()}}, time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
		}
	}

	return read(ctx, d, meta)
}
