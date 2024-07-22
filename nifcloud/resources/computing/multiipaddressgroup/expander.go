package multiipaddressgroup

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandDescribeMultiIPAddressGroupsInput(d *schema.ResourceData) *computing.DescribeMultiIpAddressGroupsInput {
	return &computing.DescribeMultiIpAddressGroupsInput{
		MultiIpAddressGroupId: []string{d.Id()},
	}
}

func expandCreateMultiIPAddressGroupInput(d *schema.ResourceData) *computing.CreateMultiIpAddressGroupInput {
	return &computing.CreateMultiIpAddressGroupInput{
		MultiIpAddressGroupName: nifcloud.String(d.Get("name").(string)),
		Placement: &types.RequestPlacementOfCreateMultiIpAddressGroup{
			AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		},
		IpAddressCount: nifcloud.Int32(int32(d.Get("ip_address_count").(int))),
		Description:    nifcloud.String(d.Get("description").(string)),
	}
}

func expandModifyMultiIPAddressGroupAttributeForNameInput(d *schema.ResourceData) *computing.ModifyMultiIpAddressGroupAttributeInput {
	return &computing.ModifyMultiIpAddressGroupAttributeInput{
		MultiIpAddressGroupId:   nifcloud.String(d.Id()),
		MultiIpAddressGroupName: nifcloud.String(d.Get("name").(string)),
	}
}

func expandModifyMultiIPAddressGroupAttributeForDescriptionInput(d *schema.ResourceData) *computing.ModifyMultiIpAddressGroupAttributeInput {
	return &computing.ModifyMultiIpAddressGroupAttributeInput{
		MultiIpAddressGroupId: nifcloud.String(d.Id()),
		Description:           nifcloud.String(d.Get("description").(string)),
	}
}

func expandIncreaseMultiIpAddressCountInput(d *schema.ResourceData, newCount int) *computing.IncreaseMultiIpAddressCountInput {
	return &computing.IncreaseMultiIpAddressCountInput{
		MultiIpAddressGroupId: nifcloud.String(d.Id()),
		IpAddressCount:        nifcloud.Int32(int32(newCount)),
	}
}

func expandReleaseMultiIpAddressesInput(d *schema.ResourceData, ipAddressesToRelease []string) *computing.ReleaseMultiIpAddressesInput {
	return &computing.ReleaseMultiIpAddressesInput{
		MultiIpAddressGroupId: nifcloud.String(d.Id()),
		IpAddress:             ipAddressesToRelease,
	}
}

func expandDeleteMultiIPAddressGroupInput(d *schema.ResourceData) *computing.DeleteMultiIpAddressGroupInput {
	return &computing.DeleteMultiIpAddressGroupInput{
		MultiIpAddressGroupId: nifcloud.String(d.Id()),
	}
}
