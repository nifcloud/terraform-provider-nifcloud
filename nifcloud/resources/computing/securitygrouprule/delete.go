package securitygrouprule

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inputList := expandRevokeSecurityGroupIngressInputList(d)
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

			err = checkSecurityGroupExist(describeSecurityGroupsOutput.SecurityGroupInfo, nifcloud.ToString(input.GroupName))
			if err != nil {
				return nil
			}

			err := computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctxt, &computing.DescribeSecurityGroupsInput{GroupName: []string{nifcloud.ToString(input.GroupName)}}, time.Until(deadline))
			if err != nil {
				return fmt.Errorf("failed wait until securityGroup applied: %s", err)
			}

			_, err = svc.RevokeSecurityGroupIngress(ctx, input)

			err = retry.RetryContext(ctxt, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
				if err != nil {
					var awsErr smithy.APIError
					if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.SecurityGroupIngress" {
						return nil
					}
					return retry.RetryableError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed deleting securityGroup rule: %s", err)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
