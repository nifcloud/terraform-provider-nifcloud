package networkinterface

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateNetworkInterfaceInput(d)

	svc := meta.(*client.Client).Computing

	if raw, ok := d.GetOk("network_id"); ok && len(raw.(string)) > 0 {
		key, err := mutexkv.LockPrivateLan(ctx, raw.(string), svc)
		if err != nil {
			return diag.FromErr(err)
		}
		defer mutexkv.UnlockPrivateLan(key)
	}

	req := svc.CreateNetworkInterfaceRequest(input)

	if err := waitForRouterOfNetworkInterfaceAvailable(ctx, d, svc); err != nil {
		return err
	}

	res, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating network interface: %s", err))
	}

	d.SetId(nifcloud.StringValue(res.NetworkInterface.NetworkInterfaceId))

	return read(ctx, d, meta)
}
