package elblistener

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyDeleteNiftyElasticLoadBalancerInput(d)
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	mutexKV.Lock(nifcloud.ToString(input.ElasticLoadBalancerId))
	defer mutexKV.Unlock(nifcloud.ToString(input.ElasticLoadBalancerId))

	_, err := svc.NiftyDeleteElasticLoadBalancer(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = computing.NewElasticLoadBalancerDeletedWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInput(d), time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted elb error: %s", err))
	}

	d.SetId("")
	return nil
}
