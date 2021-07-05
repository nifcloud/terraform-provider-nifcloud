package networkinterface

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
		input := expandModifyNetworkInterfaceAttributeInputForDescription(d)

		req := svc.ModifyNetworkInterfaceAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating network interface description: %s", err))
		}
	}

	if d.HasChange("ip_address") {
		if err := waitForRouterOfNetworkInterfaceAvailable(ctx, d, svc); err != nil {
			return err
		}

		input := expandModifyNetworkInterfaceAttributeInputForIPAddress(d)
		req := svc.ModifyNetworkInterfaceAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating network interface ip address: %s", err))
		}

		if err := waitForRouterOfNetworkInterfaceAvailable(ctx, d, svc); err != nil {
			return err
		}
	}
	return read(ctx, d, meta)
}
