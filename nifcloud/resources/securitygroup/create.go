package securitygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateSecurityGroupInput(d)

	svc := meta.(*client.Client).Computing
	req := svc.CreateSecurityGroupRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating SecurityGroup: %s", err))
	}

	groupName := d.Get("group_name").(string)
	d.SetId(groupName)

	err = svc.WaitUntilSecurityGroupApplied(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{groupName}})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
	}

	_, err = svc.UpdateSecurityGroupRequest(expandUpdateSecurityGroupInputForLogLimit(d)).Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed updating SecurityGroup: %s", err))
	}
	return read(ctx, d, meta)
}
