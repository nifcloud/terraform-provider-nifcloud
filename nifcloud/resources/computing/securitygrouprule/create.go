package securitygrouprule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inputList := expandAuthorizeSecurityGroupIngressInputList(d)

	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	describeSecurityGroupsInput := expandDescribeSecurityGroupsInput(d)
	describeSecurityGroupsOutput, err := svc.DescribeSecurityGroups(ctx, describeSecurityGroupsInput)
	if err != nil {
		return diag.Errorf("failed describe security groups: %s", err)
	}

	eg, ctxt := errgroup.WithContext(ctx)
	for _, input := range inputList {
		input := input
		eg.Go(func() error {
			mutexKV.Lock(nifcloud.ToString(input.GroupName))
			defer mutexKV.Unlock(nifcloud.ToString(input.GroupName))

			err := checkSecurityGroupExist(describeSecurityGroupsOutput.SecurityGroupInfo, nifcloud.ToString(input.GroupName))
			if err != nil {
				return err
			}

			err = computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctxt, &computing.DescribeSecurityGroupsInput{GroupName: []string{nifcloud.ToString(input.GroupName)}}, time.Until(deadline))
			if err != nil {
				return fmt.Errorf("failed wait until securityGroup applied: %s", err)
			}

			_, err = svc.AuthorizeSecurityGroupIngress(ctx, input)
			if err != nil {
				return fmt.Errorf("failed creating securityGroup rule: %s", err)
			}

			err = computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctxt, &computing.DescribeSecurityGroupsInput{GroupName: []string{nifcloud.ToString(input.GroupName)}}, time.Until(deadline))
			if err != nil {
				return fmt.Errorf("failed wait until securityGroup applied: %s", err)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return diag.FromErr(err)
	}

	id := idHash(inputList)
	d.SetId(id)

	return read(ctx, d, meta)
}
