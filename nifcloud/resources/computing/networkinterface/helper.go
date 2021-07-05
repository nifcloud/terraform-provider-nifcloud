package networkinterface

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

const waiterInitialDelay = 3

var mutexKV = mutexkv.NewMutexKV()

func getRouterSet(ctx context.Context, d *schema.ResourceData, svc *computing.Client) ([]computing.RouterSetOfNiftyDescribePrivateLans, error) {
	result := []computing.RouterSetOfNiftyDescribePrivateLans{}

	input := expandNiftyDescribePrivateLansInput(d)

	req := svc.NiftyDescribePrivateLansRequest(input)

	res, err := req.Send(ctx)
	if err != nil {
		return result, err
	}

	if res == nil || len(res.PrivateLanSet) == 0 {
		return result, nil
	}

	result = res.PrivateLanSet[0].RouterSet
	return result, nil
}

func waitForRouterOfNetworkInterfaceAvailable(ctx context.Context, d *schema.ResourceData, svc *computing.Client) diag.Diagnostics {
	// lintignore:R018
	time.Sleep(waiterInitialDelay * time.Second)

	routerSet, err := getRouterSet(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed updating network interface for get router set: %s", err))
	}

	for _, r := range routerSet {
		mutexKV.Lock(nifcloud.StringValue(r.RouterId))
		defer mutexKV.Unlock(nifcloud.StringValue(r.RouterId))

		if err := svc.WaitUntilRouterAvailable(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{nifcloud.StringValue(r.RouterId)}}); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
		}
	}

	return nil
}
