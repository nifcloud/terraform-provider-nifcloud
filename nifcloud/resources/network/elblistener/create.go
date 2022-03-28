package elblistener

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyRegisterPortWithElasticLoadBalancerInput(d)
	deadline, _ := ctx.Deadline()

	svc := meta.(*client.Client).Computing
	_, err := svc.NiftyRegisterPortWithElasticLoadBalancer(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating elb listener: %s", err))
	}

	mutexKV.Lock(nifcloud.ToString(input.ElasticLoadBalancerId))
	defer mutexKV.Unlock(nifcloud.ToString(input.ElasticLoadBalancerId))

	err = computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInputWithID(d), time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
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
