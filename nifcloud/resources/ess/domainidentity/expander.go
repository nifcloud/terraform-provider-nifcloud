package domainidentity

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
)

func expandVerifyDomainIdentityInput(d *schema.ResourceData) *ess.VerifyDomainIdentityInput {
	return &ess.VerifyDomainIdentityInput{
		Domain: nifcloud.String(d.Get("domain").(string)),
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
