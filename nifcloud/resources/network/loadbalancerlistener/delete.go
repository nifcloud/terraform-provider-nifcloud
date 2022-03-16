package loadbalancerlistener

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	ip := d.Get("instance_port").(int)
	lbp := d.Get("load_balancer_port").(int)

	_, err := svc.DeleteLoadBalancer(ctx, &computing.DeleteLoadBalancerInput{
		InstancePort:     nifcloud.Int32(int32(ip)),
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(lbp)),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting load_balancer: %s", err))
	}
	d.SetId("")
	return nil
}
