package router

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

const waiterInitialDelay = 3

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("accounting_type") {
		input := expandNiftyModifyRouterAttributeInputForAccountingType(d)

		req := svc.NiftyModifyRouterAttributeRequest(input)
		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating router accounting_type: %s", err))
		}

		if d := waitForRouterAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("description") {
		input := expandNiftyModifyRouterAttributeInputForDescription(d)

		req := svc.NiftyModifyRouterAttributeRequest(input)
		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating router description: %s", err))
		}

		if d := waitForRouterAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("name") {
		input := expandNiftyModifyRouterAttributeInputForRouterName(d)

		req := svc.NiftyModifyRouterAttributeRequest(input)

		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating router name %s", err))
		}

		if d := waitForRouterAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("type") {
		input := expandNiftyModifyRouterAttributeInputForType(d)

		req := svc.NiftyModifyRouterAttributeRequest(input)

		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating router type: %s", err))
		}

		if d := waitForRouterAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("network_interface") {
		input := expandNiftyUpdateRouterNetworkInterfacesInput(d)

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

		req := svc.NiftyUpdateRouterNetworkInterfacesRequest(input)
		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating router network_interface: %s", err))
		}

		if d := waitForRouterAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("security_group") {
		before, after := d.GetChange("security_group")

		if before != "" && after == "" {
			input := expandNiftyDeregisterRoutersFromSecurityGroupInput(d)

			req := svc.NiftyDeregisterRoutersFromSecurityGroupRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed deregistering router security_group: %s", err))
			}
		} else {
			input := expandNiftyModifyRouterAttributeInputForSecurityGroup(d)

			req := svc.NiftyModifyRouterAttributeRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed updating router security_group: %s", err))
			}
		}

		if d := waitForRouterAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("nat_table_id") {
		before, after := d.GetChange("nat_table_id")
		if before != "" && after == "" {
			input := expandNiftyDisassociateNatTableInput(d)

			req := svc.NiftyDisassociateNatTableRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed disassociating NAT table: %s", err))
			}
		} else if before == "" && after != "" {
			input := expandNiftyAssociateNatTableInput(d)

			req := svc.NiftyAssociateNatTableRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed associating NAT table: %s", err))
			}
		} else {
			input := expandNiftyReplaceNatTableAssociationInput(d)

			req := svc.NiftyReplaceNatTableAssociationRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed updating router nat_table_id: %s", err))
			}
		}

		if d := waitForRouterAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("route_table_id") {
		before, after := d.GetChange("route_table_id")
		if before != "" && after == "" {
			input := expandDisassociateRouteTableInput(d)

			req := svc.DisassociateRouteTableRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed disassociating route table: %s", err))
			}
		} else if before == "" && after != "" {
			input := expandAssociateRouteTableInput(d)

			req := svc.AssociateRouteTableRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed associating route table: %s", err))
			}
		} else {
			input := expandReplaceRouteTableAssociation(d)

			req := svc.ReplaceRouteTableAssociationRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed updating router route_table_id: %s", err))
			}
		}

		if d := waitForRouterAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	return read(ctx, d, meta)
}

func waitForRouterAvailable(ctx context.Context, d *schema.ResourceData, svc *computing.Client) diag.Diagnostics {
	// lintignore:R018
	time.Sleep(waiterInitialDelay * time.Second)

	if err := svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
	}

	return nil
}
