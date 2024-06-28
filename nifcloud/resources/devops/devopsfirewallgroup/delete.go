package devopsfirewallgroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func deleteFirewallGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	input := expandDeleteFirewallGroupInput(d)

	if _, err := svc.DeleteFirewallGroup(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete a DevOps firewall group: %s", err))
	}

	err := waitUntilFirewallGroupDeleted(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait until the DevOps firewall group is deleted: %s", err))
	}

	d.SetId("")

	return nil
}
