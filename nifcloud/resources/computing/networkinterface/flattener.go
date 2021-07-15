package networkinterface

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeNetworkInterfacesResponse) error {
	if res == nil || len(res.NetworkInterfaceSet) == 0 {
		d.SetId("")
		return nil
	}

	networkInterface := res.NetworkInterfaceSet[0]

	if nifcloud.StringValue(networkInterface.NetworkInterfaceId) != d.Id() {
		return fmt.Errorf(
			"unable to find network interface within: %#v",
			res.NetworkInterfaceSet,
		)
	}

	if err := d.Set("network_interface_id", networkInterface.NetworkInterfaceId); err != nil {
		return err
	}

	if err := d.Set("availability_zone", networkInterface.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("description", networkInterface.Description); err != nil {
		return err
	}

	if err := d.Set("network_id", networkInterface.NiftyNetworkId); err != nil {
		return err
	}

	if raw, ok := d.GetOk("ip_address"); ok {
		if raw == "static" {
			if err := d.Set("ip_address", "static"); err != nil {
				return err
			}
		} else {
			if err := d.Set("ip_address", networkInterface.PrivateIpAddress); err != nil {
				return err
			}
		}
	} else {
		if err := d.Set("ip_address", networkInterface.PrivateIpAddress); err != nil {
			return err
		}
	}

	if err := d.Set("private_ip", networkInterface.PrivateIpAddress); err != nil {
		return err
	}
	return nil
}
