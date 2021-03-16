package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.IsNewResource() {
		err := svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for instance to become ready: %s", err))
		}
	}

	if d.HasChange("accounting_type") {
		input := expandModifyInstanceAttributeInputForAccountingType(d)

		req := svc.ModifyInstanceAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance accounting_type: %s", err))
		}

		err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}

	if d.HasChange("description") {
		input := expandModifyInstanceAttributeInputForDescription(d)

		req := svc.ModifyInstanceAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance description: %s", err))
		}

		err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}

	if d.HasChange("disable_api_termination") {
		input := expandModifyInstanceAttributeInputForDisableAPITermination(d)

		req := svc.ModifyInstanceAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance disable_api_termination: %s", err))
		}

		err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}

	if d.HasChange("instance_id") && !d.IsNewResource() {
		input := expandModifyInstanceAttributeInputForInstanceID(d)

		req := svc.ModifyInstanceAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance instance_id %s", err))
		}

		d.SetId(d.Get("instance_id").(string))

		err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}

	if d.HasChange("instance_type") {
		input := expandModifyInstanceAttributeInputForInstanceType(d)

		req := svc.ModifyInstanceAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance instance_type: %s", err))
		}

		err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}

	if d.HasChange("network_interface") {
		routers, err := getRouterList(ctx, d, svc)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance network for get router set: %s", err))
		}

		for _, r := range routers {
			mutexKV.Lock(r)
			defer mutexKV.Unlock(r)

			if err := svc.WaitUntilRouterAvailable(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}); err != nil {
				return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
			}
		}

		input := expandNiftyUpdateInstanceNetworkInterfacesInput(d)

		req := svc.NiftyUpdateInstanceNetworkInterfacesRequest(input)

		_, err = req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance network_interface: %s", err))
		}

		err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}

		for _, r := range routers {
			if err := svc.WaitUntilRouterAvailable(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}); err != nil {
				return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
			}
		}

		o, n := d.GetChange("network_interface")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old network interface and detach any obsolete ones
		for _, n := range ors.List() {
			if attachmentID, ok := n.(map[string]interface{})["network_interface_attachment_id"]; ok && attachmentID != "" {
				input := expandDetachNetworkInterfaceInput(d, attachmentID.(string))
				req := svc.DetachNetworkInterfaceRequest(input)

				_, err := req.Send(ctx)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed updating instance interface to detach network interface: %s", err))
				}

				err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
				}

				for _, r := range routers {
					if err := svc.WaitUntilRouterAvailable(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}); err != nil {
						return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
					}
				}

			}
		}

		// Then loop through all the newly configured network interface and attach them
		for _, n := range nrs.List() {
			if networkInterfaceID, ok := n.(map[string]interface{})["network_interface_id"]; ok && networkInterfaceID != "" {
				input := expandAttachNetworkInterfaceInput(d, networkInterfaceID.(string))
				req := svc.AttachNetworkInterfaceRequest(input)

				_, err := req.Send(ctx)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed updating instance to attach network interface: %s", err))
				}

				err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
				}

				for _, r := range routers {
					if err := svc.WaitUntilRouterAvailable(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}); err != nil {
						return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
					}
				}
			}
		}
	}

	if d.HasChange("security_group") {
		input := expandModifyInstanceAttributeInputForSecurityGroup(d)

		req := svc.ModifyInstanceAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance security_group: %s", err))
		}

		err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}
	return read(ctx, d, meta)
}
