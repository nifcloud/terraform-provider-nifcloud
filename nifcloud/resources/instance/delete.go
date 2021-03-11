package instance

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
	svc := meta.(*client.Client).Computing

	describeInstancesInput := expandDescribeInstancesInput(d)
	describeInstancesRes, err := svc.DescribeInstancesRequest(describeInstancesInput).Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.Instance" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting for describe instances error: %s", err))
	}

	instance := describeInstancesRes.ReservationSet[0].InstancesSet[0]

	if nifcloud.StringValue(instance.InstanceState.Name) != "stopped" {
		stopInstancesInput := expandStopInstancesInput(d)
		_, err := svc.StopInstancesRequest(stopInstancesInput).Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed deleting for stop instances error: %s", err))
		}

		err = svc.WaitUntilInstanceStopped(ctx, describeInstancesInput)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed deleting for wait until stopped instances error: %s", err))
		}
	}

	routers, err := getRouterList(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed updating instance network for get router set: %s", err))
	}

	for _, r := range routers {
		mutexKV.Lock(r)
		defer mutexKV.Unlock(r)

		if err := svc.WaitUntilRouterAvailable(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
		}
	}

	terminateInstancesInput := expandTerminateInstancesInput(d)
	_, err = svc.TerminateInstancesRequest(terminateInstancesInput).Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for terminate instances error: %s", err))
	}

	err = svc.WaitUntilInstanceDeleted(ctx, describeInstancesInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted instances error: %s", err))
	}

	d.SetId("")
	return nil
}
