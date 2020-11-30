package securitygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.IsNewResource() {
		err := svc.WaitUntilSecurityGroupApplied(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{d.Id()}})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
		}
	}

	if d.HasChange("group_name") {
		input := expandUpdateSecurityGroupInputForName(d)

		req := svc.UpdateSecurityGroupRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating securityGroup name: %s", err))
		}

		groupName := d.Get("group_name").(string)
		d.SetId(groupName)

		err = svc.WaitUntilSecurityGroupApplied(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{d.Id()}})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
		}
	}

	if d.HasChange("description") {
		input := expandUpdateSecurityGroupInputForDescription(d)

		req := svc.UpdateSecurityGroupRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating securityGroup description: %s", err))
		}

		err = svc.WaitUntilSecurityGroupApplied(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{d.Id()}})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
		}
	}

	if d.HasChange("log_limit") {
		input := expandUpdateSecurityGroupInputForLogLimit(d)

		req := svc.UpdateSecurityGroupRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating securityGroup log_limit: %s", err))
		}

		err = svc.WaitUntilSecurityGroupApplied(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{d.Id()}})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
		}
	}

	return read(ctx, d, meta)
}
