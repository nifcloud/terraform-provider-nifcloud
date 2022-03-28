package nasinstance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteNASInstanceInput(d)
	svc := meta.(*client.Client).NAS
	deadline, _ := ctx.Deadline()
	_, err := svc.DeleteNASInstance(ctx, input)

	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameter.NotFound.NASInstanceIdentifer" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting NAS instance: %s", err))
	}

	if err := nas.NewNASInstanceDeletedWaiter(svc).Wait(ctx, expandDescribeNASInstancesInput(d), time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted NAS instance error: %s", err))
	}

	d.SetId("")

	return nil
}
