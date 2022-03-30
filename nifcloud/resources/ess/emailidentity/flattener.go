package emailidentity

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
)

func flatten(d *schema.ResourceData, res *ess.GetIdentityVerificationAttributesOutput) error {
	if res == nil || len(res.VerificationAttributes) == 0 {
		d.SetId("")
		return nil
	}

	entry := res.VerificationAttributes[0]

	if err := d.Set("email", entry.Key); err != nil {
		return err
	}
	return nil
}
