package webproxy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("listen_interface_network_name") {
		input := expandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkName(d)

		req := svc.NiftyModifyWebProxyAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy listen_interface_network_name: %s", err))
		}

		err = svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("listen_interface_network_id") {
		input := expandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkID(d)

		req := svc.NiftyModifyWebProxyAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy listen_interface_network_id: %s", err))
		}

		err = svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("bypass_interface_network_name") {
		input := expandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkName(d)

		req := svc.NiftyModifyWebProxyAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_name: %s", err))
		}

		err = svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("bypass_interface_network_id") {
		input := expandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkID(d)

		req := svc.NiftyModifyWebProxyAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_id: %s", err))
		}

		err = svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("name_server") {
		input := expandNiftyModifyWebProxyAttributeInputForNameServer(d)

		req := svc.NiftyModifyWebProxyAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_id: %s", err))
		}

		err = svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("description") {
		input := expandNiftyModifyWebProxyAttributeInputForDescription(d)

		req := svc.NiftyModifyWebProxyAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_id: %s", err))
		}

		err = svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	if d.HasChange("listen_port") {
		input := expandNiftyModifyWebProxyAttributeInputForListenPort(d)

		req := svc.NiftyModifyWebProxyAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating web proxy bypass_interface_network_id: %s", err))
		}

		err = svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	return read(ctx, d, meta)
}
