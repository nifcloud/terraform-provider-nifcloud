package separateinstancerule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("description") {
		input := expandNiftyUpdateSeparateInstanceRuleInputForDescription(d)

		req := svc.NiftyUpdateSeparateInstanceRuleRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating separate instance rule description: %s", err))
		}
	}

	if d.HasChange("name") && !d.IsNewResource() {
		input := expandNiftyUpdateSeparateInstanceRuleInputForName(d)

		req := svc.NiftyUpdateSeparateInstanceRuleRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating separate instance rule name %s", err))
		}

		d.SetId(d.Get("name").(string))

	}

	return read(ctx, d, meta)
}
