package mutexkv

import (
	"context"
	"fmt"
	"time"

	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

var privateLan = NewMutexKV()

func LockPrivateLan(ctx context.Context, id string, svc *computing.Client) (string, error) {
	privateLan.Lock(id)

	deadline, _ := ctx.Deadline()
	if id == "net-COMMON_PRIVATE" || id == "net-COMMON_GLOBAL" {
		return id, nil
	}

	return id, computing.NewPrivateLanAvailableWaiter(svc).Wait(ctx, &computing.NiftyDescribePrivateLansInput{NetworkId: []string{id}}, time.Until(deadline))
}

func LockPrivateLanByName(ctx context.Context, name string, svc *computing.Client) (string, error) {
	res, err := svc.NiftyDescribePrivateLans(ctx, &computing.NiftyDescribePrivateLansInput{PrivateLanName: []string{name}})
	if err != nil {
		return "", err
	}

	if res == nil || len(res.PrivateLanSet) != 1 {
		return "", fmt.Errorf("the privateLan not found: %s", name)
	}

	id := nifcloud.ToString(res.PrivateLanSet[0].NetworkId)

	return LockPrivateLan(ctx, id, svc)
}

func UnlockPrivateLan(id string) {
	privateLan.Unlock(id)
}
