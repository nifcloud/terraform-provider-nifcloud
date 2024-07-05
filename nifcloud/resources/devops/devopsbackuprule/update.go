package devopsbackuprule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func updateBackupRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	if d.HasChanges("name", "description") {
		input := expandUpdateBackupRuleInput(d)

		if _, err := svc.UpdateBackupRule(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update a DevOps backup rule: %s", err))
		}
	}

	if d.HasChange("name") {
		d.SetId(d.Get("name").(string))
	}

	return readBackupRule(ctx, d, meta)
}
