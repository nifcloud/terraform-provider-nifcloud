package keypair

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandImportKeyPairInput(d *schema.ResourceData) *computing.ImportKeyPairInput {
	return &computing.ImportKeyPairInput{
		KeyName:           nifcloud.String(d.Get("key_name").(string)),
		PublicKeyMaterial: nifcloud.String(d.Get("public_key").(string)),
		Description:       nifcloud.String(d.Get("description").(string)),
	}
}

func expandNiftyModifyKeyPairAttributeInput(d *schema.ResourceData) *computing.NiftyModifyKeyPairAttributeInput {
	return &computing.NiftyModifyKeyPairAttributeInput{
		KeyName:   nifcloud.String(d.Id()),
		Attribute: "description",
		Value:     nifcloud.String(d.Get("description").(string)),
	}
}
