package domaindkim

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
)

func flatten(d *schema.ResourceData, res *ess.GetIdentityDkimAttributesOutput) error {
	if res == nil || len(res.DkimAttributes) == 0 {
		d.SetId("")
		return nil
	}

	entry := res.DkimAttributes[0]

	if err := d.Set("domain", entry.Key); err != nil {
		return err
	}

	if err := d.Set("dkim_tokens", entry.Value.DkimTokens); err != nil {
		return err
	}

	return nil
}
