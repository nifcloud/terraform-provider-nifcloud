package elb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyCreateElasticLoadBalancerInput(d)

	svc := meta.(*client.Client).Computing

	for _, ni := range d.Get("network_interface").(*schema.Set).List() {
		if v, ok := ni.(map[string]interface{}); ok {
			if raw, ok := v["network_id"]; ok && len(raw.(string)) > 0 {
				key, err := mutexkv.LockPrivateLan(ctx, raw.(string), svc)
				if err != nil {
					return diag.FromErr(err)
				}
				defer mutexkv.UnlockPrivateLan(key)
			}
			if raw, ok := v["network_name"]; ok && len(raw.(string)) > 0 {
				key, err := mutexkv.LockPrivateLanByName(ctx, raw.(string), svc)
				if err != nil {
					return diag.FromErr(err)
				}
				defer mutexkv.UnlockPrivateLan(key)
			}
		}
	}

	req := svc.NiftyCreateElasticLoadBalancerRequest(input)
	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating elb: %s", err))
	}

	res, err := svc.NiftyDescribeElasticLoadBalancersRequest(expandNiftyDescribeElasticLoadBalancersInputWithName(d)).Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed describe elb: %s", err))
	}

	elbID := res.NiftyDescribeElasticLoadBalancersOutput.ElasticLoadBalancerDescriptions[0].ElasticLoadBalancerId
	d.SetId(nifcloud.StringValue(elbID))

	return update(ctx, d, meta)
}
