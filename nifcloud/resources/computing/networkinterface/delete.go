package networkinterface

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
	input := expandDeleteNetworkInterfaceInput(d)
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if err := waitForRouterOfNetworkInterfaceAvailable(ctx, d, svc); err != nil {
		return err
	}

	_, err := svc.DeleteNetworkInterface(ctx, input)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.NetworkInterfaceId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = computing.NewPrivateLanAvailableWaiter(svc).Wait(ctx, &computing.NiftyDescribePrivateLansInput{NetworkId: []string{d.Get("network_id").(string)}}, time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until private lan available after network interface deleted error: %s", err))
	}

	d.SetId("")
	return nil
}
