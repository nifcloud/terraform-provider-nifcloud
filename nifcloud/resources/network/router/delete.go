package router

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

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	describeRoutersInput := expandNiftyDescribeRoutersInput(d)
	if _, err := svc.NiftyDescribeRouters(ctx, describeRoutersInput); err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.RouterId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading routers: %s", err))
	}

	err := computing.NewRouterAvailableWaiter(svc).Wait(ctx, describeRoutersInput, time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
	}

	deleteRoueterInput := expandNiftyDeleteRouterInput(d)
	if _, err := svc.NiftyDeleteRouter(ctx, deleteRoueterInput); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting router: %s", err))
	}

	if err := computing.NewRouterDeletedWaiter(svc).Wait(ctx, describeRoutersInput, time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for router deleted: %s", err))
	}

	d.SetId("")

	return nil
}
