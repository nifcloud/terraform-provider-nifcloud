package router

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	describeRoutersInput := expandNiftyDescribeRoutersInput(d)
	if _, err := svc.NiftyDescribeRoutersRequest(describeRoutersInput).Send(ctx); err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.RouterId" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading routers: %s", err))
	}

	deleteRoueterInput := expandNiftyDeleteRouterInput(d)
	if _, err := svc.NiftyDeleteRouterRequest(deleteRoueterInput).Send(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting router: %s", err))
	}

	if err := svc.WaitUntilRouterDeleted(ctx, describeRoutersInput); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for router deleted: %s", err))
	}

	d.SetId("")

	return nil
}
