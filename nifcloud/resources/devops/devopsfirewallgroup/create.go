package devopsfirewallgroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func createFirewallGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	input := expandCreateFirewallGroupInput(d)

	res, err := svc.CreateFirewallGroup(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create a DevOps firewall group: %s", err))
	}

	d.SetId(nifcloud.ToString(res.FirewallGroup.FirewallGroupName))

	return updateFirewallGroup(ctx, d, meta)
}
