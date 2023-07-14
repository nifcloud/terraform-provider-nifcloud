package remoteaccessvpngateway

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDescribeRemoteAccessVpnGatewaysInput(d)
	describeClientConfigInput := expandDescribeRemoteAccessVpnGatewayClientConfigInput(d)

	svc := meta.(*client.Client).Computing

	res, err := svc.DescribeRemoteAccessVpnGateways(ctx, input)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.RemoteAccessVpnGatewayId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading remote access vpn gateway: %s", err))
	}

	describeClientConfigRes, err := svc.DescribeRemoteAccessVpnGatewayClientConfig(ctx, describeClientConfigInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed reading remote access vpn gateway client config: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	if err := flattenClientConfig(d, describeClientConfigRes); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
