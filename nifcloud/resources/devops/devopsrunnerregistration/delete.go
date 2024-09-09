package devopsrunnerregistration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func deleteRunnerRegistration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	input := expandUnregisterRunnerInput(d)

	if _, err := svc.UnregisterRunner(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed to unregister the DevOps runner from a GitLab instance: %s", err))
	}

	err := waitUntilRunnerRunning(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait for the DevOps Runner to become ready: %s", err))
	}

	d.SetId("")

	return nil
}
