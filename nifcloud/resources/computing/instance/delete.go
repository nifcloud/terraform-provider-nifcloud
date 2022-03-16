package instance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	describeInstancesInput := expandDescribeInstancesInput(d)
	describeInstancesRes, err := svc.DescribeInstances(ctx, describeInstancesInput)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Instance" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting for describe instances error: %s", err))
	}

	instance := describeInstancesRes.ReservationSet[0].InstancesSet[0]

	if nifcloud.ToString(instance.InstanceState.Name) != "stopped" {
		stopInstancesInput := expandStopInstancesInput(d)
		_, err := svc.StopInstances(ctx, stopInstancesInput)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed deleting for stop instances error: %s", err))
		}

		err = computing.NewInstanceStoppedWaiter(svc).Wait(ctx, describeInstancesInput, time.Until(deadline))
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

		if err := computing.NewRouterAvailableWaiter(svc).Wait(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{r}}, time.Until(deadline)); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
		}
	}

	terminateInstancesInput := expandTerminateInstancesInput(d)
	_, err = svc.TerminateInstances(ctx, terminateInstancesInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for terminate instances error: %s", err))
	}

	err = computing.NewInstanceDeletedWaiter(svc).Wait(ctx, describeInstancesInput, time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted instances error: %s", err))
	}

	d.SetId("")
	return nil
}
