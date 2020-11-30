package volume

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func waitUntilVolumeExtended(ctx context.Context, svc *computing.Client, input *computing.DescribeVolumesInput) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"attaching"},
		Target:     []string{"attached"},
		Refresh:    waitUntilVolumeExtendedState(ctx, svc, input),
		Timeout:    3 * time.Minute,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func waitUntilVolumeExtendedState(ctx context.Context, svc *computing.Client, input *computing.DescribeVolumesInput) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := svc.DescribeVolumesRequest(input).Send(ctx)

		if err != nil {
			return nil, "", err
		}

		if len(res.VolumeSet) == 0 {
			return res, "", nil
		}

		for _, a := range res.VolumeSet[0].AttachmentSet {
			if a.Status == nil {
				fmt.Printf("Ignoring nil attachment state for volume %#v: %v", res.VolumeSet[0].VolumeId, a)
				continue
			}
			return res, nifcloud.StringValue(a.Status), nil
		}
		return res, "", nil
	}
}
