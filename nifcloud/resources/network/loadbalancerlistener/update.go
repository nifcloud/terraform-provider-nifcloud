package loadbalancerlistener

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	if d.HasChanges(
		"balancing_type",
		"instance_port",
		"load_balancer_port",
	) && !d.IsNewResource() {
		input := expandUpdateLoadBalancer(d)
		req := svc.UpdateLoadBalancerRequest(input)
		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating load balancer %s", err))
		}

		lbID := strings.Join([]string{
			d.Get("load_balancer_name").(string),
			strconv.Itoa(d.Get("load_balancer_port").(int)),
			strconv.Itoa(d.Get("instance_port").(int)),
		}, "_")
		d.SetId(lbID)
	}
	if d.HasChanges(
		"session_stickiness_policy_enable",
		"session_stickiness_policy_expiration_period",
		"sorry_page_enable",
		"sorry_page_status_code",
	) {
		input := expandUpdateLoadBalancerOption(d)
		req := svc.UpdateLoadBalancerOptionRequest(input)
		_, err := req.Send(ctx)
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

			req := svc.RegisterInstancesWithLoadBalancerRequest(input)
			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed registering instances with load balancer: %s", err))
			}
		}

		if len(delInstances) > 0 {
			input := expandDeregisterInstancesFromLoadBalancerInput(d, delInstances)

			req := svc.DeregisterInstancesFromLoadBalancerRequest(input)

			_, err := req.Send(ctx)
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
		req := svc.ConfigureHealthCheckRequest(input)
		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating load balancer %s", err))
		}
	}
	if d.HasChange("filter_type") {
		input := expandSetFilterForLoadBalancerFilterType(d)
		req := svc.SetFilterForLoadBalancerRequest(input)
		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed setting load balancer filters %s", err))
		}
	}
	if d.HasChange("filter") {
		input := expandUnSetFilterForLoadBalancer(d)
		if len(input.IPAddresses) > 0 && *input.IPAddresses[0].IPAddress != "*.*.*.*" {
			req := svc.SetFilterForLoadBalancerRequest(input)
			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed setting load balancer filters %s", err))
			}
		}

		input = expandSetFilterForLoadBalancer(d)
		if len(input.IPAddresses) > 0 && *input.IPAddresses[0].IPAddress != "*.*.*.*" {
			req := svc.SetFilterForLoadBalancerRequest(input)
			_, err := req.Send(ctx)
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
			req := svc.UnsetLoadBalancerListenerSSLCertificateRequest(input)
			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed un setting SSLCertificate with load balancer: %s", err))
			}
		}
		input := expandSetLoadBalancerListenerSSLCertificate(d)
		req := svc.SetLoadBalancerListenerSSLCertificateRequest(input)
		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed setting SSLCertificate with load balancer: %s", err))
		}
	}
	if d.HasChanges("ssl_policy_name", "ssl_policy_id") {
		if d.Get("ssl_policy_name") == "" && d.Get("ssl_policy_id") == "" {
			input := expandNiftyUnsetLoadBalancerSSLPoliciesOfListener(d)
			req := svc.NiftyUnsetLoadBalancerSSLPoliciesOfListenerRequest(input)
			_, err := req.Send(ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating load balancer unset %s", err))
			}
		}
	}
	if d.HasChange("ssl_policy_id") && d.Get("ssl_policy_id") != "" && d.Get("ssl_policy_id") != nil {
		input := expandNiftySetLoadBalancerSSLPoliciesOfListenerForPolicyID(d)
		req := svc.NiftySetLoadBalancerSSLPoliciesOfListenerRequest(input)
		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating load balancer set ssl_policy_id %s", err))
		}
	}
	if d.HasChange("ssl_policy_name") && d.Get("ssl_policy_name") != "" && d.Get("ssl_policy_name") != nil {
		input := expandNiftySetLoadBalancerSSLPoliciesOfListenerForPolicyName(d)
		req := svc.NiftySetLoadBalancerSSLPoliciesOfListenerRequest(input)
		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating load balancer set ssl_policy_name %s", err))
		}
	}
	return read(ctx, d, meta)
}
