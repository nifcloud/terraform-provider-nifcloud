package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

var mutexKV = mutexkv.NewMutexKV()

func getRouterList(ctx context.Context, d *schema.ResourceData, svc *computing.Client) ([]string, error) {
	routers := []types.RouterSetOfNiftyDescribePrivateLans{}
	result := []string{}

	networkIDs := []string{}
	describeInstancesRes, err := svc.DescribeInstances(ctx, expandDescribeInstancesInput(d))
	if err != nil {
		return result, err
	}

	for _, ni := range describeInstancesRes.ReservationSet[0].InstancesSet[0].NetworkInterfaceSet {
		networkID := nifcloud.ToString(ni.NiftyNetworkId)
		if networkID != "net-COMMON_GLOBAL" && networkID != "net-COMMON_PRIVATE" && networkID != "net-MULTI_IP_ADDRESS" {
			networkIDs = append(networkIDs, networkID)
		}
	}

	if len(networkIDs) == 0 {
		return result, nil
	}

	input := &computing.NiftyDescribePrivateLansInput{
		NetworkId: networkIDs,
	}

	res, err := svc.NiftyDescribePrivateLans(ctx, input)
	if err != nil {
		return result, err
	}

	if res == nil || len(res.PrivateLanSet) == 0 {
		return result, nil
	}

	for _, privateLan := range res.PrivateLanSet {
		routers = append(routers, privateLan.RouterSet...)
	}

	m := make(map[string]struct{})
	for _, router := range routers {
		routerID := nifcloud.ToString(router.RouterId)
		if _, ok := m[routerID]; !ok {
			m[routerID] = struct{}{}
			result = append(result, routerID)
		}
	}

	return result, nil
}
