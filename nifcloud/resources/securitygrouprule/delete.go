package securitygrouprule

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inputList := expandRevokeSecurityGroupIngressInputList(d)
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
			err = checkSecurityGroupExist(describeSecurityGroupsOutput.SecurityGroupInfo, nifcloud.StringValue(input.GroupName))
			if err != nil {
				return nil
			}

			err := svc.WaitUntilSecurityGroupApplied(ctxt, &computing.DescribeSecurityGroupsInput{GroupName: []string{nifcloud.StringValue(input.GroupName)}})
			if err != nil {
				return fmt.Errorf("failed wait until securityGroup applied: %s", err)
			}

			req := svc.RevokeSecurityGroupIngressRequest(input)

			err = resource.RetryContext(ctxt, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
				_, err = req.Send(ctxt)
				if err != nil {
					var awsErr awserr.Error
					if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroupIngress" {
						return nil
					}
					return resource.RetryableError(err)
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
