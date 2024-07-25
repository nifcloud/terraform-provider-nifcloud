package devopsrunner

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func createRunner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	input := expandCreateRunnerInput(d)

	res, err := svc.CreateRunner(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create a DevOps Runner: %s", err))
	}

	d.SetId(nifcloud.ToString(res.Runner.RunnerName))

	err = waitUntilRunnerRunning(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait for the DevOps Runner to become ready: %s", err))
	}

	return updateRunner(ctx, d, meta)
}
