package customergateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("name") {
		input := expandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayName(d)

		_, err := svc.NiftyModifyCustomerGatewayAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating name: %s", err))
		}
	}

	if d.HasChange("description") {
		input := expandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayDescription(d)

		_, err := svc.NiftyModifyCustomerGatewayAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating description: %s", err))
		}
	}

	return read(ctx, d, meta)
}
