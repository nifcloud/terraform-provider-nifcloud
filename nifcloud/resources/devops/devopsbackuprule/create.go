package devopsbackuprule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func createBackupRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	input := expandCreateBackupRuleInput(d)

	res, err := svc.CreateBackupRule(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create a DevOps backup rule: %s", err))
	}

	d.SetId(nifcloud.ToString(res.BackupRule.BackupRuleName))

	return updateBackupRule(ctx, d, meta)
}
