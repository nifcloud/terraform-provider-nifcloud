package devopsinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func deleteInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	input := expandDeleteInstanceInput(d)

	if _, err := svc.DeleteInstance(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete a DevOps instance: %s", err))
	}

	err := waitUntilInstanceDeleted(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait until the DevOps instance is deleted: %s", err))
	}

	d.SetId("")

	return nil
}
