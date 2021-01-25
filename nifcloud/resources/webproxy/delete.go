package webproxy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyDeleteWebProxyInput(d)
	svc := meta.(*client.Client).Computing

	req := svc.NiftyDeleteWebProxyRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
	}

	d.SetId("")
	return nil
}
