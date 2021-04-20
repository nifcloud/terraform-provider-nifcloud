package volume

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeVolumesResponse) error {
	if res == nil || len(res.VolumeSet) == 0 {
		d.SetId("")
		return nil
	}

	volume := res.VolumeSet[0]

	if nifcloud.StringValue(volume.VolumeId) != d.Id() {
		return fmt.Errorf("unable to find volume within: %#v", res.VolumeSet)
	}
	if err := d.Set("volume_id", volume.VolumeId); err != nil {
		return err
	}

	volumeSize, err := strconv.Atoi(nifcloud.StringValue(volume.Size))
	if err != nil {
		return fmt.Errorf("failed converting volume size")
	}
	if err := d.Set("size", volumeSize); err != nil {
		return err
	}

	if err := d.Set("disk_type", volume.DiskType); err != nil {
		return err
	}

	if err := d.Set("accounting_type", volume.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("description", volume.Description); err != nil {
		return err
	}

	if len(res.VolumeSet[0].AttachmentSet) != 0 {
		instance := res.VolumeSet[0].AttachmentSet[0]

		if _, ok := d.GetOk("instance_id"); ok {
			if err := d.Set("instance_id", instance.InstanceId); err != nil {
				return err
			}
		}

		if _, ok := d.GetOk("instance_unique_id"); ok {
			if err := d.Set("instance_unique_id", instance.InstanceUniqueId); err != nil {
				return err
			}
		}
	}

	return nil
}
