package webproxy

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyDescribeWebProxiesInput(d)

	svc := meta.(*client.Client).Computing

	if d.IsNewResource() {
		err := svc.WaitUntilRouterAvailable(ctx, expandNiftyDescribeRoutersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for router to become ready: %s", err))
		}
	}

	req := svc.NiftyDescribeWebProxiesRequest(input)

	res, err := req.Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.RouterId" {
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
