package networkinterface

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateNetworkInterfaceInput(d)

	svc := meta.(*client.Client).Computing
	req := svc.CreateNetworkInterfaceRequest(input)

	routerSet, err := getRouterSet(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating network interface for get router set: %s", err))
	}

	for _, r := range routerSet {
		mutexKV.Lock(nifcloud.StringValue(r.RouterId))
		defer mutexKV.Unlock(nifcloud.StringValue(r.RouterId))

		if err := svc.WaitUntilRouterAvailable(ctx, &computing.NiftyDescribeRoutersInput{RouterId: []string{nifcloud.StringValue(r.RouterId)}}); err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router available: %s", err))
		}
	}

	res, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating network interface: %s", err))
	}

	d.SetId(nifcloud.StringValue(res.NetworkInterface.NetworkInterfaceId))

	return read(ctx, d, meta)
}
