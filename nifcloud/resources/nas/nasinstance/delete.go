package nasinstance

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
	input := expandDeleteNASInstanceInput(d)
	svc := meta.(*client.Client).NAS
	req := svc.DeleteNASInstanceRequest(input)

	if _, err := req.Send(ctx); err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameter.NotFound.NASInstanceIdentifer" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting NAS instance: %s", err))
	}

	if err := svc.WaitUntilNASInstanceDeleted(ctx, expandDescribeNASInstancesInput(d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted NAS instance error: %s", err))
	}

	d.SetId("")

	return nil
}
