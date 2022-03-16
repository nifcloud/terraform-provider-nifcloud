package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	// lintignore:R019
	if d.HasChanges(
		"accounting_type",
		"network_volume",
		"balancing_type",
		"instance_port",
		"load_balancer_port",
		"load_balancer_name",
	) && !d.IsNewResource() {
		input := expandUpdateLoadBalancer(d)
		_, err := svc.UpdateLoadBalancer(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating load balancer %s", err))
		}
		d.SetId(d.Get("load_balancer_name").(string))
	}
	if d.HasChanges(
		"session_stickiness_policy_enable",
		"session_stickiness_policy_expiration_period",
		"sorry_page_enable",
		"sorry_page_status_code",
	) {
		input := expandUpdateLoadBalancerOption(d)
		_, err := svc.UpdateLoadBalancerOption(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating load balancer %s", err))
		}
	}
	if d.HasChange("instances") {
		o, n := d.GetChange("instances")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		addInstances := ns.Difference(os).List()
		delInstances := os.Difference(ns).List()

		if len(addInstances) > 0 {
			input := expandRegisterInstancesWithLoadBalancerInput(d, addInstances)

			_, err := svc.RegisterInstancesWithLoadBalancer(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed registering instances with load balancer: %s", err))
			}
		}

		if len(delInstances) > 0 {
			input := expandDeregisterInstancesFromLoadBalancerInput(d, delInstances)

			_, err := svc.DeregisterInstancesFromLoadBalancer(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed deregistering instances with elb: %s", err))
			}
		}
	}
	if d.HasChanges(
		"unhealthy_threshold",
		"health_check_target",
		"health_check_interval",
		"health_check_timeout",
	) {
		input := expandConfigureHealthCheck(d)
		_, err := svc.ConfigureHealthCheck(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating load balancer %s", err))
		}
	}
	if d.HasChange("filter_type") {
		input := expandSetFilterForLoadBalancerFilterType(d)
		_, err := svc.SetFilterForLoadBalancer(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed setting load balancer filters %s", err))
		}
	}
	if d.HasChange("filter") {
		input := expandUnSetFilterForLoadBalancer(d)
		if len(input.IPAddresses.Member) > 0 && nifcloud.ToString(input.IPAddresses.Member[0].IPAddress) != "*.*.*.*" {
			_, err := svc.SetFilterForLoadBalancer(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed setting load balancer filters %s", err))
			}
		}

		input = expandSetFilterForLoadBalancer(d)
		if len(input.IPAddresses.Member) > 0 && nifcloud.ToString(input.IPAddresses.Member[0].IPAddress) != "*.*.*.*" {
			_, err := svc.SetFilterForLoadBalancer(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed setting load balancer filters %s", err))
			}
		}
	}
	if d.HasChange("ssl_certificate_id") {
		o, _ := d.GetChange("ssl_certificate_id")
		oc := o.(string)
		if oc != "" {
			input := expandUnsetLoadBalancerListenerSSLCertificate(d)
			_, err := svc.UnsetLoadBalancerListenerSSLCertificate(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed un setting SSLCertificate with load balancer: %s", err))
			}
		}
		input := expandSetLoadBalancerListenerSSLCertificate(d)
		_, err := svc.SetLoadBalancerListenerSSLCertificate(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed setting SSLCertificate with load balancer: %s", err))
		}
	}
	if d.HasChanges("ssl_policy_name", "ssl_policy_id") {
		if d.Get("ssl_policy_name") == "" && d.Get("ssl_policy_id") == "" {
			input := expandNiftyUnsetLoadBalancerSSLPoliciesOfListener(d)
			_, err := svc.NiftyUnsetLoadBalancerSSLPoliciesOfListener(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating load balancer unset %s", err))
			}
		}
	}
	if d.HasChange("ssl_policy_id") && d.Get("ssl_policy_id") != "" && d.Get("ssl_policy_id") != nil {
		input := expandNiftySetLoadBalancerSSLPoliciesOfListenerForPolicyID(d)
		_, err := svc.NiftySetLoadBalancerSSLPoliciesOfListener(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating load balancer set ssl_policy_id %s", err))
		}
	}
	if d.HasChange("ssl_policy_name") && d.Get("ssl_policy_name") != "" && d.Get("ssl_policy_name") != nil {
		input := expandNiftySetLoadBalancerSSLPoliciesOfListenerForPolicyName(d)
		_, err := svc.NiftySetLoadBalancerSSLPoliciesOfListener(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating load balancer set ssl_policy_name %s", err))
		}
	}
	return read(ctx, d, meta)
}
