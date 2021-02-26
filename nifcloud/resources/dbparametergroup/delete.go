package dbparametergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func deletegroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	input := expandDeleteDBParameterGroupInput(d)
	req := svc.DeleteDBParameterGroupRequest(input)

	if _, err := req.Send(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting DBParameterGroup: %s", err))
	}

	d.SetId("")

	return nil
}
