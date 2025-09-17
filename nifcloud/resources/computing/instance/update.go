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

		o, n := d.GetChange("network_interface")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old network interface
		//   network interface: detach any obsolete ones
		//   multi ip address group: disassociate from the instance
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

			if networkID, ok := n.(map[string]interface{})["network_id"]; ok && networkID.(string) == "net-MULTI_IP_ADDRESS" {
				if err := disassociateMultiIPAddressGroup(ctx, d, svc, true); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// Then loop through all the newly configured network interface
		//   network interface: attach any new ones
		//   multi ip address group: associate with the instance
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

			if networkID, ok := n.(map[string]interface{})["network_id"]; ok && networkID.(string) == "net-MULTI_IP_ADDRESS" {
				if err := associateMultiIPAddressGroup(ctx, d, svc, n); err != nil {
					return diag.FromErr(err)
				}
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

func associateMultiIPAddressGroup(ctx context.Context, d *schema.ResourceData, svc *computing.Client, networkInterface interface{}) error {
	deadline, _ := ctx.Deadline()

	id, ok := networkInterface.(map[string]interface{})["multi_ip_address_group_id"]
	if !ok || id.(string) == "" {
		return fmt.Errorf("when specifying net-MULTI_IP_ADDRESS, multi_ip_address_group_id is required")
	}
	multiIPAddressGroupID := id.(string)

	instanceRes, err := svc.DescribeInstances(ctx, expandDescribeInstancesInput(d))
	if err != nil {
		return fmt.Errorf("failed describing instance info: %w", err)
	}

	if len(instanceRes.ReservationSet) == 0 || len(instanceRes.ReservationSet[0].InstancesSet) == 0 {
		return fmt.Errorf("instance %s not found", d.Id())
	}
	instance := instanceRes.ReservationSet[0].InstancesSet[0]

	if instance.MultiIpAddressGroup != nil && nifcloud.ToString(instance.MultiIpAddressGroup.MultiIpAddressGroupId) == multiIPAddressGroupID {
		return nil
	}

	if nifcloud.ToString(instance.InstanceState.Name) != "stopped" {
		if _, err := svc.StopInstances(ctx, expandStopInstancesInput(d, false)); err != nil {
			return fmt.Errorf("failed stopping instance before associating with multi IP address group: %w", err)
		}

		if err := computing.NewInstanceStoppedWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline)); err != nil {
			return fmt.Errorf("failed wait until instance stopped: %w", err)
		}
	}

	mutexKV.Lock(multiIPAddressGroupID)
	defer mutexKV.Unlock(multiIPAddressGroupID)

	input := expandAssociateMultiIpAddressGroupInput(multiIPAddressGroupID, nifcloud.ToString(instance.InstanceUniqueId))
	if _, err := svc.AssociateMultiIpAddressGroup(ctx, input); err != nil {
		return fmt.Errorf("failed associating instance with multi IP address group: %w", err)
	}

	if err := computing.NewInstanceStoppedWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline)); err != nil {
		return fmt.Errorf("failed wait until multi IP address group associated: %w", err)
	}

	if _, err := svc.StartInstances(ctx, expandStartInstancesInputWithMultiIPAddressConfigurationUserData(d)); err != nil {
		return fmt.Errorf("failed starting instance after associating with multi IP address group: %w", err)
	}

	if err := computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline)); err != nil {
		return fmt.Errorf("failed wait until instance running: %w", err)
	}

	return nil
}

func disassociateMultiIPAddressGroup(ctx context.Context, d *schema.ResourceData, svc *computing.Client, startInstanceAfterDisassociate bool) error {
	deadline, _ := ctx.Deadline()

	instanceRes, err := svc.DescribeInstances(ctx, expandDescribeInstancesInput(d))
	if err != nil {
		return fmt.Errorf("failed describing instance info: %w", err)
	}

	if len(instanceRes.ReservationSet) == 0 || len(instanceRes.ReservationSet[0].InstancesSet) == 0 {
		return fmt.Errorf("instance %s not found", d.Id())
	}
	instance := instanceRes.ReservationSet[0].InstancesSet[0]

	if instance.MultiIpAddressGroup != nil && instance.MultiIpAddressGroup.MultiIpAddressGroupId == nil {
		return nil
	}

	if nifcloud.ToString(instance.InstanceState.Name) != "stopped" {
		if _, err := svc.StopInstances(ctx, expandStopInstancesInput(d, false)); err != nil {
			return fmt.Errorf("failed stopping instance before associating with multi IP address group: %w", err)
		}

		if err := computing.NewInstanceStoppedWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline)); err != nil {
			return fmt.Errorf("failed wait until instance stopped: %w", err)
		}
	}

	multiIPAddressGroupID := nifcloud.ToString(instance.MultiIpAddressGroup.MultiIpAddressGroupId)
	mutexKV.Lock(multiIPAddressGroupID)
	defer mutexKV.Unlock(multiIPAddressGroupID)

	input := expandDisassociateMultiIpAddressGroupInput(multiIPAddressGroupID, nifcloud.ToString(instance.InstanceUniqueId))
	if _, err := svc.DisassociateMultiIpAddressGroup(ctx, input); err != nil {
		return fmt.Errorf("failed disassociating instance with multi IP address group: %w", err)
	}

	if err := computing.NewInstanceStoppedWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline)); err != nil {
		return fmt.Errorf("failed wait until multi IP address group disassociated: %w", err)
	}

	if !startInstanceAfterDisassociate {
		return nil
	}

	if _, err := svc.StartInstances(ctx, expandStartInstancesInputWithMultiIPAddressConfigurationUserData(d)); err != nil {
		return fmt.Errorf("failed starting instance after disassociating with multi IP address group: %w", err)
	}

	if err := computing.NewInstanceRunningWaiter(svc).Wait(ctx, expandDescribeInstancesInput(d), time.Until(deadline)); err != nil {
		return fmt.Errorf("failed wait until instance running: %w", err)
	}

	return nil
}
