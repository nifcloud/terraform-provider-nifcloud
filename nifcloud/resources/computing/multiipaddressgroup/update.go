package multiipaddressgroup

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

	if d.HasChange("name") {
		input := expandModifyMultiIPAddressGroupAttributeForNameInput(d)
		if _, err := svc.ModifyMultiIpAddressGroupAttribute(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating multi IP address group name: %s", err))
		}
	}

	waiter := computing.NewMultiIpAddressGroupAvailableWaiter(svc)
	if err := waiter.Wait(ctx, expandDescribeMultiIPAddressGroupsInput(d), 5*time.Minute); err != nil {
		return diag.Errorf("failed waiting multi IP address group to be available: %s", err)
	}

	if d.HasChange("description") {
		input := expandModifyMultiIPAddressGroupAttributeForDescriptionInput(d)
		if _, err := svc.ModifyMultiIpAddressGroupAttribute(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating multi IP address group description: %s", err))
		}
	}

	if err := waiter.Wait(ctx, expandDescribeMultiIPAddressGroupsInput(d), 5*time.Minute); err != nil {
		return diag.Errorf("failed waiting multi IP address group to be available: %s", err)
	}

	if d.HasChange("ip_address_count") {
		o, n := d.GetChange("ip_address_count")
		diff := n.(int) - o.(int)
		if diff < 0 {
			releaseCount := -1 * diff
			releaseTargets := make([]string, releaseCount)
			ipAddresses := d.Get("ip_addresses").([]interface{})
			for i := 0; i < releaseCount; i++ {
				releaseTargets[i] = ipAddresses[len(ipAddresses)-1-i].(string)
			}

			input := expandReleaseMultiIpAddressesInput(d, releaseTargets)
			if _, err := svc.ReleaseMultiIpAddresses(ctx, input); err != nil {
				return diag.FromErr(fmt.Errorf("failed updating multi IP address count (release count: %d): %s", diff, err))
			}
		} else {
			input := expandIncreaseMultiIpAddressCountInput(d, diff)
			if _, err := svc.IncreaseMultiIpAddressCount(ctx, input); err != nil {
				return diag.FromErr(fmt.Errorf("failed updating multi IP address count (increase count: %d): %s", diff, err))
			}
		}
	}

	if err := waiter.Wait(ctx, expandDescribeMultiIPAddressGroupsInput(d), 5*time.Minute); err != nil {
		return diag.Errorf("failed waiting multi IP address group to be available: %s", err)
	}

	return read(ctx, d, meta)
}
