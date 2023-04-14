package securitygroup

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
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteSecurityGroupInput(d)
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if v := d.Get("revoke_rules_on_delete").(bool); v {
		err := forceRevokeSecurityGroupRules(ctx, svc, d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed revoke_rules_on_delete: %s", err))
		}
	}

	_, err := svc.DeleteSecurityGroup(ctx, input)
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.SecurityGroup" {
				return nil
			}

			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.Inoperable.SecurityGroup.InUse" {
				// If it is a dependency violation, we want to retry
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = computing.NewSecurityGroupDeletedWaiter(svc).Wait(ctx, expandDescribeSecurityGroupsInput(d), time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted security group error: %s", err))
	}

	d.SetId("")
	return nil
}

func forceRevokeSecurityGroupRules(ctx context.Context, svc *computing.Client, d *schema.ResourceData) error {
	describeSecurityGroupsinput := expandDescribeSecurityGroupsInput(d)
	deadline, _ := ctx.Deadline()
	res, err := svc.DescribeSecurityGroups(ctx, describeSecurityGroupsinput)
	if err != nil {
		return err
	}

	if res == nil || len(res.SecurityGroupInfo) == 0 {
		return nil
	}

	securityGroup := res.SecurityGroupInfo[0]

	if len(securityGroup.IpPermissions) > 0 {
		ipPermissions := []types.RequestIpPermissionsOfRevokeSecurityGroupIngress{}
		for _, i := range securityGroup.IpPermissions {
			ipPermission := types.RequestIpPermissionsOfRevokeSecurityGroupIngress{
				IpProtocol: types.IpProtocolOfIpPermissionsForRevokeSecurityGroupIngress(nifcloud.ToString(i.IpProtocol)),
				FromPort:   i.FromPort,
				ToPort:     i.ToPort,
				InOut:      types.InOutOfIpPermissionsForRevokeSecurityGroupIngress(nifcloud.ToString(i.InOut)),
			}
			for _, ipRange := range i.IpRanges {
				ipPermission.ListOfRequestIpRanges = append(
					ipPermission.ListOfRequestIpRanges,
					types.RequestIpRanges{CidrIp: ipRange.CidrIp},
				)
			}
			for _, group := range i.Groups {
				ipPermission.ListOfRequestGroups = append(
					ipPermission.ListOfRequestGroups,
					types.RequestGroups{GroupName: group.GroupName},
				)
			}
			ipPermissions = append(ipPermissions, ipPermission)
		}

		err = computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, describeSecurityGroupsinput, time.Until(deadline))
		if err != nil {
			return err
		}

		input := expandRevokeSecurityGroupIngressInput(d, ipPermissions)
		_, err := svc.RevokeSecurityGroupIngress(ctx, input)
		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.SecurityGroupIngress" {
				return nil
			}
			return fmt.Errorf(
				"Error revoking security group %s rules: %s",
				nifcloud.ToString(securityGroup.GroupName), err)
		}

		err = computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, describeSecurityGroupsinput, time.Until(deadline))
		if err != nil {
			return err
		}
	}
	return nil
}
