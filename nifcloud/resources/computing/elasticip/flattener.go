package elasticip

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeAddressesOutput) error {
	if res == nil || len(res.AddressesSet) == 0 {
		d.SetId("")
		return nil
	}

	elasticIP := res.AddressesSet[0]

	if nifcloud.ToString(elasticIP.PrivateIpAddress) == d.Id() {
		if err := d.Set("ip_type", true); err != nil {
			return err
		}
		if err := d.Set("private_ip", elasticIP.PrivateIpAddress); err != nil {
			return err
		}
	} else if nifcloud.ToString(elasticIP.PublicIp) == d.Id() {
		if err := d.Set("ip_type", false); err != nil {
			return err
		}
		if err := d.Set("public_ip", elasticIP.PublicIp); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unable to find elastic ip within: %#v", res.AddressesSet)
	}

	if err := d.Set("availability_zone", elasticIP.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("description", elasticIP.Description); err != nil {
		return err
	}
	return nil
}
