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

	if d.HasChange("customer_gateway_name") {
		input := expandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayName(d)

		req := svc.NiftyModifyCustomerGatewayAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating customer gateway customer_gateway_name: %s", err))
		}
	}

	if d.HasChange("customer_gateway_description") {
		input := expandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayDescription(d)

		req := svc.NiftyModifyCustomerGatewayAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating customer gateway customer_gateway_description: %s", err))
		}
	}

	return read(ctx, d, meta)
}
