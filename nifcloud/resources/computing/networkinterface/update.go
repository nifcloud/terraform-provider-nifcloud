package networkinterface

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

const waiterInitialDelayForUpdate = 60

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("description") {
		input := expandModifyNetworkInterfaceAttributeInputForDescription(d)

		req := svc.ModifyNetworkInterfaceAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating network interface description: %s", err))
		}
	}

	if d.HasChange("ip_address") {
		input := expandModifyNetworkInterfaceAttributeInputForIPAddress(d)

		req := svc.ModifyNetworkInterfaceAttributeRequest(input)

		routerSet, err := getRouterSet(ctx, d, svc)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating network interface for get router set: %s", err))
		}

		for _, r := range routerSet {
			mutexKV.Lock(nifcloud.StringValue(r.RouterId))
			defer mutexKV.Unlock(nifcloud.StringValue(r.RouterId))

			if err := svc.WaitUntilRouterAvailable(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{nifcloud.StringValue(r.RouterId)}}); err != nil {
				return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
			}
		}

		_, err = req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating network interface ip address: %s", err))
		}

		waitForNetworkInterfaceAvailable()
	}
	return read(ctx, d, meta)
}

func waitForNetworkInterfaceAvailable() diag.Diagnostics {
	// The status of the NetworkInterface changes shortly after calling the modify API.
	// So, wait a few seconds as initial delay.
	time.Sleep(waiterInitialDelayForUpdate * time.Second)
	return nil
}
