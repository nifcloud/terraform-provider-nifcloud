package devopsrunner

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func updateRunner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	if d.HasChange("instance_type") && !d.IsNewResource() {
		input := expandModifyRunnerInstanceTypeInput(d)

		if _, err := svc.ModifyRunnerInstanceType(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed to modify instance type of the DevOps Runner: %s", err))
		}

		err := waitUntilRunnerRunning(ctx, d, svc)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to wait for the DevOps Runner to become ready: %s", err))
		}
	}

	if d.HasChanges("name", "concurrent", "description") && !d.IsNewResource() {
		input := expandUpdateRunnerInput(d)

		if _, err := svc.UpdateRunner(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update a DevOps Runner: %s", err))
		}
	}

	if d.HasChange("name") {
		d.SetId(d.Get("name").(string))
	}

	return readRunner(ctx, d, meta)
}
