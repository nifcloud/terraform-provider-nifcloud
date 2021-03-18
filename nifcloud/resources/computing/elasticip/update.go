package elasticip

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("description") {

		input := expandNiftyModifyAddressAttributeInput(d)

		svc := meta.(*client.Client).Computing
		req := svc.NiftyModifyAddressAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating elastic ip: %s", err))
		}
	}
	return read(ctx, d, meta)
}
