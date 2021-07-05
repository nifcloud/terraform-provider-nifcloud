package networkinterface

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteNetworkInterfaceInput(d)
	svc := meta.(*client.Client).Computing

	req := svc.DeleteNetworkInterfaceRequest(input)

	if err := waitForRouterOfNetworkInterfaceAvailable(ctx, d, svc); err != nil {
		return err
	}

	_, err := req.Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.NetworkInterfaceId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = svc.WaitUntilPrivateLanAvailable(ctx, &computing.NiftyDescribePrivateLansInput{NetworkId: []string{d.Get("network_id").(string)}})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until private lan available after network interface deleted error: %s", err))
	}

	d.SetId("")
	return nil
}
