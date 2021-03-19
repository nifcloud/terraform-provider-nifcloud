package elblistener

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyRegisterPortWithElasticLoadBalancerInput(d)

	svc := meta.(*client.Client).Computing
	req := svc.NiftyRegisterPortWithElasticLoadBalancerRequest(input)

	mutexKV.Lock(nifcloud.StringValue(input.ElasticLoadBalancerId))
	defer mutexKV.Unlock(nifcloud.StringValue(input.ElasticLoadBalancerId))

	err := svc.WaitUntilElasticLoadBalancerAvailable(ctx, expandNiftyDescribeElasticLoadBalancersInputWithID(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
	}

	_, err = req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating elb listener: %s", err))
	}

	elbID := strings.Join([]string{
		d.Get("elb_id").(string),
		d.Get("protocol").(string),
		strconv.Itoa(d.Get("lb_port").(int)),
		strconv.Itoa(d.Get("instance_port").(int)),
	}, "_")
	d.SetId(elbID)

	return update(ctx, d, meta)
}
