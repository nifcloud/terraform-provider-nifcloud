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

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateSecurityGroupInput(d)
	deadline, _ := ctx.Deadline()

	svc := meta.(*client.Client).Computing
	_, err := svc.CreateSecurityGroup(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating SecurityGroup: %s", err))
	}

	groupName := d.Get("group_name").(string)
	d.SetId(groupName)

	err = computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{groupName}}, time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
	}

	_, err = svc.UpdateSecurityGroup(ctx, expandUpdateSecurityGroupInputForLogLimit(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed updating SecurityGroup: %s", err))
	}
	return read(ctx, d, meta)
}
