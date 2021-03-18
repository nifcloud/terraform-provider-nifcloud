package elblistener

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyDeleteNiftyElasticLoadBalancerInput(d)
	svc := meta.(*client.Client).Computing

	mutexKV.Lock(nifcloud.StringValue(input.ElasticLoadBalancerId))
	defer mutexKV.Unlock(nifcloud.StringValue(input.ElasticLoadBalancerId))

	req := svc.NiftyDeleteElasticLoadBalancerRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = svc.WaitUntilElasticLoadBalancerDeleted(ctx, expandNiftyDescribeElasticLoadBalancersInput(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted elb error: %s", err))
	}

	d.SetId("")
	return nil
}
