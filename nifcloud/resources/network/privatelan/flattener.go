package privatelan

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.NiftyDescribePrivateLansOutput) error {

	if res == nil || len(res.PrivateLanSet) == 0 {
		d.SetId("")
		return nil
	}
	privateLan := res.PrivateLanSet[0]

	if nifcloud.ToString(privateLan.NetworkId) != d.Id() {
		return fmt.Errorf("unable to find private lan within: %#v", res.PrivateLanSet)
	}

	if err := d.Set("accounting_type", privateLan.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("availability_zone", privateLan.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("description", privateLan.Description); err != nil {
		return err
	}

	if err := d.Set("network_id", privateLan.NetworkId); err != nil {
		return err
	}

	if err := d.Set("private_lan_name", privateLan.PrivateLanName); err != nil {
		return err
	}

	if err := d.Set("state", privateLan.State); err != nil {
		return err
	}

	if err := d.Set("cidr_block", privateLan.CidrBlock); err != nil {
		return err
	}

	return nil
}
