package securitygroup

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
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteSecurityGroupInput(d)
	svc := meta.(*client.Client).Computing

	if v := d.Get("revoke_rules_on_delete").(bool); v {
		err := forceRevokeSecurityGroupRules(ctx, svc, d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed revoke_rules_on_delete: %s", err))
		}
	}

	req := svc.DeleteSecurityGroupRequest(input)
	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := req.Send(ctx)
		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroup" {
				return nil
			}

			if errors.As(err, &awsErr) && awsErr.Code() == "Client.Inoperable.SecurityGroup.InUse" {
				// If it is a dependency violation, we want to retry
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = svc.WaitUntilSecurityGroupDeleted(ctx, expandDescribeSecurityGroupsInput(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted security group error: %s", err))
	}

	d.SetId("")
	return nil
}

func forceRevokeSecurityGroupRules(ctx context.Context, svc *computing.Client, d *schema.ResourceData) error {
	describeSecurityGroupsinput := expandDescribeSecurityGroupsInput(d)
	req := svc.DescribeSecurityGroupsRequest(describeSecurityGroupsinput)

	res, err := req.Send(ctx)
	if err != nil {
		return err
	}

	if res == nil || len(res.SecurityGroupInfo) == 0 {
		return nil
	}

	securityGroup := res.SecurityGroupInfo[0]

	if len(securityGroup.IpPermissions) > 0 {
		ipPermissions := []computing.RequestIpPermissionsOfRevokeSecurityGroupIngress{}
		for _, i := range securityGroup.IpPermissions {
			ipPermission := computing.RequestIpPermissionsOfRevokeSecurityGroupIngress{
				IpProtocol: computing.IpProtocolOfIpPermissionsForRevokeSecurityGroupIngress(nifcloud.StringValue(i.IpProtocol)),
				FromPort:   i.FromPort,
				ToPort:     i.ToPort,
				InOut:      computing.InOutOfIpPermissionsForRevokeSecurityGroupIngress(nifcloud.StringValue(i.InOut)),
			}
			for _, ipRange := range i.IpRanges {
				ipPermission.ListOfRequestIpRanges = append(
					ipPermission.ListOfRequestIpRanges,
					computing.RequestIpRanges{CidrIp: ipRange.CidrIp},
				)
			}
			for _, group := range i.Groups {
				ipPermission.ListOfRequestGroups = append(
					ipPermission.ListOfRequestGroups,
					computing.RequestGroups{GroupName: group.GroupName},
				)
			}
			ipPermissions = append(ipPermissions, ipPermission)
		}

		err = svc.WaitUntilSecurityGroupApplied(ctx, describeSecurityGroupsinput)
		if err != nil {
			return err
		}

		input := expandRevokeSecurityGroupIngressInput(d, ipPermissions)
		_, err := svc.RevokeSecurityGroupIngressRequest(input).Send(ctx)
		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroupIngress" {
				return nil
			}
			return fmt.Errorf(
				"Error revoking security group %s rules: %s",
				nifcloud.StringValue(securityGroup.GroupName), err)
		}

		err = svc.WaitUntilSecurityGroupApplied(ctx, describeSecurityGroupsinput)
		if err != nil {
			return err
		}
	}
	return nil
}
