package emailidentity

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
)

func expandVerifyEmailIdentityInput(d *schema.ResourceData) *ess.VerifyEmailIdentityInput {
	return &ess.VerifyEmailIdentityInput{
		EmailAddress: nifcloud.String(d.Get("email").(string)),
	}
}

func expandGetIdentityVerificationAttributesInput(d *schema.ResourceData) *ess.GetIdentityVerificationAttributesInput {
	return &ess.GetIdentityVerificationAttributesInput{
		Identities: []string{d.Id()},
	}
}

func expandDeleteIdentityInput(d *schema.ResourceData) *ess.DeleteIdentityInput {
	return &ess.DeleteIdentityInput{
		Identity: nifcloud.String(d.Id()),
	}
}
