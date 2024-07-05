package devopsinstance

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

// waitUntilInstanceRunning waits until the state of the instance become RUNNING.
// DevOps SDK does not provide a waiter.
func waitUntilInstanceRunning(ctx context.Context, d *schema.ResourceData, svc *devops.Client) error {
	const timeout = 3600 * time.Second

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		input := expandGetInstanceInput(d)
		res, err := svc.GetInstance(ctx, input)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		if nifcloud.ToString(res.Instance.State) == "RUNNING" {
			return nil
		}

		return retry.RetryableError(fmt.Errorf("expected the DevOps instance was in state RUNNING"))
	})

	return err
}

// waitUntilInstanceDeleted waits until the state of the instance is deleted.
// DevOps SDK does not provide a waiter.
func waitUntilInstanceDeleted(ctx context.Context, d *schema.ResourceData, svc *devops.Client) error {
	const timeout = 200 * time.Second

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		input := expandGetInstanceInput(d)
		_, err := svc.GetInstance(ctx, input)
		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Instance" {
				return nil
			}
			return retry.RetryableError(fmt.Errorf("failed to read a DevOps instance: %s", err))
		}

		return retry.RetryableError(fmt.Errorf("expected the DevOps instance was deleted"))
	})

	return err
}
