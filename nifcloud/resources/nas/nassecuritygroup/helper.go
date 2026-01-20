package nassecuritygroup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
)

func waitUntilNASSecurityGroupRuleRevoked(ctx context.Context, d *schema.ResourceData, svc *nas.Client, rule map[string]interface{}) error {
	const timeout = 200 * time.Second

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		input := expandDescribeNASSecurityGroupsInput(d)
		res, err := svc.DescribeNASSecurityGroups(ctx, input)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		targetExists := false
		if rule["cidr_ip"] != "" {
			target := rule["cidr_ip"].(string)
			for _, ip := range res.NASSecurityGroups[0].IPRanges {
				if nifcloud.ToString(ip.CIDRIP) == target && nifcloud.ToString(ip.Status) == "revoking" {
					targetExists = true
					break
				}
			}
		} else {
			target := rule["security_group_name"].(string)
			for _, group := range res.NASSecurityGroups[0].SecurityGroups {
				if nifcloud.ToString(group.SecurityGroupName) == target && nifcloud.ToString(group.Status) == "revoking" {
					targetExists = true
					break
				}
			}
		}

		if !targetExists {
			return nil
		}

		return retry.RetryableError(fmt.Errorf("expected rule to revoked but was in state revoking"))
	})

	return err
}
