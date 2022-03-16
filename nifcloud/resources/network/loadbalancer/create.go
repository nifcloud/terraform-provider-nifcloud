package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateLoadBalancerInput(d)

	svc := meta.(*client.Client).Computing

	_, err := svc.CreateLoadBalancer(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating load_balancer: %s", err))
	}
	d.SetId(d.Get("load_balancer_name").(string))
	return update(ctx, d, meta)
}
