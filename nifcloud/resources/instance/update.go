package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

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

	if d.HasChange("instance_id") {
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
		input := expandNiftyUpdateInstanceNetworkInterfacesInput(d)

		req := svc.NiftyUpdateInstanceNetworkInterfacesRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance network_interface: %s", err))
		}

		err = svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
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
