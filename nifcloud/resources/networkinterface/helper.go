package networkinterface

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

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
