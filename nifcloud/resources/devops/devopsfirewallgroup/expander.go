package devopsfirewallgroup

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
)

func expandCreateFirewallGroupInput(d *schema.ResourceData) *devops.CreateFirewallGroupInput {
	return &devops.CreateFirewallGroupInput{
		FirewallGroupName: nifcloud.String(d.Get("name").(string)),
		AvailabilityZone:  types.AvailabilityZoneOfCreateFirewallGroupRequest(d.Get("availability_zone").(string)),
		Description:       nifcloud.String(d.Get("description").(string)),
	}
}

func expandUpdateFirewallGroupInput(d *schema.ResourceData) *devops.UpdateFirewallGroupInput {
	return &devops.UpdateFirewallGroupInput{
		FirewallGroupName:        nifcloud.String(d.Id()),
		ChangedFirewallGroupName: nifcloud.String(d.Get("name").(string)),
		Description:              nifcloud.String(d.Get("description").(string)),
	}
}

func expandGetFirewallGroupInput(d *schema.ResourceData) *devops.GetFirewallGroupInput {
	return &devops.GetFirewallGroupInput{
		FirewallGroupName: nifcloud.String(d.Id()),
	}
}

func expandDeleteFirewallGroupInput(d *schema.ResourceData) *devops.DeleteFirewallGroupInput {
	return &devops.DeleteFirewallGroupInput{
		FirewallGroupName: nifcloud.String(d.Id()),
	}
}

func expandAuthorizeFirewallRulesInput(d *schema.ResourceData, rules []types.RequestRules) *devops.AuthorizeFirewallRulesInput {
	return &devops.AuthorizeFirewallRulesInput{
		FirewallGroupName: nifcloud.String(d.Id()),
		Rules:             rules,
	}
}

func expandRevokeFirewallRulesInput(d *schema.ResourceData, ruleIds []string) *devops.RevokeFirewallRulesInput {
	return &devops.RevokeFirewallRulesInput{
		FirewallGroupName: nifcloud.String(d.Id()),
		Ids:               nifcloud.String(strings.Join(ruleIds, ",")),
	}
}
