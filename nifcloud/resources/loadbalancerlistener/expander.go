package loadbalancerlistener

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandRegisterPortWithLoadBalancerInput(d *schema.ResourceData) *computing.RegisterPortWithLoadBalancerInput {
	return &computing.RegisterPortWithLoadBalancerInput{
		LoadBalancerName: nifcloud.String(d.Get("load_balancer_name").(string)),
		Listeners: []computing.RequestListenersOfRegisterPortWithLoadBalancer{{
			BalancingType:    nifcloud.Int64(int64(d.Get("balancing_type").(int))),
			InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
			LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		}},
	}
}

func expandDescribeLoadBalancersInput(d *schema.ResourceData) *computing.DescribeLoadBalancersInput {
	return &computing.DescribeLoadBalancersInput{
		LoadBalancerNames: []computing.RequestLoadBalancerNames{
			{
				InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
				LoadBalancerName: nifcloud.String(getLBID(d)),
				LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
			},
		},
	}
}

func expandUpdateLoadBalancer(d *schema.ResourceData) *computing.UpdateLoadBalancerInput {
	input := computing.UpdateLoadBalancerInput{
		LoadBalancerName: nifcloud.String(d.Get("load_balancer_name").(string)),
	}

	lu := &computing.RequestListenerUpdate{
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		RequestListener:  &computing.RequestListener{},
	}

	if d.HasChange("balancing_type") && !d.IsNewResource() {
		lu.RequestListener.BalancingType = nifcloud.Int64(int64(d.Get("balancing_type").(int)))
		input.ListenerUpdate = lu
	}

	if d.HasChange("load_balancer_port") && !d.IsNewResource() {
		lbPortBefore, lbPortAfter := d.GetChange("load_balancer_port")
		lu.LoadBalancerPort = nifcloud.Int64(int64(lbPortBefore.(int)))
		lu.RequestListener.LoadBalancerPort = nifcloud.Int64(int64(lbPortAfter.(int)))
		input.ListenerUpdate = lu
	}

	if d.HasChange("instance_port") && !d.IsNewResource() {
		instancePortBefore, instancePortAfter := d.GetChange("instance_port")
		lu.InstancePort = nifcloud.Int64(int64(instancePortBefore.(int)))
		lu.RequestListener.InstancePort = nifcloud.Int64(int64(instancePortAfter.(int)))
		input.ListenerUpdate = lu
	}
	return &input
}

func expandUpdateLoadBalancerOption(d *schema.ResourceData) *computing.UpdateLoadBalancerOptionInput {
	input := computing.UpdateLoadBalancerOptionInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
	}
	if d.HasChanges(
		"session_stickiness_policy_enable",
		"session_stickiness_policy_expiration_period",
	) {
		sspu := computing.RequestSessionStickinessPolicyUpdate{}
		if d.Get("session_stickiness_policy_enable").(bool) {
			sspu.Enable = nifcloud.Bool(true)
			sspu.ExpirationPeriod = nifcloud.Int64(int64(d.Get("session_stickiness_policy_expiration_period").(int)))
		} else {
			sspu.Enable = nifcloud.Bool(false)
		}
		input.SessionStickinessPolicyUpdate = &sspu
	}
	if d.HasChanges(
		"sorry_page_enable",
		"sorry_page_status_code",
	) {
		spu := computing.RequestSorryPageUpdate{}
		if d.Get("sorry_page_enable").(bool) {
			spu.Enable = nifcloud.Bool(true)
			spu.StatusCode = nifcloud.Int64(int64(d.Get("sorry_page_status_code").(int)))
		} else {
			spu.Enable = nifcloud.Bool(false)
		}
		input.SorryPageUpdate = &spu
	}
	return &input
}

func expandRegisterInstancesWithLoadBalancerInput(
	d *schema.ResourceData,
	list []interface{},
) *computing.RegisterInstancesWithLoadBalancerInput {
	var instances []computing.RequestInstancesOfRegisterInstancesWithLoadBalancer
	for _, i := range list {
		instances = append(instances, computing.RequestInstancesOfRegisterInstancesWithLoadBalancer{
			InstanceId: nifcloud.String(i.(string)),
		})
	}

	return &computing.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		Instances:        instances,
	}
}

