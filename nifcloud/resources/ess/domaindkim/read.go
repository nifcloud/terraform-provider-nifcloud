package domaindkim

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandGetIdentityDkimAttributesInput(d)
	svc := meta.(*client.Client).ESS
	res, err := svc.GetIdentityDkimAttributes(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed reading domain dkim: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
