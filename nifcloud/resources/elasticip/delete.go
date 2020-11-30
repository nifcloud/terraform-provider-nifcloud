package elasticip

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandReleaseAddressInput(d)

	svc := meta.(*client.Client).Computing
	req := svc.ReleaseAddressRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	d.SetId("")
	return nil
}
