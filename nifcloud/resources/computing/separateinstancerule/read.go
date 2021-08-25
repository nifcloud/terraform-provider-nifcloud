package separateinstancerule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyDescribeSeparateInstanceRulesInput(d)
	svc := meta.(*client.Client).Computing

	req := svc.NiftyDescribeSeparateInstanceRulesRequest(input)

	res, err := req.Send(ctx)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed reading: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
