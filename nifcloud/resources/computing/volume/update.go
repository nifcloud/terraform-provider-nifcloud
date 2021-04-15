package volume

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.HasChange("accounting_type") {
		input := expandModifyVolumeAttributeInputForAccountingType(d)

		req := svc.ModifyVolumeAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating volume accounting_type: %s", err))
		}
	}

	if d.HasChange("volume_id") {
		input := expandModifyVolumeAttributeInputForVolumeName(d)

		req := svc.ModifyVolumeAttributeRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating volume volume_id %s", err))
		}

		d.SetId(d.Get("volume_id").(string))
	}

	if d.HasChange("description") {
		input := expandModifyVolumeAttributeInputForDescription(d)

		req := svc.ModifyVolumeAttributeRequest(input)

		_, err := req.Send(ctx)
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
			req := svc.ExtendVolumeSizeRequest(extendVolumeInput)
			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed extending volume size: %s", err))
			}

			describeVolumeInput := expandDescribeVolumesInput(d)

			err = waitUntilVolumeExtended(ctx, svc, describeVolumeInput)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed extending volume size: %s", err))
			}

			res, err := svc.DescribeVolumesRequest(describeVolumeInput).Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed reading: %s", err))
			}

			extendSize, err := strconv.Atoi(nifcloud.StringValue(res.VolumeSet[0].Size))
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

		// Check if the instance has changed
		res, err := svc.DescribeInstancesRequest(
			&computing.DescribeInstancesInput{
				InstanceId: []string{afterID.(string)},
			},
		).Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed reading instance: %s", err))
		}

		afterUniqueID := res.DescribeInstancesOutput.ReservationSet[0].InstancesSet[0].InstanceUniqueId
		if *afterUniqueID == d.Get("instance_unique_id").(string) {
			return read(ctx, d, meta)
		}

		detachVolumeInput := expandDetachVolumeInput(d)
		detachVolumeInput.InstanceId = nifcloud.String(beforeID.(string))
		_, err = svc.DetachVolumeRequest(detachVolumeInput).Send(ctx)
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

		attachVolumeInput := expandAttachVolumeInput(d)
		_, err = svc.AttachVolumeRequest(attachVolumeInput).Send(ctx)
		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.Volume" {
				d.SetId("")
				return nil
			}
			return diag.FromErr(fmt.Errorf("failed attaching volume: %s", err))
		}

		err = svc.WaitUntilVolumeInUse(ctx, describeVolumeInput)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for volume attached: %s", err))
		}
	}

	return read(ctx, d, meta)
}
