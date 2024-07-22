package multiipaddressgroup

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeMultiIpAddressGroupsOutput) error {
	if res == nil || len(res.MultiIpAddressGroupsSet) == 0 {
		d.SetId("")
		return nil
	}

	group := res.MultiIpAddressGroupsSet[0]

	if nifcloud.ToString(group.MultiIpAddressGroupId) != d.Id() {
		return fmt.Errorf("unable to find multi IP address group within: %#v", res.MultiIpAddressGroupsSet)
	}

	if err := d.Set("name", group.MultiIpAddressGroupName); err != nil {
		return err
	}

	if err := d.Set("description", group.Description); err != nil {
		return err
	}

	if err := d.Set("availability_zone", group.AvailabilityZone); err != nil {
		return err
	}

	if group.MultiIpAddressNetwork == nil {
		return nil
	}

	if err := d.Set("default_gateway", group.MultiIpAddressNetwork.DefaultGateway); err != nil {
		return err
	}

	if err := d.Set("subnet_mask", group.MultiIpAddressNetwork.SubnetMask); err != nil {
		return err
	}

	ipAddresses := make([]string, len(group.MultiIpAddressNetwork.IpAddressesSet))
	for i, set := range group.MultiIpAddressNetwork.IpAddressesSet {
		ipAddresses[i] = nifcloud.ToString(set.IpAddress)
	}
	if err := d.Set("ip_addresses", ipAddresses); err != nil {
		return err
	}

	if err := d.Set("ip_address_count", len(ipAddresses)); err != nil {
		return err
	}

	return nil
}
