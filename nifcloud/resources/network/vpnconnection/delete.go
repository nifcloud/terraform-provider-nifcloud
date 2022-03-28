package vpnconnection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteVpnConnectionInput(d)
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	_, err := svc.DeleteVpnConnection(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting vpn connection: %s", err))
	}

	describeVpnConnectionsInput := expandDescribeVpnConnectionsInput(d)
	err = computing.NewVpnConnectionDeletedWaiter(svc).Wait(ctx, describeVpnConnectionsInput, time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for vpn connection deleted: %s", err))
	}

	d.SetId("")
	return nil
}
