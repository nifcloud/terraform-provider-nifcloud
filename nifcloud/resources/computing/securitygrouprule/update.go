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

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	inputList := expandAuthorizeSecurityGroupIngressInputList(d)

	if d.HasChange("security_group_names") {
		before, after := d.GetChange("security_group_names")
		var addList, removeList []string

		for _, b := range before.([]interface{}) {
			found := false
			for _, a := range after.([]interface{}) {
				if a.(string) == b.(string) {
					found = true
					break
				}
			}
			if !found {
				removeList = append(removeList, b.(string))
			}
		}

		for _, a := range after.([]interface{}) {
			found := false
			for _, b := range before.([]interface{}) {
				if a.(string) == b.(string) {
					found = true
					break
				}
			}
			if !found {
				addList = append(addList, a.(string))
			}
		}

		// Add rule in security group if the target security group has been added
		if len(addList) > 0 {
			if err := d.Set("security_group_names", addList); err != nil {
				return diag.FromErr(err)
			}

			describeSecurityGroupsInput := expandDescribeSecurityGroupsInput(d)
			describeSecurityGroupsOutput, err := svc.DescribeSecurityGroupsRequest(describeSecurityGroupsInput).Send(ctx)
			if err != nil {
				return diag.Errorf("failed describe security groups: %s", err)
			}

			authorizeInputList := expandAuthorizeSecurityGroupIngressInputList(d)

			eg, ctxt := errgroup.WithContext(ctx)
			for _, input := range authorizeInputList {
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
		}

		// Remove rule in security group if the target security group has been deleted
		if len(removeList) > 0 {
			if err := d.Set("security_group_names", removeList); err != nil {
				return diag.FromErr(err)
			}

			describeSecurityGroupsInput := expandDescribeSecurityGroupsInput(d)
			describeSecurityGroupsOutput, err := svc.DescribeSecurityGroupsRequest(describeSecurityGroupsInput).Send(ctx)
			if err != nil {
				return diag.Errorf("failed describe security groups: %s", err)
			}

			revokeInputList := expandRevokeSecurityGroupIngressInputList(d)

			eg, ctxt := errgroup.WithContext(ctx)
			for _, input := range revokeInputList {
				input := input
				eg.Go(func() error {
					mutexKV.Lock(nifcloud.StringValue(input.GroupName))
					defer mutexKV.Unlock(nifcloud.StringValue(input.GroupName))

					err = checkSecurityGroupExist(describeSecurityGroupsOutput.SecurityGroupInfo, nifcloud.StringValue(input.GroupName))
					if err != nil {
						return nil
					}

					err = svc.WaitUntilSecurityGroupApplied(ctxt, &computing.DescribeSecurityGroupsInput{GroupName: []string{nifcloud.StringValue(input.GroupName)}})
					if err != nil {
						return fmt.Errorf("failed wait until securityGroup applied: %s", err)
					}

					req := svc.RevokeSecurityGroupIngressRequest(input)

					err := resource.RetryContext(ctxt, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
						_, err := req.Send(ctxt)
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

		}

		if err := d.Set("security_group_names", after); err != nil {
			return diag.FromErr(err)
		}

		id := idHash(inputList)
		d.SetId(id)
	}
	return read(ctx, d, meta)
}
