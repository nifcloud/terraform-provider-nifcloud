package router

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
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
	// The status of the router changes shortly after calling the modify API.
	// So, wait a few seconds as initial delay.
	time.Sleep(waiterInitialDelay * time.Second)

	if err := svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
	}

	return nil
}
