package securitygroup

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeSecurityGroupsResponse) error {
	if res == nil || len(res.SecurityGroupInfo) == 0 {
		d.SetId("")
		return nil
	}

	securityGroup := res.SecurityGroupInfo[0]

	if nifcloud.StringValue(securityGroup.GroupName) != d.Id() {
		return fmt.Errorf("unable to find key pair within: %#v", res.SecurityGroupInfo)
	}

	if err := d.Set("group_name", securityGroup.GroupName); err != nil {
		return err
	}

	if err := d.Set("description", securityGroup.GroupDescription); err != nil {
		return err
	}

	if err := d.Set("availability_zone", securityGroup.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("log_limit", securityGroup.GroupLogLimit); err != nil {
		return err
	}
	return nil
}
