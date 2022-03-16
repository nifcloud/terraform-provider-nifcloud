package zone

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteHostedZoneInput(d)

	svc := meta.(*client.Client).DNS
	_, err := svc.DeleteHostedZone(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting hosted zone error: %s", err))
	}

	d.SetId("")

	return nil
}
