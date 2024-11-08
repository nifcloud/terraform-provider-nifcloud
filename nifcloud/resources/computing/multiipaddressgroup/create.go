package multiipaddressgroup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	input := expandCreateMultiIPAddressGroupInput(d)
	res, err := svc.CreateMultiIpAddressGroup(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating multi IP address group: %s", err))
	}

	d.SetId(nifcloud.ToString(res.MultiIpAddressGroup.MultiIpAddressGroupId))

	waiter := computing.NewMultiIpAddressGroupAvailableWaiter(svc)
	if err := waiter.Wait(ctx, expandDescribeMultiIPAddressGroupsInput(d), 5*time.Minute); err != nil {
		return diag.Errorf("failed waiting multi IP address group to be available: %s", err)
	}

	return read(ctx, d, meta)
}
