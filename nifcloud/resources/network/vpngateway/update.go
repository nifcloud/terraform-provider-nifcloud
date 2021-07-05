package vpngateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

const waiterInitialDelayForUpdate = 3

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("accounting_type") {
		input := expandNiftyModifyVpnGatewayAttributeInputForAccountingType(d)

		req := svc.NiftyModifyVpnGatewayAttributeRequest(input)
		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway accounting_type: %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("description") {
		input := expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayDescription(d)

		req := svc.NiftyModifyVpnGatewayAttributeRequest(input)
		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway description: %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("name") {
		input := expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayName(d)

		req := svc.NiftyModifyVpnGatewayAttributeRequest(input)

		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway name %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("type") {
		input := expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayType(d)

		req := svc.NiftyModifyVpnGatewayAttributeRequest(input)

		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway type: %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("ip_address") {
		input := expandNiftyUpdateVpnGatewayNetworkInterfacesInput(d)

		req := svc.NiftyUpdateVpnGatewayNetworkInterfacesRequest(input)

		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway ip_address: %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("security_group") {
		input := expandNiftyModifyVpnGatewayAttributeInputForSecurityGroup(d)

		req := svc.NiftyModifyVpnGatewayAttributeRequest(input)

		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway security_group: %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("route_table_id") {
		before, after := d.GetChange("route_table_id")
		if before != "" && after == "" {
			input := expandNiftyDisassociateRouteTableFromVpnGatewayInput(d)

			req := svc.NiftyDisassociateRouteTableFromVpnGatewayRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed disassociating vpngateway table: %s", err))
			}
		} else if before == "" && after != "" {
			input := expandNiftyAssociateRouteTableWithVpnGatewayInput(d)

			req := svc.NiftyAssociateRouteTableWithVpnGatewayRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed associating vpngateway table: %s", err))
			}
		} else {
			input := expandNiftyReplaceRouteTableAssociationWithVpnGatewayInput(d)

			req := svc.NiftyReplaceRouteTableAssociationWithVpnGatewayRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed updating vpngateway route_table_id: %s", err))
			}
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	return read(ctx, d, meta)
}
func waitForVpnGatewayAvailable(ctx context.Context, d *schema.ResourceData, svc *computing.Client) diag.Diagnostics {
	// lintignore:R018
	time.Sleep(waiterInitialDelayForUpdate * time.Second)

	if err := svc.WaitUntilVpnGatewayAvailable(ctx, expandDescribeVpnGatewaysInput(d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for vpngateway available: %s", err))
	}

	return nil
}
