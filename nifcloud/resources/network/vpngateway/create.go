package vpngateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

const waiterInitialDelayForCreate = 3

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateVpnGatewayInput(d)
	svc := meta.(*client.Client).Computing
	req := svc.CreateVpnGatewayRequest(input)

	res, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating VPN Gateway: %s", err))
	}

	d.SetId(nifcloud.StringValue(res.VpnGateway.VpnGatewayId))

	describeVpnGatewaysInput := expandDescribeVpnGatewaysInput(d)

	if _, ok := d.GetOk("route_table_id"); ok {
		if err := svc.WaitUntilVpnGatewayAvailable(ctx, describeVpnGatewaysInput); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for vpngateway to become ready: %s", err))
		}
		niftyAssociateRouteTableWithVpnGatewayInput := expandNiftyAssociateRouteTableWithVpnGatewayInput(d)
		niftyAssociateRouteTableWithVpnGatewayRequest := svc.NiftyAssociateRouteTableWithVpnGatewayRequest(niftyAssociateRouteTableWithVpnGatewayInput)
		if _, err := niftyAssociateRouteTableWithVpnGatewayRequest.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed associating route table to vpngateway: %s", err))
		}
	}

	// lintignore:R018
	time.Sleep(waiterInitialDelayForCreate * time.Second)

	// wait for AssociateId.
	if err := svc.WaitUntilVpnGatewayAvailable(ctx, expandDescribeVpnGatewaysInput(d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for vpngateway available: %s", err))
	}

	return read(ctx, d, meta)
}
