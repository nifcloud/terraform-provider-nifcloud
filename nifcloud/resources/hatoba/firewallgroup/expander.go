package firewallgroup

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
)

func expandCreateFirewallGroupInput(d *schema.ResourceData) *hatoba.CreateFirewallGroupInput {
	return &hatoba.CreateFirewallGroupInput{
		FirewallGroup: &hatoba.CreateFirewallGroupRequestFirewallGroup{
			Name:        nifcloud.String(d.Get("name").(string)),
			Description: nifcloud.String(d.Get("description").(string)),
		},
	}
}

func expandAuthorizeFirewallGroupInput(d *schema.ResourceData, rule map[string]interface{}) *hatoba.AuthorizeFirewallGroupInput {
	r := hatoba.AuthorizeFirewallGroupRequestFirewallRule{
		Protocol:    nifcloud.String(rule["protocol"].(string)),
		Direction:   nifcloud.String(rule["direction"].(string)),
		CidrIp:      nifcloud.String(rule["cidr_ip"].(string)),
		Description: nifcloud.String(rule["description"].(string)),
	}

	if rule["from_port"] != nil && rule["from_port"] != 0 {
		r.FromPort = nifcloud.Int64(int64(rule["from_port"].(int)))
	}
	if rule["to_port"] != nil && rule["to_port"] != 0 {
		r.ToPort = nifcloud.Int64(int64(rule["to_port"].(int)))
	}

	return &hatoba.AuthorizeFirewallGroupInput{
		FirewallGroupName: nifcloud.String(d.Id()),
		Rules:             []hatoba.AuthorizeFirewallGroupRequestFirewallRule{r},
	}
}

func expandRevokeFirewallGroupInput(d *schema.ResourceData, rule map[string]interface{}) *hatoba.RevokeFirewallGroupInput {
	return &hatoba.RevokeFirewallGroupInput{
		FirewallGroupName: nifcloud.String(d.Id()),
		Ids:               nifcloud.String(rule["id"].(string)),
	}
}

func expandGetFirewallGroupInput(d *schema.ResourceData) *hatoba.GetFirewallGroupInput {
	return &hatoba.GetFirewallGroupInput{
		FirewallGroupName: nifcloud.String(d.Id()),
	}
}

func expandUpdateFirewallGroupInput(d *schema.ResourceData) *hatoba.UpdateFirewallGroupInput {
	input := &hatoba.UpdateFirewallGroupInput{
		FirewallGroupName: nifcloud.String(d.Id()),
		FirewallGroup: &hatoba.UpdateFirewallGroupRequestFirewallGroup{
			Description: nifcloud.String(d.Get("description").(string)),
		},
	}

	if d.HasChange("name") && !d.IsNewResource() {
		input.FirewallGroup.Name = nifcloud.String(d.Get("name").(string))
	}

	return input
}

func expandDeleteFirewallGroupInput(d *schema.ResourceData) *hatoba.DeleteFirewallGroupInput {
	return &hatoba.DeleteFirewallGroupInput{
		FirewallGroupName: nifcloud.String(d.Id()),
	}
}
