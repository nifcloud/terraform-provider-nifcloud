package loadbalancerlistener

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandRegisterPortWithLoadBalancerInput(d)

	svc := meta.(*client.Client).Computing

	_, err := svc.RegisterPortWithLoadBalancer(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating load_balancer: %s", err))
	}

	lbID := strings.Join([]string{
		d.Get("load_balancer_name").(string),
		strconv.Itoa(d.Get("load_balancer_port").(int)),
		strconv.Itoa(d.Get("instance_port").(int)),
	}, "_")
	d.SetId(lbID)
	return update(ctx, d, meta)
}
