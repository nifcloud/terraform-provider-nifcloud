package elb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyCreateElasticLoadBalancerInput(d)

	svc := meta.(*client.Client).Computing
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