func expandDeregisterInstancesFromLoadBalancerInput(
	d *schema.ResourceData,
	list []interface{},
) *computing.DeregisterInstancesFromLoadBalancerInput {
	var instances []computing.RequestInstances
	for _, i := range list {
		instances = append(instances, computing.RequestInstances{
			InstanceId: nifcloud.String(i.(string)),
		})
	}

	return &computing.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		Instances:        instances,
	}
}

func expandSetLoadBalancerListenerSSLCertificate(d *schema.ResourceData) *computing.SetLoadBalancerListenerSSLCertificateInput {
	return &computing.SetLoadBalancerListenerSSLCertificateInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		SSLCertificateId: nifcloud.String(d.Get("ssl_certificate_id").(string)),
	}
}

func expandNiftySetLoadBalancerSSLPoliciesOfListenerForPolicyID(d *schema.ResourceData) *computing.NiftySetLoadBalancerSSLPoliciesOfListenerInput {
	return &computing.NiftySetLoadBalancerSSLPoliciesOfListenerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		SSLPolicyId:      nifcloud.String(d.Get("ssl_policy_id").(string)),
	}
}

func expandNiftySetLoadBalancerSSLPoliciesOfListenerForPolicyName(d *schema.ResourceData) *computing.NiftySetLoadBalancerSSLPoliciesOfListenerInput {
	return &computing.NiftySetLoadBalancerSSLPoliciesOfListenerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		SSLPolicyId:      nifcloud.String(d.Get("ssl_policy_id").(string)),
	}
}

func expandNiftyUnsetLoadBalancerSSLPoliciesOfListener(d *schema.ResourceData) *computing.NiftyUnsetLoadBalancerSSLPoliciesOfListenerInput {
	return &computing.NiftyUnsetLoadBalancerSSLPoliciesOfListenerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
	}
}

func expandUnsetLoadBalancerListenerSSLCertificate(d *schema.ResourceData) *computing.UnsetLoadBalancerListenerSSLCertificateInput {
	return &computing.UnsetLoadBalancerListenerSSLCertificateInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
	}
}

func expandConfigureHealthCheck(d *schema.ResourceData) *computing.ConfigureHealthCheckInput {
	input := computing.RequestHealthCheck{
		Interval:           nifcloud.Int64(int64(d.Get("health_check_interval").(int))),
		Target:             nifcloud.String(d.Get("health_check_target").(string)),
		UnhealthyThreshold: nifcloud.Int64(int64(d.Get("unhealthy_threshold").(int))),
	}
	return &computing.ConfigureHealthCheckInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		HealthCheck:      &input,
	}
}

func expandSetFilterForLoadBalancerFilterType(d *schema.ResourceData) *computing.SetFilterForLoadBalancerInput {
	return &computing.SetFilterForLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		FilterType:       computing.FilterTypeOfSetFilterForLoadBalancerRequest(d.Get("filter_type").(string)),
	}
}

func expandSetFilterForLoadBalancer(d *schema.ResourceData) *computing.SetFilterForLoadBalancerInput {
	var filters []computing.RequestIPAddresses
	fl := d.Get("filter").(*schema.Set).List()
	for _, i := range fl {
		filters = append(filters, computing.RequestIPAddresses{
			IPAddress:   nifcloud.String(i.(string)),
			AddOnFilter: nifcloud.Bool(true),
		})
	}
	return &computing.SetFilterForLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		IPAddresses:      filters,
		FilterType:       computing.FilterTypeOfSetFilterForLoadBalancerRequest(d.Get("filter_type").(string)),
	}
}

func expandUnSetFilterForLoadBalancer(d *schema.ResourceData) *computing.SetFilterForLoadBalancerInput {
	o, _ := d.GetChange("filter")
	var filters []computing.RequestIPAddresses
	fl := o.(*schema.Set).List()
	for _, i := range fl {
		filters = append(filters, computing.RequestIPAddresses{
			IPAddress:   nifcloud.String(i.(string)),
			AddOnFilter: nifcloud.Bool(false),
		})
	}
	return &computing.SetFilterForLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int64(int64(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int64(int64(d.Get("instance_port").(int))),
		IPAddresses:      filters,
		FilterType:       computing.FilterTypeOfSetFilterForLoadBalancerRequest(d.Get("filter_type").(string)),
	}
}
