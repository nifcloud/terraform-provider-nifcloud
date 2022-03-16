package vpngateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

const waiterInitialDelayForCreate = 3

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateVpnGatewayInput(d)
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if raw, ok := d.GetOk("network_name"); ok && len(raw.(string)) > 0 {
		key, err := mutexkv.LockPrivateLanByName(ctx, raw.(string), svc)
		if err != nil {
			return diag.FromErr(err)
		}
		defer mutexkv.UnlockPrivateLan(key)
	}
	if raw, ok := d.GetOk("network_id"); ok && len(raw.(string)) > 0 {
		key, err := mutexkv.LockPrivateLan(ctx, raw.(string), svc)
		if err != nil {
			return diag.FromErr(err)
		}
		defer mutexkv.UnlockPrivateLan(key)
	}

	res, err := svc.CreateVpnGateway(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating VPN Gateway: %s", err))
	}

	d.SetId(nifcloud.ToString(res.VpnGateway.VpnGatewayId))

	describeVpnGatewaysInput := expandDescribeVpnGatewaysInput(d)

	if _, ok := d.GetOk("route_table_id"); ok {
		if err := computing.NewVpnGatewayAvailableWaiter(svc).Wait(ctx, describeVpnGatewaysInput, time.Until(deadline)); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for vpngateway to become ready: %s", err))
		}
		niftyAssociateRouteTableWithVpnGatewayInput := expandNiftyAssociateRouteTableWithVpnGatewayInput(d)
		if _, err := svc.NiftyAssociateRouteTableWithVpnGateway(ctx, niftyAssociateRouteTableWithVpnGatewayInput); err != nil {
			return diag.FromErr(fmt.Errorf("failed associating route table to vpngateway: %s", err))
		}
	}

	// lintignore:R018
	time.Sleep(waiterInitialDelayForCreate * time.Second)

	// wait for AssociateId.
	if err := computing.NewVpnGatewayAvailableWaiter(svc).Wait(ctx, expandDescribeVpnGatewaysInput(d), time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for vpngateway available: %s", err))
	}

	return read(ctx, d, meta)
}
