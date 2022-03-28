package vpngateway

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

	describeVpnGatewaysInput := expandDescribeVpnGatewaysInput(d)
	if _, err := svc.DescribeVpnGateways(ctx, describeVpnGatewaysInput); err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.VpnGatewayId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading vpngateways: %s", err))
	}

	deleteVpnGatewayInput := expandDeleteVpnGatewayInput(d)
	if _, err := svc.DeleteVpnGateway(ctx, deleteVpnGatewayInput); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting vpngateway: %s", err))
	}

	if err := computing.NewVpnGatewayDeletedWaiter(svc).Wait(ctx, describeVpnGatewaysInput, time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for vpngateway deleted: %s", err))
	}

	d.SetId("")

	return nil
}
