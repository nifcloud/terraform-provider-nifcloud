package webproxy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if d.HasChange("listen_interface_network_name") {
		input := expandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkName(d)

		_, err := svc.NiftyModifyWebProxyAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy listen_interface_network_name: %s", err))
		}

		err = computing.NewRouterAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeRoutersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("listen_interface_network_id") {
		input := expandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkID(d)

		_, err := svc.NiftyModifyWebProxyAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy listen_interface_network_id: %s", err))
		}

		err = computing.NewRouterAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeRoutersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("bypass_interface_network_name") {
		input := expandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkName(d)

		_, err := svc.NiftyModifyWebProxyAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_name: %s", err))
		}

		err = computing.NewRouterAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeRoutersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("bypass_interface_network_id") {
		input := expandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkID(d)

		_, err := svc.NiftyModifyWebProxyAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_id: %s", err))
		}

		err = computing.NewRouterAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeRoutersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("name_server") {
		input := expandNiftyModifyWebProxyAttributeInputForNameServer(d)

		_, err := svc.NiftyModifyWebProxyAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_id: %s", err))
		}

		err = computing.NewRouterAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeRoutersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("description") {
		input := expandNiftyModifyWebProxyAttributeInputForDescription(d)

		_, err := svc.NiftyModifyWebProxyAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_id: %s", err))
		}

		err = computing.NewRouterAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeRoutersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("listen_port") {
		input := expandNiftyModifyWebProxyAttributeInputForListenPort(d)

		_, err := svc.NiftyModifyWebProxyAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_id: %s", err))
		}

		err = computing.NewRouterAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeRoutersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	return read(ctx, d, meta)
}
