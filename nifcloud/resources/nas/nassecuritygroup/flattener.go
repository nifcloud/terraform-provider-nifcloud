package nassecuritygroup

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
)

func flatten(d *schema.ResourceData, res *nas.DescribeNASSecurityGroupsResponse) error {
	if res == nil || len(res.NASSecurityGroups) == 0 {
		d.SetId("")
		return nil
	}

	nasSecurityGroup := res.NASSecurityGroups[0]

	if nifcloud.StringValue(nasSecurityGroup.NASSecurityGroupName) != d.Id() {
		return fmt.Errorf("unable to find NAS security group within: %#v", res.NASSecurityGroups)
	}

	if err := d.Set("group_name", nasSecurityGroup.NASSecurityGroupName); err != nil {
		return err
	}

	var rules []map[string]interface{}

	if len(nasSecurityGroup.IPRanges) != 0 {
		for _, r := range nasSecurityGroup.IPRanges {
			var findElm map[string]interface{}

			for _, e := range d.Get("rule").(*schema.Set).List() {
				elm := e.(map[string]interface{})

				if elm["cidr_ip"] == nifcloud.StringValue(r.CIDRIP) {
					findElm = elm
					break
				}
			}

			if findElm != nil {
				rule := map[string]interface{}{
					"cidr_ip": r.CIDRIP,
				}
				rules = append(rules, rule)
			}
		}
	}

	if len(nasSecurityGroup.SecurityGroups) != 0 {
		for _, r := range nasSecurityGroup.SecurityGroups {
			var findElm map[string]interface{}

			for _, e := range d.Get("rule").(*schema.Set).List() {
				elm := e.(map[string]interface{})

				if elm["security_group_name"] == nifcloud.StringValue(r.SecurityGroupName) {
					findElm = elm
					break
				}
			}

			if findElm != nil {
				rule := map[string]interface{}{
					"security_group_name": r.SecurityGroupName,
				}
				rules = append(rules, rule)
			}
		}
	}

	if err := d.Set("rule", rules); err != nil {
		return err
	}

	if err := d.Set("availability_zone", nasSecurityGroup.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("description", nasSecurityGroup.NASSecurityGroupDescription); err != nil {
		return err
	}

	return nil
}
