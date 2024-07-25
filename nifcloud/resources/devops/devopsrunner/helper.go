package devopsrunner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
)

// waitUntilRunnerRunning waits until the state of the runner become RUNNING.
// DevOps SDK does not provide a waiter.
func waitUntilRunnerRunning(ctx context.Context, d *schema.ResourceData, svc *devopsrunner.Client) error {
	const timeout = 3600 * time.Second

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		input := expandGetRunnerInput(d)
		res, err := svc.GetRunner(ctx, input)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		if nifcloud.ToString(res.Runner.Status.Name) == "RUNNING" {
			return nil
		}

		return retry.RetryableError(fmt.Errorf("expected the DevOps Runner was in state RUNNING"))
	})

	return err
}

// waitUntilRunnerDeleted waits until the state of the runner is deleted.
// DevOps SDK does not provide a waiter.
func waitUntilRunnerDeleted(ctx context.Context, d *schema.ResourceData, svc *devopsrunner.Client) error {
	const timeout = 200 * time.Second

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		input := expandGetRunnerInput(d)
		_, err := svc.GetRunner(ctx, input)
		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Runner" {
				return nil
			}
			return retry.RetryableError(fmt.Errorf("failed to read a DevOps Runner: %s", err))
		}

		return retry.RetryableError(fmt.Errorf("expected the DevOps Runner was deleted"))
	})

	return err
}
