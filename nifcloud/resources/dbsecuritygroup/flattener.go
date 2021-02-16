package dbsecuritygroup

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
)

func flatten(d *schema.ResourceData, res *rdb.DescribeDBSecurityGroupsResponse) error {
	if res == nil || len(res.DBSecurityGroups) == 0 {
		d.SetId("")
		return nil
	}

	dbSecurityGroup := res.DBSecurityGroups[0]

	if nifcloud.StringValue(dbSecurityGroup.DBSecurityGroupName) != d.Id() {
		return fmt.Errorf("unable to find db security group within: %#v", res.DBSecurityGroups)
	}

	if err := d.Set("group_name", dbSecurityGroup.DBSecurityGroupName); err != nil {
		return err
	}

	var rules []map[string]interface{}

	if len(dbSecurityGroup.IPRanges) != 0 {
		for _, r := range dbSecurityGroup.IPRanges {
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

	if len(dbSecurityGroup.EC2SecurityGroups) != 0 {
		for _, r := range dbSecurityGroup.EC2SecurityGroups {
			var findElm map[string]interface{}

			for _, e := range d.Get("rule").(*schema.Set).List() {
				elm := e.(map[string]interface{})

				if elm["security_group_name"] == nifcloud.StringValue(r.EC2SecurityGroupName) {
					findElm = elm
					break
				}
			}

			if findElm != nil {
				rule := map[string]interface{}{
					"security_group_name": r.EC2SecurityGroupName,
				}
				rules = append(rules, rule)
			}
		}
	}

	if err := d.Set("rule", rules); err != nil {
		return err
	}

	if err := d.Set("availability_zone", dbSecurityGroup.NiftyAvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("description", dbSecurityGroup.DBSecurityGroupDescription); err != nil {
		return err
	}

	return nil
}
