package zone

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
)

func flatten(d *schema.ResourceData, res *dns.GetHostedZoneResponse) error {
	if res == nil {
		d.SetId("")
		return nil
	}

	hostedZone := res.HostedZone
	delegationSet := res.DelegationSet

	if nifcloud.StringValue(hostedZone.Name) != d.Id() {
		return fmt.Errorf("unable to find hosted zone within: %#v", hostedZone)
	}

	if err := d.Set("name", hostedZone.Name); err != nil {
		return err
	}

	if err := d.Set("comment", hostedZone.Config.Comment); err != nil {
		return err
	}

	if err := d.Set("name_servers", delegationSet.NameServers); err != nil {
		return err
	}

	return nil
}
