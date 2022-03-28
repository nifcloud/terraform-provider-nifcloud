package vpnconnection

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDescribeVpnConnectionsInput(d)
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if d.IsNewResource() {
		err := computing.NewVpnConnectionAvailableWaiter(svc).Wait(ctx, input, time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for vpn connection to become ready: %s", err))
		}
	}

	res, err := svc.DescribeVpnConnections(ctx, input)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.VpnConnectionId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading vpn connection: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
