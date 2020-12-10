package loadbalancer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := computing.UpdateLoadBalancerInput{
		LoadBalancerName: nifcloud.String(d.Id()),
	}
	if d.HasChange("accounting_type") {
		ac, _ := strconv.Atoi(d.Get("accounting_type").(string))
		input.AccountingTypeUpdate = nifcloud.Int64(int64(ac))
	}
	if d.HasChange("network_volume") {
		input.NetworkVolumeUpdate = nifcloud.Int64(int64(d.Get("network_volume").(int)))
	}
	svc := meta.(*client.Client).Computing
	req := svc.UpdateLoadBalancerRequest(&input)

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed updating load balancer %s", err))
	}
	return read(ctx, d, meta)
}
