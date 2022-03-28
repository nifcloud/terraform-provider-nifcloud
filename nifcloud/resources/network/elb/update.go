package elb

import (
	"context"
	"fmt"
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
		err := computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInput(d), time.Until(deadline))
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

		_, err := svc.NiftyModifyElasticLoadBalancerAttributes(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating elb attributes: %s", err))
		}

		err = computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
		}
	}

	if d.HasChanges("elb_name", "accounting_type", "network_volume") {
		input := expandNiftyUpdateElasticLoadBalancerInput(d)

		_, err := svc.NiftyUpdateElasticLoadBalancer(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating elb: %s", err))
		}

		err = computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInput(d), time.Until(deadline))
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

	if d.HasChange("route_table_id") {
		o, n := d.GetChange("route_table_id")
		ors := o.(string)
		nrs := n.(string)

		if ors != "" && nrs != "" {
			input := expandNiftyReplaceRouteTableAssociationWithElasticLoadBalancerInput(d)
			_, err := svc.NiftyReplaceRouteTableAssociationWithElasticLoadBalancer(ctx, input)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed replace route table with elb: %s", err))
			}
		}
		if ors == "" && nrs != "" {
			input := expandNiftyAssociateRouteTableWithElasticLoadBalancerInput(d)
			_, err := svc.NiftyAssociateRouteTableWithElasticLoadBalancer(ctx, input)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed associating route table with elb: %s", err))
			}
		}
		if ors != "" && nrs == "" {
			input := expandNiftyDisassociateRouteTableFromElasticLoadBalancerInput(d)
			_, err := svc.NiftyDisassociateRouteTableFromElasticLoadBalancer(ctx, input)

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed disassociating route table with elb: %s", err))
			}
		}

		err := computing.NewElasticLoadBalancerAvailableWaiter(svc).Wait(ctx, expandNiftyDescribeElasticLoadBalancersInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed wait until elb available: %s", err))
		}
	}
	return read(ctx, d, meta)
}
