package loadbalancerlistener

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandRegisterPortWithLoadBalancerInput(d *schema.ResourceData) *computing.RegisterPortWithLoadBalancerInput {
	return &computing.RegisterPortWithLoadBalancerInput{
		LoadBalancerName: nifcloud.String(d.Get("load_balancer_name").(string)),
		Listeners: &types.ListOfRequestListenersOfRegisterPortWithLoadBalancer{
			Member: []types.RequestListenersOfRegisterPortWithLoadBalancer{{
				BalancingType:    nifcloud.Int32(int32(d.Get("balancing_type").(int))),
				InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
				LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
			}}},
	}
}

func expandDescribeLoadBalancersInput(d *schema.ResourceData) *computing.DescribeLoadBalancersInput {
	return &computing.DescribeLoadBalancersInput{
		LoadBalancerNames: &types.ListOfRequestLoadBalancerNames{
			Member: []types.RequestLoadBalancerNames{
				{
					InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
					LoadBalancerName: nifcloud.String(getLBID(d)),
					LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
				},
			},
		},
	}
}

func expandUpdateLoadBalancer(d *schema.ResourceData) *computing.UpdateLoadBalancerInput {
	input := computing.UpdateLoadBalancerInput{
		LoadBalancerName: nifcloud.String(d.Get("load_balancer_name").(string)),
	}

	lu := &types.RequestListenerUpdate{
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		RequestListener:  &types.RequestListener{},
	}

	if d.HasChange("balancing_type") && !d.IsNewResource() {
		lu.RequestListener.BalancingType = nifcloud.Int32(int32(d.Get("balancing_type").(int)))
		input.ListenerUpdate = lu
	}

	if d.HasChange("load_balancer_port") && !d.IsNewResource() {
		lbPortBefore, lbPortAfter := d.GetChange("load_balancer_port")
		lu.LoadBalancerPort = nifcloud.Int32(int32(lbPortBefore.(int)))
		lu.RequestListener.LoadBalancerPort = nifcloud.Int32(int32(lbPortAfter.(int)))
		input.ListenerUpdate = lu
	}

	if d.HasChange("instance_port") && !d.IsNewResource() {
		instancePortBefore, instancePortAfter := d.GetChange("instance_port")
		lu.InstancePort = nifcloud.Int32(int32(instancePortBefore.(int)))
		lu.RequestListener.InstancePort = nifcloud.Int32(int32(instancePortAfter.(int)))
		input.ListenerUpdate = lu
	}
	return &input
}

func expandUpdateLoadBalancerOption(d *schema.ResourceData) *computing.UpdateLoadBalancerOptionInput {
	input := computing.UpdateLoadBalancerOptionInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
	}
	if d.HasChanges(
		"session_stickiness_policy_enable",
		"session_stickiness_policy_expiration_period",
	) {
		sspu := types.RequestSessionStickinessPolicyUpdate{}
		if d.Get("session_stickiness_policy_enable").(bool) {
			sspu.Enable = nifcloud.Bool(true)
			sspu.ExpirationPeriod = nifcloud.Int32(int32(d.Get("session_stickiness_policy_expiration_period").(int)))
		} else {
			sspu.Enable = nifcloud.Bool(false)
		}
		input.SessionStickinessPolicyUpdate = &sspu
	}
	if d.HasChanges(
		"sorry_page_enable",
		"sorry_page_status_code",
	) {
		spu := types.RequestSorryPageUpdate{}
		if d.Get("sorry_page_enable").(bool) {
			spu.Enable = nifcloud.Bool(true)
			spu.StatusCode = nifcloud.Int32(int32(d.Get("sorry_page_status_code").(int)))
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
	var instances []types.RequestInstances
	for _, i := range list {
		instances = append(instances, types.RequestInstances{
			InstanceId: nifcloud.String(i.(string)),
		})
	}

	return &computing.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		Instances:        &types.ListOfRequestInstances{Member: instances},
	}
}

func expandDeregisterInstancesFromLoadBalancerInput(
	d *schema.ResourceData,
	list []interface{},
) *computing.DeregisterInstancesFromLoadBalancerInput {
	var instances []types.RequestInstances
	for _, i := range list {
		instances = append(instances, types.RequestInstances{
			InstanceId: nifcloud.String(i.(string)),
		})
	}

	return &computing.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		Instances:        &types.ListOfRequestInstances{Member: instances},
	}
}

