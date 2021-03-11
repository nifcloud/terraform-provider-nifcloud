package networkinterface

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandCreateNetworkInterfaceInput(d *schema.ResourceData) *computing.CreateNetworkInterfaceInput {
	input := &computing.CreateNetworkInterfaceInput{
		NiftyNetworkId: nifcloud.String(d.Get("network_id").(string)),
		IpAddress:      nifcloud.String(d.Get("ip_address").(string)),
		Description:    nifcloud.String(d.Get("description").(string)),
		Placement: &computing.RequestPlacementOfCreateNetworkInterface{
			AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		},
	}
	return input
}

func expandDescribeNetworkInterfacesInput(d *schema.ResourceData) *computing.DescribeNetworkInterfacesInput {
	input := &computing.DescribeNetworkInterfacesInput{
		NetworkInterfaceId: []string{d.Id()},
	}
	return input
}

func expandModifyNetworkInterfaceAttributeInputForDescription(d *schema.ResourceData) *computing.ModifyNetworkInterfaceAttributeInput {
	return &computing.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: nifcloud.String(d.Id()),
		Description:        nifcloud.String(d.Get("description").(string)),
	}
}

func expandModifyNetworkInterfaceAttributeInputForIPAddress(d *schema.ResourceData) *computing.ModifyNetworkInterfaceAttributeInput {
	return &computing.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: nifcloud.String(d.Id()),
		IpAddress:          nifcloud.String(d.Get("ip_address").(string)),
	}
}

func expandDeleteNetworkInterfaceInput(d *schema.ResourceData) *computing.DeleteNetworkInterfaceInput {
	input := &computing.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: nifcloud.String(d.Id()),
	}
	return input
}

func expandNiftyDescribePrivateLansInput(d *schema.ResourceData) *computing.NiftyDescribePrivateLansInput {
	return &computing.NiftyDescribePrivateLansInput{
		NetworkId: []string{d.Get("network_id").(string)},
	}
}
