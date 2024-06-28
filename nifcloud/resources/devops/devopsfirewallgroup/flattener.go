package devopsfirewallgroup

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
)

func flatten(d *schema.ResourceData, res *devops.GetFirewallGroupOutput) error {
	if res == nil || res.FirewallGroup == nil {
		d.SetId("")
		return nil
	}

	group := res.FirewallGroup

	if nifcloud.ToString(group.FirewallGroupName) != d.Id() {
		return fmt.Errorf("unable to find the DevOps firewall group within: %#v", group)
	}

	if err := d.Set("name", group.FirewallGroupName); err != nil {
		return err
	}

	if err := d.Set("availability_zone", group.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("description", group.Description); err != nil {
		return err
	}

	var rules []map[string]interface{}
	for _, r := range group.Rules {
		rule := map[string]interface{}{
			"id":          r.Id,
			"protocol":    r.Protocol,
			"port":        r.Port,
			"cidr_ip":     r.CidrIp,
			"description": r.Description,
		}
		rules = append(rules, rule)
	}

	if err := d.Set("rule", rules); err != nil {
		return err
	}

	return nil
}
