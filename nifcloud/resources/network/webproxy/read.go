package webproxy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyDescribeWebProxiesInput(d)

	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if d.IsNewResource() {
		err := computing.NewRouterAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeRoutersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	res, err := svc.NiftyDescribeWebProxies(ctx, input)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.RouterId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
