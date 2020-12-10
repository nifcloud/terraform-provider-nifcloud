package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateLoadBalancerInput(d)

	svc := meta.(*client.Client).Computing
	req := svc.CreateLoadBalancerRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating load_balancer: %s", err))
	}

	d.SetId(nifcloud.StringValue(req.Input.LoadBalancerName))
	return read(ctx, d, meta)
}
