package devopsrunnerparametergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func deleteRunnerParameterGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	input := expandDeleteRunnerParameterGroupInput(d)

	if _, err := svc.DeleteRunnerParameterGroup(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete a DevOps Runner parameter group: %s", err))
	}

	d.SetId("")

	return nil
}
