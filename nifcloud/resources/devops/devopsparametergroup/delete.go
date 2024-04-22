package devopsparametergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func deleteParameterGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	input := expandDeleteParameterGroupInput(d)

	if _, err := svc.DeleteParameterGroup(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting ParameterGroup: %s", err))
	}

	d.SetId("")

	return nil
}
