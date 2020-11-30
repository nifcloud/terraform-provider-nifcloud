package volume

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

	detachVolumeInput := expandDetachVolumeInput(d)
	_, err := svc.DetachVolumeRequest(detachVolumeInput).Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.Volume" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed detaching volume: %s", err))
	}

	describeVolumeInput := expandDescribeVolumesInput(d)
	err = svc.WaitUntilVolumeAvailable(ctx, describeVolumeInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for volume detached: %s", err))
	}

	deleteVolumeInput := expandDeleteVolumeInput(d)
	_, err = svc.DeleteVolumeRequest(deleteVolumeInput).Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting volume: %s", err))
	}

	err = svc.WaitUntilVolumeDeleted(ctx, describeVolumeInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for volume deleted: %s", err))
	}

	d.SetId("")
	return nil
}
