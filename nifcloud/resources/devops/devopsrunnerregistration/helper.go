package devopsrunnerregistration

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func validateRunnerRegistrationImportString(importStr string) ([]string, error) {
	// example: runnerName_registrationId

	importParts := strings.Split(importStr, "_")
	errStr := "unexpected format of import string (%q), expected RUNNERNAME_REGISTRATIONID: %s"
	if len(importParts) < 2 {
		return nil, fmt.Errorf(errStr, importStr, "invalid parts")
	}

	runnerName := importParts[0]
	id := importParts[1]

	if runnerName == "" {
		return nil, fmt.Errorf(errStr, importStr, "runner_name must be required")
	}

	if id == "" {
		return nil, fmt.Errorf(errStr, importStr, "id must be required")
	}

	return importParts, nil
}

func populateRunnerRegistrationFromImport(d *schema.ResourceData, importParts []string) error {
	runnerName := importParts[0]
	id := importParts[1]

	if err := d.Set("runner_name", runnerName); err != nil {
		return err
	}

	d.SetId(id)

	return nil
}
