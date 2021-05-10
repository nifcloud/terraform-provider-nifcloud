package nassecuritygroup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
)

func waitUntilNASSecurityGroupRuleRevoked(ctx context.Context, d *schema.ResourceData, svc *nas.Client, rule map[string]interface{}) error {
	const maxRetryCount = 20
	const pollInterval = 10 * time.Second

	retryCount := 0
	for {
		if retryCount > maxRetryCount {
			return fmt.Errorf("max retry count exceeded about waiting NAS security group rule revoked: rule -> %v", rule)
		}

		input := expandDescribeNASSecurityGroupsInput(d)
		req := svc.DescribeNASSecurityGroupsRequest(input)
		res, err := req.Send(ctx)
		if err != nil {
			return err
		}

		targetExists := false
		if rule["cidr_ip"] != "" {
			target := rule["cidr_ip"].(string)
			for _, ip := range res.DescribeNASSecurityGroupsOutput.NASSecurityGroups[0].IPRanges {
				if nifcloud.StringValue(ip.CIDRIP) == target && nifcloud.StringValue(ip.Status) == "revoking" {
					targetExists = true
					break
				}
			}
		} else {
			target := rule["security_group_name"].(string)
			for _, group := range res.DescribeNASSecurityGroupsOutput.NASSecurityGroups[0].SecurityGroups {
				if nifcloud.StringValue(group.SecurityGroupName) == target && nifcloud.StringValue(group.Status) == "revoking" {
					targetExists = true
					break
				}
			}
		}

		if !targetExists {
			return nil
		}

		time.Sleep(pollInterval)
		retryCount++
	}
}
