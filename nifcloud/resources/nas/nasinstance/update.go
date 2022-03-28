package nasinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).NAS
	deadline, _ := ctx.Deadline()

	if d.IsNewResource() {
		err := nas.NewNASInstanceAvailableWaiter(svc).Wait(ctx, expandDescribeNASInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until NAS instance available: %s", err))
		}
	}

	input := expandModifyNASInstanceInput(d)
	_, err := svc.ModifyNASInstance(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed updating NAS instance: %s", err))
	}

	d.SetId(d.Get("identifier").(string))

	if err := nas.NewNASInstanceAvailableWaiter(svc).Wait(ctx, expandDescribeNASInstancesInput(d), time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for NAS instance to become ready: %s", err))
	}

	return read(ctx, d, meta)
}
