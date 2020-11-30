package privatelan

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("private_lan_name") {
		input := expandNiftyModifyPrivateLanAttributeInputForPrivateLanName(d)

		svc := meta.(*client.Client).Computing
		req := svc.NiftyModifyPrivateLanAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating private_lan private_lan_name: %s", err))
		}
	}

	if d.HasChange("cidr_block") {
		input := expandNiftyModifyPrivateLanAttributeInputForCidrBlock(d)

		svc := meta.(*client.Client).Computing
		req := svc.NiftyModifyPrivateLanAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating private_lan cidr_block: %s", err))
		}
	}

	if d.HasChange("accounting_type") {
		input := expandNiftyModifyPrivateLanAttributeInputForAccountingType(d)

		svc := meta.(*client.Client).Computing
		req := svc.NiftyModifyPrivateLanAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating private_lan accounting_type: %s", err))
		}
	}

	if d.HasChange("description") {
		input := expandNiftyModifyPrivateLanAttributeInputForDescription(d)

		svc := meta.(*client.Client).Computing
		req := svc.NiftyModifyPrivateLanAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating private_lan description: %s", err))
		}
	}
	return read(ctx, d, meta)
}
