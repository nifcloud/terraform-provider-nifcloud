package zone

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
)

func expandCreateHostedZoneInput(d *schema.ResourceData) *dns.CreateHostedZoneInput {
	return &dns.CreateHostedZoneInput{
		Name: nifcloud.String(d.Get("name").(string)),
		RequestHostedZoneConfig: &dns.RequestHostedZoneConfig{
			Comment: nifcloud.String(d.Get("comment").(string)),
		},
	}
}

func expandGetHostedZoneInput(d *schema.ResourceData) *dns.GetHostedZoneInput {
	return &dns.GetHostedZoneInput{
		ZoneID: nifcloud.String(d.Id()),
	}
}

func expandDeleteHostedZoneInput(d *schema.ResourceData) *dns.DeleteHostedZoneInput {
	return &dns.DeleteHostedZoneInput{
		ZoneID: nifcloud.String(d.Id()),
	}
}
