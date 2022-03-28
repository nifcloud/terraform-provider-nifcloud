package webproxy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyDeleteWebProxyInput(d)
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	_, err := svc.NiftyDeleteWebProxy(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = computing.NewRouterAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeRoutersInput(d), time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
	}

	d.SetId("")
	return nil
}