func expandSetLoadBalancerListenerSSLCertificate(d *schema.ResourceData) *computing.SetLoadBalancerListenerSSLCertificateInput {
	return &computing.SetLoadBalancerListenerSSLCertificateInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		SSLCertificateId: nifcloud.String(d.Get("ssl_certificate_id").(string)),
	}
}

func expandNiftySetLoadBalancerSSLPoliciesOfListenerForPolicyID(d *schema.ResourceData) *computing.NiftySetLoadBalancerSSLPoliciesOfListenerInput {
	return &computing.NiftySetLoadBalancerSSLPoliciesOfListenerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		SSLPolicyId:      nifcloud.String(d.Get("ssl_policy_id").(string)),
	}
}

func expandNiftySetLoadBalancerSSLPoliciesOfListenerForPolicyName(d *schema.ResourceData) *computing.NiftySetLoadBalancerSSLPoliciesOfListenerInput {
	return &computing.NiftySetLoadBalancerSSLPoliciesOfListenerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		SSLPolicyId:      nifcloud.String(d.Get("ssl_policy_id").(string)),
	}
}

func expandNiftyUnsetLoadBalancerSSLPoliciesOfListener(d *schema.ResourceData) *computing.NiftyUnsetLoadBalancerSSLPoliciesOfListenerInput {
	return &computing.NiftyUnsetLoadBalancerSSLPoliciesOfListenerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
	}
}

func expandUnsetLoadBalancerListenerSSLCertificate(d *schema.ResourceData) *computing.UnsetLoadBalancerListenerSSLCertificateInput {
	return &computing.UnsetLoadBalancerListenerSSLCertificateInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
	}
}

func expandConfigureHealthCheck(d *schema.ResourceData) *computing.ConfigureHealthCheckInput {
	input := types.RequestHealthCheck{
		Interval:           nifcloud.Int32(int32(d.Get("health_check_interval").(int))),
		Target:             nifcloud.String(d.Get("health_check_target").(string)),
		UnhealthyThreshold: nifcloud.Int32(int32(d.Get("unhealthy_threshold").(int))),
	}
	return &computing.ConfigureHealthCheckInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		HealthCheck:      &input,
	}
}

func expandSetFilterForLoadBalancerFilterType(d *schema.ResourceData) *computing.SetFilterForLoadBalancerInput {
	return &computing.SetFilterForLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		FilterType:       types.FilterTypeOfSetFilterForLoadBalancerRequest(d.Get("filter_type").(string)),
	}
}

func expandSetFilterForLoadBalancer(d *schema.ResourceData) *computing.SetFilterForLoadBalancerInput {
	var filters []types.RequestIPAddresses
	fl := d.Get("filter").(*schema.Set).List()
	for _, i := range fl {
		filters = append(filters, types.RequestIPAddresses{
			IPAddress:   nifcloud.String(i.(string)),
			AddOnFilter: nifcloud.Bool(true),
		})
	}
	return &computing.SetFilterForLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		IPAddresses:      &types.ListOfRequestIPAddresses{Member: filters},
		FilterType:       types.FilterTypeOfSetFilterForLoadBalancerRequest(d.Get("filter_type").(string)),
	}
}

func expandUnSetFilterForLoadBalancer(d *schema.ResourceData) *computing.SetFilterForLoadBalancerInput {
	o, _ := d.GetChange("filter")
	var filters []types.RequestIPAddresses
	fl := o.(*schema.Set).List()
	for _, i := range fl {
		filters = append(filters, types.RequestIPAddresses{
			IPAddress:   nifcloud.String(i.(string)),
			AddOnFilter: nifcloud.Bool(false),
		})
	}
	return &computing.SetFilterForLoadBalancerInput{
		LoadBalancerName: nifcloud.String(getLBID(d)),
		LoadBalancerPort: nifcloud.Int32(int32(d.Get("load_balancer_port").(int))),
		InstancePort:     nifcloud.Int32(int32(d.Get("instance_port").(int))),
		IPAddresses:      &types.ListOfRequestIPAddresses{Member: filters},
		FilterType:       types.FilterTypeOfSetFilterForLoadBalancerRequest(d.Get("filter_type").(string)),
	}
}
