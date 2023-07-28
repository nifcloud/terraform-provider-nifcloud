package remoteaccessvpngateway

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
	input := expandCreateRemoteAccessVpnGatewayInput(d)

	svc := meta.(*client.Client).Computing

	niMap := d.Get("network_interface").([]interface{})[0].(map[string]interface{})
	if v, ok := niMap["network_id"].(string); ok && v != "" {
		key, err := mutexkv.LockPrivateLan(ctx, v, svc)
		if err != nil {
			return diag.FromErr(err)
		}
		defer mutexkv.UnlockPrivateLan(key)
	}

	res, err := svc.CreateRemoteAccessVpnGateway(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating remote access vpn gateway: %s", err))
	}

	d.SetId(nifcloud.ToString(res.RemoteAccessVpnGateway.RemoteAccessVpnGatewayId))

	return update(ctx, d, meta)
}
