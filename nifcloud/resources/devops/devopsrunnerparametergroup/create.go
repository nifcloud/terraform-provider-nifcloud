package devopsrunnerparametergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func createRunnerParameterGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	input := expandCreateRunnerParameterGroupInput(d)

	res, err := svc.CreateRunnerParameterGroup(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create a DevOps Runner parameter group: %s", err))
	}

	d.SetId(nifcloud.ToString(res.ParameterGroup.ParameterGroupName))

	return updateRunnerParameterGroup(ctx, d, meta)
}
