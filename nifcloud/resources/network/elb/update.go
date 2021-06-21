package elb

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
		err := svc.WaitUntilElasticLoadBalancerAvailable(ctx, expandNiftyDescribeElasticLoadBalancersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
		}
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

	if d.HasChanges("elb_name", "accounting_type", "network_volume") {
		input := expandNiftyUpdateElasticLoadBalancerInput(d)

		req := svc.NiftyUpdateElasticLoadBalancerRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating elb: %s", err))
		}

		err = svc.WaitUntilElasticLoadBalancerAvailable(ctx, expandNiftyDescribeElasticLoadBalancersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
		}
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

	if d.HasChange("route_table_id") {
		o, n := d.GetChange("route_table_id")
		ors := o.(string)
		nrs := n.(string)

		if ors != "" && nrs != "" {
			input := expandNiftyReplaceRouteTableAssociationWithElasticLoadBalancerInput(d)
			req := svc.NiftyReplaceRouteTableAssociationWithElasticLoadBalancerRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed replace route table with elb: %s", err))
			}
		}
		if ors == "" && nrs != "" {
			input := expandNiftyAssociateRouteTableWithElasticLoadBalancerInput(d)
			req := svc.NiftyAssociateRouteTableWithElasticLoadBalancerRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed associating route table with elb: %s", err))
			}
		}
		if ors != "" && nrs == "" {
			input := expandNiftyDisassociateRouteTableFromElasticLoadBalancerInput(d)
			req := svc.NiftyDisassociateRouteTableFromElasticLoadBalancerRequest(input)

			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed disassociating route table with elb: %s", err))
			}
		}

		err := svc.WaitUntilElasticLoadBalancerAvailable(ctx, expandNiftyDescribeElasticLoadBalancersInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
		}
	}
	return read(ctx, d, meta)
}
