package remoteaccessvpngateway

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

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	describeRemoteAccessVpnGatewaysInput := expandDescribeRemoteAccessVpnGatewaysInput(d)
	if _, err := svc.DescribeRemoteAccessVpnGateways(ctx, describeRemoteAccessVpnGatewaysInput); err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.RemoteAccessVpnGatewayId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading remote access vpn gateways: %s", err))
	}

	err := computing.NewRemoteAccessVpnGatewayAvailableWaiter(svc).Wait(ctx, describeRemoteAccessVpnGatewaysInput, time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for remote access vpn gateway to become ready: %s", err))
	}

	input := expandDeleteRemoteAccessVpnGatewayInput(d)
	if _, err := svc.DeleteRemoteAccessVpnGateway(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting remote access vpn gateway: %s", err))
	}

	if err := computing.NewRemoteAccessVpnGatewayDeletedWaiter(svc).Wait(ctx, describeRemoteAccessVpnGatewaysInput, time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for remote access vpn gateway deleted: %s", err))
	}

	d.SetId("")

	return nil
}
