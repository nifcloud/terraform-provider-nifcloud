package keypair

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("description") {
		input := expandNiftyModifyKeyPairAttributeInput(d)

		svc := meta.(*client.Client).Computing
		_, err := svc.NiftyModifyKeyPairAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating KeyPair: %s", err))
		}
	}
	return read(ctx, d, meta)
}
