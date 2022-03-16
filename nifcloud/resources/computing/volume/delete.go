package volume

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

	detachVolumeInput := expandDetachVolumeInput(d)
	_, err := svc.DetachVolume(ctx, detachVolumeInput)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Volume" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed detaching volume: %s", err))
	}

	describeVolumeInput := expandDescribeVolumesInput(d)
	err = computing.NewVolumeAvailableWaiter(svc).Wait(ctx, describeVolumeInput, time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for volume detached: %s", err))
	}

	deleteVolumeInput := expandDeleteVolumeInput(d)
	_, err = svc.DeleteVolume(ctx, deleteVolumeInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting volume: %s", err))
	}

	err = computing.NewVolumeDeletedWaiter(svc).Wait(ctx, describeVolumeInput, time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for volume deleted: %s", err))
	}

	d.SetId("")
	return nil
}
