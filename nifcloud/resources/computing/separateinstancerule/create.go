package separateinstancerule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyCreateSeparateInstanceRuleInput(d)

	svc := meta.(*client.Client).Computing
	req := svc.NiftyCreateSeparateInstanceRuleRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating SeparateInstanceRule: %s", err))
	}

	SeparateInstanceRuleName := d.Get("name").(string)
	d.SetId(SeparateInstanceRuleName)

	return read(ctx, d, meta)
}
