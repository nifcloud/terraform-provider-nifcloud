package nasinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).NAS

	if d.IsNewResource() {
		err := svc.WaitUntilNASInstanceAvailable(ctx, expandDescribeNASInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until NAS instance available: %s", err))
		}
	}

	input := expandModifyNASInstanceInput(d)
	req := svc.ModifyNASInstanceRequest(input)
	if _, err := req.Send(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("failed updating NAS instance: %s", err))
	}

	d.SetId(d.Get("identifier").(string))

	if err := svc.WaitUntilNASInstanceAvailable(ctx, expandDescribeNASInstancesInput(d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for NAS instance to become ready: %s", err))
	}

	return read(ctx, d, meta)
}
