package securitygroup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDescribeSecurityGroupsInput(d)
	deadline, _ := ctx.Deadline()

	svc := meta.(*client.Client).Computing

	if d.IsNewResource() {
		err := computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{d.Id()}}, time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
		}
	}

	res, err := svc.DescribeSecurityGroups(ctx, input)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.SecurityGroup" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
