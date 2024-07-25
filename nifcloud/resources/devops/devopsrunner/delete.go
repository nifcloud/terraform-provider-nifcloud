package devopsrunner

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func deleteRunner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	input := expandDeleteRunnerInput(d)

	if _, err := svc.DeleteRunner(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete a DevOps Runner: %s", err))
	}

	err := waitUntilRunnerDeleted(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait until the DevOps Runner is deleted: %s", err))
	}

	d.SetId("")

	return nil
}
