package emailidentity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandVerifyEmailIdentityInput(d)

	svc := meta.(*client.Client).ESS
	_, err := svc.VerifyEmailIdentity(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating email: %s", err))
	}

	d.SetId(d.Get("email").(string))

	return read(ctx, d, meta)
}
