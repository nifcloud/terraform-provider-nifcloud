package dhcpconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyDeleteDhcpConfigInput(d)
	svc := meta.(*client.Client).Computing

	_, err := svc.NiftyDeleteDhcpConfig(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	d.SetId("")
	return nil
}
