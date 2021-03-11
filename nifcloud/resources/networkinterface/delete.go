package networkinterface

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteNetworkInterfaceInput(d)
	svc := meta.(*client.Client).Computing

	req := svc.DeleteNetworkInterfaceRequest(input)

	routerSet, err := getRouterSet(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting network interface for get router set: %s", err))
	}

	for _, r := range routerSet {
		mutexKV.Lock(nifcloud.StringValue(r.RouterId))
		defer mutexKV.Unlock(nifcloud.StringValue(r.RouterId))

		if err := svc.WaitUntilRouterAvailable(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{nifcloud.StringValue(r.RouterId)}}); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
		}
	}

	_, err = req.Send(ctx)
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
