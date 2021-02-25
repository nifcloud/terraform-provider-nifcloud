package vpngateway

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	describeVpnGatewaysInput := expandDescribeVpnGatewaysInput(d)
	if _, err := svc.DescribeVpnGatewaysRequest(describeVpnGatewaysInput).Send(ctx); err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.VpnGatewayId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading vpngateways: %s", err))
	}

	deleteVpnGatewayInput := expandDeleteVpnGatewaysInput(d)
	if _, err := svc.DeleteVpnGatewayRequest(deleteVpnGatewayInput).Send(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting vpngateway: %s", err))
	}

	if err := svc.WaitUntilVpnGatewayDeleted(ctx, describeVpnGatewaysInput); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for vpngateway deleted: %s", err))
	}

	d.SetId("")

	return nil
}
