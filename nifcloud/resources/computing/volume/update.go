package volume

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if d.HasChange("accounting_type") {
		input := expandModifyVolumeAttributeInputForAccountingType(d)

		_, err := svc.ModifyVolumeAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating volume accounting_type: %s", err))
		}
	}

	if d.HasChange("volume_id") {
		input := expandModifyVolumeAttributeInputForVolumeName(d)

		_, err := svc.ModifyVolumeAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating volume volume_id %s", err))
		}

		d.SetId(d.Get("volume_id").(string))
	}

	if d.HasChange("description") {
		input := expandModifyVolumeAttributeInputForDescription(d)

		_, err := svc.ModifyVolumeAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating volume description: %s", err))
		}
	}

	if d.HasChange("size") {
		beforeSize, afterSize := d.GetChange("size")

		if afterSize.(int)%100 != 0 {
			return diag.FromErr(fmt.Errorf("could not resize volume because the size is outside the specifiable range"))
		}
		if afterSize.(int)-beforeSize.(int) < 0 {
			return diag.FromErr(fmt.Errorf("could not resize volume because it is smaller than current size"))
		}

		// NIFCLOUD ExtendVolumeSize API can only grow in size by 100GiB.
		// so, it loops until volume size reached the target size.
		for {
			extendVolumeInput := expandExtendVolumeSizeInput(d)
			_, err := svc.ExtendVolumeSize(ctx, extendVolumeInput)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed extending volume size: %s", err))
			}

			describeVolumeInput := expandDescribeVolumesInput(d)

			err = computing.NewVolumeAttachedWaiter(svc).Wait(ctx, describeVolumeInput, time.Until(deadline))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed extending volume size: %s", err))
			}

			res, err := svc.DescribeVolumes(ctx, describeVolumeInput)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed reading: %s", err))
			}

			extendSize, err := strconv.Atoi(nifcloud.ToString(res.VolumeSet[0].Size))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed convert volume size: %s", err))
			}

			if extendSize == afterSize.(int) {
				break
			}
		}
	}

	if d.HasChange("instance_id") {
		beforeID, afterID := d.GetChange("instance_id")

		describeVolumeInput := expandDescribeVolumesInput(d)

		if beforeID != "" {
			detachVolumeInput := expandDetachVolumeInput(d)
			detachVolumeInput.InstanceId = nifcloud.String(beforeID.(string))
			_, err := svc.DetachVolume(ctx, detachVolumeInput)
			if err != nil {
				var awsErr smithy.APIError
				if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Volume" {
					d.SetId("")
					return nil
				}
				return diag.FromErr(fmt.Errorf("failed detaching volume: %s", err))
			}

			err = computing.NewVolumeAvailableWaiter(svc).Wait(ctx, describeVolumeInput, time.Until(deadline))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed waiting for volume detached: %s", err))
			}
		}

		if afterID != "" {
			attachVolumeInput := expandAttachVolumeInput(d)
			_, err := svc.AttachVolume(ctx, attachVolumeInput)
			if err != nil {
				var awsErr smithy.APIError
				if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Volume" {
					d.SetId("")
					return nil
				}
				return diag.FromErr(fmt.Errorf("failed attaching volume: %s", err))
			}

			err = computing.NewVolumeInUseWaiter(svc).Wait(ctx, describeVolumeInput, time.Until(deadline))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed waiting for volume attached: %s", err))
			}
		}
	}

	return read(ctx, d, meta)
}
