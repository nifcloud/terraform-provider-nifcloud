package elblistener

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	deadline, _ := ctx.Deadline()

	if d.IsNewResource() {
		input := expandNiftyDescribeElasticLoadBalancersInput(d)

		err := computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, input, time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
		}
	} else {
		mutexKV.Lock(getELBID(d))
		defer mutexKV.Unlock(getELBID(d))
	}

	// lintignore:R019
	if d.HasChanges(
		"description",
		"balancing_type",
		"instance_port",
		"protocol",
		"lb_port",
		"ssl_certificate_id",
		"session_stickiness_policy_enable",
		"session_stickiness_policy_method",
		"session_stickiness_policy_expiration_period",
		"sorry_page_enable",
		"sorry_page_redirect_url",
	) {
		input := expandNiftyModifyElasticLoadBalancerAttributesInput(d)

		_, err := svc.NiftyModifyElasticLoadBalancerAttributes(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating elb attributes: %s", err))
		}

		err = computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInput(d), time.Until(deadline))
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
	}

	// lintignore:R019
	if d.HasChanges(
		"unhealthy_threshold",
		"health_check_target",
		"health_check_interval",
		"health_check_path",
		"health_check_expectation_http_code",
	) {
		input := expandNiftyConfigureElasticLoadBalancerHealthCheckInput(d)

		_, err := svc.NiftyConfigureElasticLoadBalancerHealthCheck(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating elb health check: %s", err))
		}

		err = computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
		}
	}

	if d.HasChange("instances") {
		o, n := d.GetChange("instances")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		addInstances := ns.Difference(os).List()
		delInstances := os.Difference(ns).List()

		if len(addInstances) > 0 {
			input := expandNiftyRegisterInstancesWithElasticLoadBalancerInput(d, addInstances)

			_, err := svc.NiftyRegisterInstancesWithElasticLoadBalancer(ctx, input)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed registering instances with elb: %s", err))
			}

			err = computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInput(d), time.Until(deadline))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
			}
		}

		if len(delInstances) > 0 {
			input := expandNiftyDeregisterInstancesFromElasticLoadBalancerInput(d, delInstances)

			_, err := svc.NiftyDeregisterInstancesFromElasticLoadBalancer(ctx, input)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed deregistering instances with elb: %s", err))
			}

			err = computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInput(d), time.Until(deadline))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
			}
		}
	}
	return read(ctx, d, meta)
}
