package emailidentity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteIdentityInput(d)

	svc := meta.(*client.Client).ESS
	_, err := svc.DeleteIdentity(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting email error: %s", err))
	}

	d.SetId("")
	return nil
}
