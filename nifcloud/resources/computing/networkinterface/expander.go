package networkinterface

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandCreateNetworkInterfaceInput(d *schema.ResourceData) *computing.CreateNetworkInterfaceInput {
	return &computing.CreateNetworkInterfaceInput{
		NiftyNetworkId: nifcloud.String(d.Get("network_id").(string)),
		IpAddress:      nifcloud.String(d.Get("ip_address").(string)),
		Description:    nifcloud.String(d.Get("description").(string)),
		Placement: &types.RequestPlacementOfCreateNetworkInterface{
			AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		},
	}
}

func expandDescribeNetworkInterfacesInput(d *schema.ResourceData) *computing.DescribeNetworkInterfacesInput {
	return &computing.DescribeNetworkInterfacesInput{
		NetworkInterfaceId: []string{d.Id()},
	}
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
	return &computing.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: nifcloud.String(d.Id()),
	}
}

func expandNiftyDescribePrivateLansInput(d *schema.ResourceData) *computing.NiftyDescribePrivateLansInput {
	return &computing.NiftyDescribePrivateLansInput{
		NetworkId: []string{d.Get("network_id").(string)},
	}
}
