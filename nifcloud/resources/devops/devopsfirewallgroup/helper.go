package devopsfirewallgroup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
)

// waitUntilFirewallGroupApplied waits until the state of the firewall group become APPLIED.
// DevOps SDK does not provide a waiter.
func waitUntilFirewallGroupApplied(ctx context.Context, d *schema.ResourceData, svc *devops.Client) error {
	const timeout = 200 * time.Second

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		input := expandGetFirewallGroupInput(d)
		res, err := svc.GetFirewallGroup(ctx, input)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		if nifcloud.ToString(res.FirewallGroup.State) == "APPLIED" {
			return nil
		}

		return retry.RetryableError(fmt.Errorf("expected the DevOps firewall group was in state APPLIED"))
	})

	return err
}

// waitUntilFirewallGroupDeleted waits until the state of the firewall group is deleted.
// DevOps SDK does not provide a waiter.
func waitUntilFirewallGroupDeleted(ctx context.Context, d *schema.ResourceData, svc *devops.Client) error {
	const timeout = 200 * time.Second

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		input := expandGetFirewallGroupInput(d)
		_, err := svc.GetFirewallGroup(ctx, input)
		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.FirewallGroup" {
				return nil
			}
			return retry.RetryableError(fmt.Errorf("failed to read a DevOps firewall group: %s", err))
		}

		return retry.RetryableError(fmt.Errorf("expected the DevOps firewall group was deleted"))
	})

	return err
}

// waitUntilFirewallRulesRevoked waits until all of the rules specified are revoked.
// DevOps SDK does not provide a waiter.
func waitUntilFirewallRulesRevoked(ctx context.Context, d *schema.ResourceData, svc *devops.Client, ruleIds []string) error {
	const timeout = 200 * time.Second

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		input := expandGetFirewallGroupInput(d)
		res, err := svc.GetFirewallGroup(ctx, input)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		targetsExist := false
		for _, ruleId := range ruleIds {
			for _, rule := range res.FirewallGroup.Rules {
				if nifcloud.ToString(rule.Id) == ruleId && nifcloud.ToString(rule.State) == "REVOKING" {
					targetsExist = true
				}
			}
		}

		if !targetsExist {
			return nil
		}

		return retry.RetryableError(fmt.Errorf("expected the DevOps firewall rules have been revoked but were in state REVOKING"))
	})

	return err
}
