package elasticip

import (
	"net"

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

	ip := net.ParseIP(d.Id())
	if ip.IsPrivate() {
		input.PrivateIpAddress = nifcloud.String(ip.String())
	} else {
		input.PublicIp = nifcloud.String(ip.String())
	}

	input.Attribute = "description"
	input.Value = nifcloud.String(d.Get("description").(string))
	return input
}

func expandDescribeAddressesInput(d *schema.ResourceData) *computing.DescribeAddressesInput {
	input := &computing.DescribeAddressesInput{}

	ip := net.ParseIP(d.Id())
	if ip.IsPrivate() {
		input.PrivateIpAddress = []string{ip.String()}
	} else {
		input.PublicIp = []string{ip.String()}
	}
	return input
}

func expandReleaseAddressInput(d *schema.ResourceData) *computing.ReleaseAddressInput {
	input := &computing.ReleaseAddressInput{}

	ip := net.ParseIP(d.Id())
	if ip.IsPrivate() {
		input.PrivateIpAddress = nifcloud.String(ip.String())
	} else {
		input.PublicIp = nifcloud.String(ip.String())
	}
	return input
}
