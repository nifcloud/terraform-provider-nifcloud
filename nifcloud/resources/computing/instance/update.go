package instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if d.IsNewResource() {
		err := computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for instance to become ready: %s", err))
		}
	}

	if d.HasChange("accounting_type") {
		input := expandModifyInstanceAttributeInputForAccountingType(d)

		_, err := svc.ModifyInstanceAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance accounting_type: %s", err))
		}

		err = computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}

	if d.HasChange("description") {
		input := expandModifyInstanceAttributeInputForDescription(d)

		_, err := svc.ModifyInstanceAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance description: %s", err))
		}

		err = computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}

	if d.HasChange("disable_api_termination") {
		input := expandModifyInstanceAttributeInputForDisableAPITermination(d)

		_, err := svc.ModifyInstanceAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance disable_api_termination: %s", err))
		}

		err = computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}

	if d.HasChange("instance_id") && !d.IsNewResource() {
		input := expandModifyInstanceAttributeInputForInstanceID(d)

		_, err := svc.ModifyInstanceAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance instance_id %s", err))
		}

		d.SetId(d.Get("instance_id").(string))

		err = computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}

	if d.HasChange("instance_type") {
		input := expandModifyInstanceAttributeInputForInstanceType(d)

		_, err := svc.ModifyInstanceAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance instance_type: %s", err))
		}

		err = computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
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

			if err := computing.NewRouterAvailableWaiter(svc).Wait(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}, time.Until(deadline)); err != nil {
				return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
			}
		}

		input := expandNiftyUpdateInstanceNetworkInterfacesInput(d)

		_, err = svc.NiftyUpdateInstanceNetworkInterfaces(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating instance network_interface: %s", err))
		}

		err = computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}

		for _, r := range routers {
			if err := computing.NewRouterAvailableWaiter(svc).Wait(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}, time.Until(deadline)); err != nil {
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
				_, err := svc.DetachNetworkInterface(ctx, input)

				if err != nil {
					return diag.FromErr(fmt.Errorf("failed updating instance interface to detach network interface: %s", err))
				}

				err = computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
				}

				for _, r := range routers {
					if err := computing.NewRouterAvailableWaiter(svc).Wait(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}, time.Until(deadline)); err != nil {
						return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
					}
				}

			}
		}

		// Then loop through all the newly configured network interface and attach them
		for _, n := range nrs.List() {
			if networkInterfaceID, ok := n.(map[string]interface{})["network_interface_id"]; ok && networkInterfaceID != "" {
				input := expandAttachNetworkInterfaceInput(d, networkInterfaceID.(string))
				_, err := svc.AttachNetworkInterface(ctx, input)

				if err != nil {
					return diag.FromErr(fmt.Errorf("failed updating instance to attach network interface: %s", err))
				}

				err = computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
				}

				for _, r := range routers {
					if err := computing.NewRouterAvailableWaiter(svc).Wait(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}, time.Until(deadline)); err != nil {
						return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
					}
				}
			}
		}
	}

	if d.HasChange("security_group") {
		if d.Get("security_group").(string) == "" {
			input := expandDeregisterInstancesFromSecurityGroupInput(d)

			mutexKV.Lock(nifcloud.ToString(input.GroupName))
			defer mutexKV.Unlock(nifcloud.ToString(input.GroupName))

			err := computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{nifcloud.ToString(input.GroupName)}}, time.Until(deadline))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
			}

			_, err = svc.DeregisterInstancesFromSecurityGroup(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating instance security_group: %s", err))
			}
		} else {
			input := expandModifyInstanceAttributeInputForSecurityGroup(d)

			mutexKV.Lock(nifcloud.ToString(input.Value))
			defer mutexKV.Unlock(nifcloud.ToString(input.Value))

			err := computing.NewSecurityGroupAppliedWaiter(svc).Wait(ctx, &computing.DescribeSecurityGroupsInput{GroupName: []string{nifcloud.ToString(input.Value)}}, time.Until(deadline))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed wait until securityGroup applied: %s", err))
			}
			_, err = svc.ModifyInstanceAttribute(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating instance security_group: %s", err))
			}
		}

		err := computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until instance running: %s", err))
		}
	}
	return read(ctx, d, meta)
}
