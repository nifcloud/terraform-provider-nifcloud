package dbsecuritygroup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
)

func waitUntilDBSecurityGroupRuleRevoked(ctx context.Context, d *schema.ResourceData, svc *rdb.Client, rule map[string]interface{}) error {
	const maxRetryCount = 20
	const pollInterval = 10 * time.Second

	retryCount := 0
	for {
		if retryCount > maxRetryCount {
			return fmt.Errorf("max retry count exceeded about waiting db security group rule revoked: rule -> %v", rule)
		}

		input := expandDescribeDBSecurityGroupsInput(d)
		req := svc.DescribeDBSecurityGroupsRequest(input)
		res, err := req.Send(ctx)
		if err != nil {
			return err
		}

		targetExists := false
		if rule["cidr_ip"] != "" {
			target := rule["cidr_ip"].(string)
			for _, ip := range res.DescribeDBSecurityGroupsOutput.DBSecurityGroups[0].IPRanges {
				if nifcloud.StringValue(ip.CIDRIP) == target && nifcloud.StringValue(ip.Status) == "revoking" {
					targetExists = true
					break
				}
			}
		} else {
			target := rule["security_group_name"].(string)
			for _, group := range res.DescribeDBSecurityGroupsOutput.DBSecurityGroups[0].EC2SecurityGroups {
				if nifcloud.StringValue(group.EC2SecurityGroupName) == target && nifcloud.StringValue(group.Status) == "revoking" {
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
