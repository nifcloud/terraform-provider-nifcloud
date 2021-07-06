package firewallgroup

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
)

func flatten(d *schema.ResourceData, res *hatoba.GetFirewallGroupResponse) error {
	if res == nil {
		d.SetId("")
		return nil
	}

	firewallGroup := res.FirewallGroup

	if nifcloud.StringValue(firewallGroup.Name) != d.Id() {
		return fmt.Errorf("unable to find Hatoba firewall group within: %#v", firewallGroup.Name)
	}

	if err := d.Set("name", firewallGroup.Name); err != nil {
		return err
	}

	var rules []map[string]interface{}

	if len(firewallGroup.Rules) != 0 {
		for _, r := range firewallGroup.Rules {
			var findElm map[string]interface{}

			for _, e := range d.Get("rule").(*schema.Set).List() {
				elm := e.(map[string]interface{})
				if elm["id"] == nifcloud.StringValue(r.Id) {
					findElm = elm
					break
				}
			}
			rule := map[string]interface{}{
				"cidr_ip":     r.CidrIp,
				"id":          r.Id,
				"protocol":    r.Protocol,
				"direction":   r.Direction,
				"from_port":   r.FromPort,
				"description": r.Description,
			}

			if findElm != nil {
				if findElm["to_port"] != nil && findElm["to_port"].(int) != 0 {
					rule["to_port"] = r.ToPort
				}
			} else {
				if nifcloud.Int64Value(r.FromPort) != nifcloud.Int64Value(r.ToPort) {
					rule["to_port"] = r.ToPort
				}
			}
			rules = append(rules, rule)
		}
	}

	if err := d.Set("rule", rules); err != nil {
		return err
	}

	if err := d.Set("description", firewallGroup.Description); err != nil {
		return err
	}

	return nil
}
