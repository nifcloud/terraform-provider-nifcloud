package router

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

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	createRouterInput := expandNiftyCreateRouterInput(d)

	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	for _, ni := range d.Get("network_interface").(*schema.Set).List() {
		if v, ok := ni.(map[string]interface{}); ok {
			if raw, ok := v["network_id"]; ok && len(raw.(string)) > 0 {
				key, err := mutexkv.LockPrivateLan(ctx, raw.(string), svc)
				if err != nil {
					return diag.FromErr(err)
				}
				defer mutexkv.UnlockPrivateLan(key)
			}
			if raw, ok := v["network_name"]; ok && len(raw.(string)) > 0 {
				key, err := mutexkv.LockPrivateLanByName(ctx, raw.(string), svc)
				if err != nil {
					return diag.FromErr(err)
				}
				defer mutexkv.UnlockPrivateLan(key)
			}
		}
	}

	res, err := svc.NiftyCreateRouter(ctx, createRouterInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating router: %s", err))
	}

	d.SetId(nifcloud.ToString(res.Router.RouterId))

	describeRoutersInput := expandNiftyDescribeRoutersInput(d)

	if _, ok := d.GetOk("route_table_id"); ok {
		if err := computing.NewRouterAvailableWaiter(svc).Wait(ctx, describeRoutersInput, time.Until(deadline)); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}

		associateRouteTableInput := expandAssociateRouteTableInput(d)
		if _, err := svc.AssociateRouteTable(ctx, associateRouteTableInput); err != nil {
			return diag.FromErr(fmt.Errorf("failed associating route table to router: %s", err))
		}
	}

	if _, ok := d.GetOk("nat_table_id"); ok {
		if err := computing.NewRouterAvailableWaiter(svc).Wait(ctx, describeRoutersInput, time.Until(deadline)); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}

		associateNatTableInput := expandNiftyAssociateNatTableInput(d)
		if _, err := svc.NiftyAssociateNatTable(ctx, associateNatTableInput); err != nil {
			return diag.FromErr(fmt.Errorf("failed associating NAT table to router: %s", err))
		}
	}

	return read(ctx, d, meta)
}
