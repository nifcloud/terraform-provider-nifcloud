package securitygrouprule

import (
	"context"
	"fmt"

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

	describeSecurityGroupsInput := expandDescribeSecurityGroupsInput(d)
	describeSecurityGroupsOutput, err := svc.DescribeSecurityGroupsRequest(describeSecurityGroupsInput).Send(ctx)
	if err != nil {
		return diag.Errorf("failed describe security groups: %s", err)
	}

	eg, ctxt := errgroup.WithContext(ctx)
	for _, input := range inputList {
		input := input
		eg.Go(func() error {
			mutexKV.Lock(nifcloud.StringValue(input.GroupName))
			defer mutexKV.Unlock(nifcloud.StringValue(input.GroupName))

			err := checkSecurityGroupExist(describeSecurityGroupsOutput.SecurityGroupInfo, nifcloud.StringValue(input.GroupName))
			if err != nil {
				return err
			}

			err = svc.WaitUntilSecurityGroupApplied(ctxt, &computing.DescribeSecurityGroupsInput{GroupName: []string{nifcloud.StringValue(input.GroupName)}})
			if err != nil {
				return fmt.Errorf("failed wait until securityGroup applied: %s", err)
			}

			req := svc.AuthorizeSecurityGroupIngressRequest(input)

			_, err = req.Send(ctxt)
			if err != nil {
				return fmt.Errorf("failed creating securityGroup rule: %s", err)
			}

			err = svc.WaitUntilSecurityGroupApplied(ctxt, &computing.DescribeSecurityGroupsInput{GroupName: []string{nifcloud.StringValue(input.GroupName)}})
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
