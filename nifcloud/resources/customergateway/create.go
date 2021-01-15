package customergateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateCustomerGatewayInput(d)
	svc := meta.(*client.Client).Computing
	req := svc.CreateCustomerGatewayRequest(input)

	res, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating customerGateway: %s", err))
	}

	customerGatewayID := res.CustomerGateway.CustomerGatewayId
	d.SetId(nifcloud.StringValue(customerGatewayID))

	return update(ctx, d, meta)
}
