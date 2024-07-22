package devopsinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func updateInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	if d.HasChanges("instance_type", "firewall_group_name", "description") && !d.IsNewResource() {
		input := expandUpdateInstanceInput(d)

		_, err := svc.UpdateInstance(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update a DevOps instance: %s", err))
		}

		if err := waitUntilInstanceRunning(ctx, d, svc); err != nil {
			return diag.FromErr(fmt.Errorf("failed to wait for the DevOps instance to become ready: %s", err))
		}
	}

	if d.HasChange("disk_size") && !d.IsNewResource() {
		o, n := d.GetChange("disk_size")

		// New value has been validated with ValidateFunc and CustomizeDiff already
		nTimesToInvoke := (n.(int) - o.(int)) / 100

		// Make request multiple times to satisfy configured values as a single ExtendDisk invocation extends 100GB of disk
		for i := 0; i < nTimesToInvoke; i++ {
			input := expandExtendDiskInput(d)

			_, err := svc.ExtendDisk(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to extend a disk of the DevOps instance: %s", err))
			}

			if err := waitUntilInstanceRunning(ctx, d, svc); err != nil {
				return diag.FromErr(fmt.Errorf("failed to wait for the DevOps instance to become ready: %s", err))
			}
		}
	}

	if d.HasChanges("network_id", "private_address") && !d.IsNewResource() {
		input := expandUpdateNetworkInterfaceInput(d)

		_, err := svc.UpdateNetworkInterface(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update a network interface of the DevOps instance: %s", err))
		}

		if err := waitUntilInstanceRunning(ctx, d, svc); err != nil {
			return diag.FromErr(fmt.Errorf("failed to wait for the DevOps instance to become ready: %s", err))
		}
	}

	if d.HasChange("to") {
		input := expandSetupAlertInput(d)

		_, err := svc.SetupAlert(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to setup alert for the DevOps instance: %s", err))
		}

		if err := waitUntilInstanceRunning(ctx, d, svc); err != nil {
			return diag.FromErr(fmt.Errorf("failed to wait for the DevOps instance to become ready: %s", err))
		}
	}

	return readInstance(ctx, d, meta)
}
