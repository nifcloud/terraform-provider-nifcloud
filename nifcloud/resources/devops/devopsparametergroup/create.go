package devopsparametergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func createParameterGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	input := expandCreateParameterGroupInput(d)

	res, err := svc.CreateParameterGroup(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating ParameterGroup: %s", err))
	}

	d.SetId(nifcloud.ToString(res.ParameterGroup.ParameterGroupName))

	return updateParameterGroup(ctx, d, meta)
}
