package elasticip

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandAllocateAddressInput(d *schema.ResourceData) *computing.AllocateAddressInput {
	return &computing.AllocateAddressInput{
		NiftyPrivateIp: nifcloud.Bool(d.Get("ip_type").(bool)),
		Placement: &computing.RequestPlacementOfAllocateAddress{
			AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		},
	}
}

func expandNiftyModifyAddressAttributeInput(d *schema.ResourceData) *computing.NiftyModifyAddressAttributeInput {
	input := &computing.NiftyModifyAddressAttributeInput{}

	if d.Get("ip_type").(bool) {
		input.PrivateIpAddress = nifcloud.String(d.Id())
	}
	if !(d.Get("ip_type").(bool)) {
		input.PublicIp = nifcloud.String(d.Id())
	}

	input.Attribute = "description"
	input.Value = nifcloud.String(d.Get("description").(string))

	return input
}

func expandDescribeAddressesInput(d *schema.ResourceData) *computing.DescribeAddressesInput {
	input := &computing.DescribeAddressesInput{}

	if d.Get("ip_type").(bool) {
		input.PrivateIpAddress = []string{d.Id()}
	}
	if !(d.Get("ip_type").(bool)) {
		input.PublicIp = []string{d.Id()}
	}
	return input
}

func expandReleaseAddressInput(d *schema.ResourceData) *computing.ReleaseAddressInput {
	input := &computing.ReleaseAddressInput{}

	if d.Get("ip_type").(bool) {
		input.PrivateIpAddress = nifcloud.String(d.Id())
	}
	if !(d.Get("ip_type").(bool)) {
		input.PublicIp = nifcloud.String(d.Id())
	}
	return input
}
