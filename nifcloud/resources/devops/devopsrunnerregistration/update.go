package devopsrunnerregistration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func updateRunnerRegistration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	if d.HasChange("parameter_group_name") {
		input := expandUpdateRunnerRegistrationInput(d)

		if _, err := svc.UpdateRunnerRegistration(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update a DevOps Runner registration: %s", err))
		}

		err := waitUntilRunnerRunning(ctx, d, svc)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to wait for the DevOps Runner to become ready: %s", err))
		}
	}

	return readRunnerRegistration(ctx, d, meta)
}
