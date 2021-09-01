package router

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	createRouterInput := expandNiftyCreateRouterInput(d)

	svc := meta.(*client.Client).Computing

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

	createRouterReq := svc.NiftyCreateRouterRequest(createRouterInput)

	res, err := createRouterReq.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating router: %s", err))
	}

	d.SetId(nifcloud.StringValue(res.Router.RouterId))

	describeRoutersInput := expandNiftyDescribeRoutersInput(d)

	if _, ok := d.GetOk("route_table_id"); ok {
		if err := svc.WaitUntilRouterAvailable(ctx, describeRoutersInput); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}

		associateRouteTableInput := expandAssociateRouteTableInput(d)
		associateRouteTableReq := svc.AssociateRouteTableRequest(associateRouteTableInput)
		if _, err := associateRouteTableReq.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed associating route table to router: %s", err))
		}
	}

	if _, ok := d.GetOk("nat_table_id"); ok {
		if err := svc.WaitUntilRouterAvailable(ctx, describeRoutersInput); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}

		associateNatTableInput := expandNiftyAssociateNatTableInput(d)
		associateNatTableReq := svc.NiftyAssociateNatTableRequest(associateNatTableInput)
		if _, err := associateNatTableReq.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed associating NAT table to router: %s", err))
		}
	}

	return read(ctx, d, meta)
}
