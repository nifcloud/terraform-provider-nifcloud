package devopsbackuprule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func deleteBackupRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	input := expandDeleteBackupRuleInput(d)

	if _, err := svc.DeleteBackupRule(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete a DevOps backup rule: %s", err))
	}

	d.SetId("")

	return nil
}
