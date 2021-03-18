package elblistener

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.IsNewResource() {
		input := expandNiftyDescribeElasticLoadBalancersInput(d)

		err := svc.WaitUntilElasticLoadBalancerAvailable(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
		}
	} else {

		mutexKV.Lock(getELBID(d))
		defer mutexKV.Unlock(getELBID(d))
	}

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

		req := svc.NiftyModifyElasticLoadBalancerAttributesRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating elb attributes: %s", err))
		}

		err = svc.WaitUntilElasticLoadBalancerAvailable(ctx, expandNiftyDescribeElasticLoadBalancersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
		}
	}

	if d.HasChanges(
		"unhealthy_threshold",
		"health_check_target",
		"health_check_interval",
		"health_check_path",
		"health_check_expectation_http_code",
	) {
		input := expandNiftyConfigureElasticLoadBalancerHealthCheckInput(d)

		req := svc.NiftyConfigureElasticLoadBalancerHealthCheckRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating elb health check: %s", err))
		}

		err = svc.WaitUntilElasticLoadBalancerAvailable(ctx, expandNiftyDescribeElasticLoadBalancersInput(d))
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

			req := svc.NiftyRegisterInstancesWithElasticLoadBalancerRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed registering instances with elb: %s", err))
			}

			err = svc.WaitUntilElasticLoadBalancerAvailable(ctx, expandNiftyDescribeElasticLoadBalancersInput(d))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
			}
		}

		if len(delInstances) > 0 {
			input := expandNiftyDeregisterInstancesFromElasticLoadBalancerInput(d, delInstances)

			req := svc.NiftyDeregisterInstancesFromElasticLoadBalancerRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed deregistering instances with elb: %s", err))
			}

			err = svc.WaitUntilElasticLoadBalancerAvailable(ctx, expandNiftyDescribeElasticLoadBalancersInput(d))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
			}
		}
	}
	return read(ctx, d, meta)
}
