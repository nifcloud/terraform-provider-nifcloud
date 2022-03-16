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

		_, err := svc.NiftyModifyVpnGatewayAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway accounting_type: %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("description") {
		input := expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayDescription(d)

		_, err := svc.NiftyModifyVpnGatewayAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway description: %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("name") {
		input := expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayName(d)

		_, err := svc.NiftyModifyVpnGatewayAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway name %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("type") {
		input := expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayType(d)

		_, err := svc.NiftyModifyVpnGatewayAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway type: %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("ip_address") {
		input := expandNiftyUpdateVpnGatewayNetworkInterfacesInput(d)

		_, err := svc.NiftyUpdateVpnGatewayNetworkInterfaces(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating vpngateway ip_address: %s", err))
		}

		if d := waitForVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("security_group") {
		input := expandNiftyModifyVpnGatewayAttributeInputForSecurityGroup(d)

		_, err := svc.NiftyModifyVpnGatewayAttribute(ctx, input)

		if err != nil {
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

			_, err := svc.NiftyDisassociateRouteTableFromVpnGateway(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed disassociating vpngateway table: %s", err))
			}
		} else if before == "" && after != "" {
			input := expandNiftyAssociateRouteTableWithVpnGatewayInput(d)

			_, err := svc.NiftyAssociateRouteTableWithVpnGateway(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed associating vpngateway table: %s", err))
			}
		} else {
			input := expandNiftyReplaceRouteTableAssociationWithVpnGatewayInput(d)

			_, err := svc.NiftyReplaceRouteTableAssociationWithVpnGateway(ctx, input)
			if err != nil {
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
	deadline, _ := ctx.Deadline()

	if err := computing.NewVpnGatewayAvailableWaiter(svc).Wait(ctx, expandDescribeVpnGatewaysInput(d), time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for vpngateway available: %s", err))
	}

	return nil
}
