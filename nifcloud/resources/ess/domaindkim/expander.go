package domaindkim

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
)

func expandVerifyDomainDkimInput(d *schema.ResourceData) *ess.VerifyDomainDkimInput {
	return &ess.VerifyDomainDkimInput{
		Domain: nifcloud.String(d.Get("domain").(string)),
	}
}

func expandGetIdentityDkimAttributesInput(d *schema.ResourceData) *ess.GetIdentityDkimAttributesInput {
	return &ess.GetIdentityDkimAttributesInput{
		Identities: []string{d.Id()},
	}
}
