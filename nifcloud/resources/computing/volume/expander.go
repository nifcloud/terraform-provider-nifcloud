package volume

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

const (
	// Available disk types.
	// doc: https://pfs.nifcloud.com/service/disk.htm

	// VolumeTypeStandard represents a general purpose volume.
	volumeTypeStandard = "Standard Storage"
	// VolumeTypeHighSpeedA represents a high speed volume (only use type A).
	volumeTypeHighSpeedA = "High-Speed Storage A"
	// VolumeTypeHighSpeedB represents a high speed volume (only use type B).
	volumeTypeHighSpeedB = "High-Speed Storage B"
	// VolumeTypeFlash represents a flash volume.
	volumeTypeFlash = "Flash Storage"
	// VolumeTypeStandardFlashA represents a standard flash volume (only use type A).
	volumeTypeStandardFlashA = "Standard Flash Storage A"
	// VolumeTypeStandardFlashB represents a standard flash volume (only use type B).
	volumeTypeStandardFlashB = "Standard Flash Storage B"
	// VolumeTypeHighSpeedFlashA represents a high spped flash volume (only use type A).
	volumeTypeHighSpeedFlashA = "High-Speed Flash Storage A"
	// VolumeTypeHighSpeedFlashB represents a high spped flash volume (only use type B).
	volumeTypeHighSpeedFlashB = "High-Speed Flash Storage B"
)

var (
	// volumeTypeMapping converts the volume identifier from volume type.
	// More info: https://pfs.nifcloud.com/api/rest/CreateVolume.htm
	volumeTypeMapping = map[string]types.DiskTypeOfCreateVolumeRequest{
		volumeTypeStandard:        types.DiskTypeOfCreateVolumeRequestStandardStorage,
		volumeTypeHighSpeedA:      types.DiskTypeOfCreateVolumeRequestHighSpeedStorageA,
		volumeTypeHighSpeedB:      types.DiskTypeOfCreateVolumeRequestHighSpeedStorageB,
		volumeTypeFlash:           types.DiskTypeOfCreateVolumeRequestFlashDrive,
		volumeTypeStandardFlashA:  types.DiskTypeOfCreateVolumeRequestStandardFlashStorageA,
		volumeTypeStandardFlashB:  types.DiskTypeOfCreateVolumeRequestStandardFlashStorageB,
		volumeTypeHighSpeedFlashA: types.DiskTypeOfCreateVolumeRequestHighSpeedFlashStorageA,
		volumeTypeHighSpeedFlashB: types.DiskTypeOfCreateVolumeRequestHighSpeedFlashStorageB,
	}
)

func expandCreateVolumeInput(d *schema.ResourceData) *computing.CreateVolumeInput {
	input := &computing.CreateVolumeInput{
		Size:           nifcloud.Int32(int32(d.Get("size").(int))),
		VolumeId:       nifcloud.String(d.Get("volume_id").(string)),
		DiskType:       volumeTypeMapping[d.Get("disk_type").(string)],
		AccountingType: types.AccountingTypeOfCreateVolumeRequest(d.Get("accounting_type").(string)),
		Description:    nifcloud.String(d.Get("description").(string)),
	}

	if len(d.Get("instance_id").(string)) != 0 {
		input.InstanceId = nifcloud.String(d.Get("instance_id").(string))
	}
	if len(d.Get("instance_unique_id").(string)) != 0 {
		input.InstanceUniqueId = nifcloud.String(d.Get("instance_unique_id").(string))
	}

	return input
}

func expandModifyVolumeAttributeInputForAccountingType(d *schema.ResourceData) *computing.ModifyVolumeAttributeInput {
	return &computing.ModifyVolumeAttributeInput{
		VolumeId:  nifcloud.String(d.Id()),
		Attribute: types.AttributeOfModifyVolumeAttributeRequestAccountingType,
		Value:     nifcloud.String(d.Get("accounting_type").(string)),
	}
}

func expandModifyVolumeAttributeInputForVolumeName(d *schema.ResourceData) *computing.ModifyVolumeAttributeInput {
	before, after := d.GetChange("volume_id")

	return &computing.ModifyVolumeAttributeInput{
		VolumeId:  nifcloud.String(before.(string)),
		Attribute: types.AttributeOfModifyVolumeAttributeRequestVolumeName,
		Value:     nifcloud.String(after.(string)),
	}
}

func expandModifyVolumeAttributeInputForDescription(d *schema.ResourceData) *computing.ModifyVolumeAttributeInput {
	return &computing.ModifyVolumeAttributeInput{
		VolumeId:  nifcloud.String(d.Id()),
		Attribute: types.AttributeOfModifyVolumeAttributeRequestDescription,
		Value:     nifcloud.String(d.Get("description").(string)),
	}
}

func expandExtendVolumeSizeInput(d *schema.ResourceData) *computing.ExtendVolumeSizeInput {
	return &computing.ExtendVolumeSizeInput{
		VolumeId:    nifcloud.String(d.Id()),
		NiftyReboot: types.NiftyRebootOfExtendVolumeSizeRequest(d.Get("reboot").(string)),
	}
}

func expandDescribeVolumesInput(d *schema.ResourceData) *computing.DescribeVolumesInput {
	return &computing.DescribeVolumesInput{
		VolumeId: []string{d.Id()},
	}
}

func expandAttachVolumeInput(d *schema.ResourceData) *computing.AttachVolumeInput {
	return &computing.AttachVolumeInput{
		VolumeId:   nifcloud.String(d.Id()),
		InstanceId: nifcloud.String(d.Get("instance_id").(string)),
	}
}

func expandDetachVolumeInput(d *schema.ResourceData) *computing.DetachVolumeInput {
	return &computing.DetachVolumeInput{
		VolumeId:   nifcloud.String(d.Id()),
		InstanceId: nifcloud.String(d.Get("instance_id").(string)),
		Agreement:  nifcloud.Bool(true),
	}
}

func expandDeleteVolumeInput(d *schema.ResourceData) *computing.DeleteVolumeInput {
	return &computing.DeleteVolumeInput{
		VolumeId: nifcloud.String(d.Id()),
	}
}
