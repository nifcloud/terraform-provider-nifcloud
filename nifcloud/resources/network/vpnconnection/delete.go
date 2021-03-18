package vpnconnection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteVpnConnectionInput(d)
	svc := meta.(*client.Client).Computing
	req := svc.DeleteVpnConnectionRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting vpn connection: %s", err))
	}

	describeVpnConnectionsInput := expandDescribeVpnConnectionsInput(d)
	err = svc.WaitUntilVpnConnectionDeleted(ctx, describeVpnConnectionsInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for vpn connection deleted: %s", err))
	}

	d.SetId("")
	return nil
}
