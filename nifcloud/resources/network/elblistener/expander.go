package elblistener

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandNiftyRegisterPortWithElasticLoadBalancerInput(d *schema.ResourceData) *computing.NiftyRegisterPortWithElasticLoadBalancerInput {
	return &computing.NiftyRegisterPortWithElasticLoadBalancerInput{
		ElasticLoadBalancerId: nifcloud.String(d.Get("elb_id").(string)),
		Listeners: &types.ListOfRequestListenersOfNiftyRegisterPortWithElasticLoadBalancer{
			Member: []types.RequestListenersOfNiftyRegisterPortWithElasticLoadBalancer{
				{
					Protocol:                types.ProtocolOfListenersForNiftyRegisterPortWithElasticLoadBalancer(d.Get("protocol").(string)),
					ElasticLoadBalancerPort: nifcloud.Int32(int32(d.Get("lb_port").(int))),
					InstancePort:            nifcloud.Int32(int32(d.Get("instance_port").(int))),
					BalancingType:           nifcloud.Int32(int32(d.Get("balancing_type").(int))),
					Description:             nifcloud.String(d.Get("description").(string)),
					SSLCertificateId:        nifcloud.String(d.Get("ssl_certificate_id").(string)),
				},
			},
		},
	}
}

func expandNiftyDescribeElasticLoadBalancersInput(d *schema.ResourceData) *computing.NiftyDescribeElasticLoadBalancersInput {
	return &computing.NiftyDescribeElasticLoadBalancersInput{
		ElasticLoadBalancers: &types.RequestElasticLoadBalancers{
			ListOfRequestElasticLoadBalancerId:   []string{getELBID(d)},
			ListOfRequestElasticLoadBalancerPort: []int32{int32(d.Get("lb_port").(int))},
			ListOfRequestInstancePort:            []int32{int32(d.Get("instance_port").(int))},
			ListOfRequestProtocol:                []string{d.Get("protocol").(string)},
		},
	}
}

func expandNiftyDescribeElasticLoadBalancersInputWithID(d *schema.ResourceData) *computing.NiftyDescribeElasticLoadBalancersInput {
	return &computing.NiftyDescribeElasticLoadBalancersInput{
		ElasticLoadBalancers: &types.RequestElasticLoadBalancers{
			ListOfRequestElasticLoadBalancerId: []string{d.Get("elb_id").(string)},
		},
	}
}

func expandNiftyConfigureElasticLoadBalancerHealthCheckInput(d *schema.ResourceData) *computing.NiftyConfigureElasticLoadBalancerHealthCheckInput {
	var expectations []types.RequestExpectation
	for _, expectation := range d.Get("health_check_expectation_http_code").(*schema.Set).List() {
		expectations = append(expectations, types.RequestExpectation{
			HttpCode: nifcloud.Int32(int32(expectation.(int))),
		})
	}

	input := &computing.NiftyConfigureElasticLoadBalancerHealthCheckInput{
		ElasticLoadBalancerId:   nifcloud.String(getELBID(d)),
		ElasticLoadBalancerPort: nifcloud.Int32(int32(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int32(int32(d.Get("instance_port").(int))),
		Protocol:                types.ProtocolOfNiftyConfigureElasticLoadBalancerHealthCheckRequest(d.Get("protocol").(string)),
		HealthCheck: &types.RequestHealthCheckOfNiftyConfigureElasticLoadBalancerHealthCheck{
			Target:             nifcloud.String(d.Get("health_check_target").(string)),
			Interval:           nifcloud.Int32(int32(d.Get("health_check_interval").(int))),
			UnhealthyThreshold: nifcloud.Int32(int32(d.Get("unhealthy_threshold").(int))),
		},
	}

	if strings.HasPrefix(nifcloud.ToString(input.HealthCheck.Target), "HTTP") {
		input.HealthCheck.ListOfRequestExpectation = &types.ListOfRequestExpectation{Member: expectations}
		input.HealthCheck.Path = nifcloud.String(d.Get("health_check_path").(string))
	}
	return input
}

func expandNiftyModifyElasticLoadBalancerAttributesInput(d *schema.ResourceData) *computing.NiftyModifyElasticLoadBalancerAttributesInput {

	lbPortBefore := d.Get("lb_port")
	lbPortAfter := d.Get("lb_port")
	if d.HasChange("lb_port") && !d.IsNewResource() {
		lbPortBefore, lbPortAfter = d.GetChange("lb_port")
	}

	instancePortBefore := d.Get("instance_port")
	instancePortAfter := d.Get("instance_port")
	if d.HasChange("instance_port") && !d.IsNewResource() {
		instancePortBefore, instancePortAfter = d.GetChange("instance_port")
	}

	protocolBefore := d.Get("protocol")
	protocolAfter := d.Get("protocol")
	if d.HasChange("protocol") && !d.IsNewResource() {
		protocolBefore, protocolAfter = d.GetChange("protocol")
	}

	input := &computing.NiftyModifyElasticLoadBalancerAttributesInput{
		ElasticLoadBalancerId:   nifcloud.String(getELBID(d)),
		ElasticLoadBalancerPort: nifcloud.Int32(int32(lbPortBefore.(int))),
		InstancePort:            nifcloud.Int32(int32(instancePortBefore.(int))),
		Protocol:                types.ProtocolOfNiftyModifyElasticLoadBalancerAttributesRequest(protocolBefore.(string)),
		LoadBalancerAttributes: &types.RequestLoadBalancerAttributes{
			RequestSession: &types.RequestSessionOfNiftyModifyElasticLoadBalancerAttributes{
				RequestStickinessPolicy: &types.RequestStickinessPolicyOfNiftyModifyElasticLoadBalancerAttributes{
					Enable: nifcloud.Bool(d.Get("session_stickiness_policy_enable").(bool)),
					Method: types.MethodOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributes(
						d.Get("session_stickiness_policy_method").(string),
					),
					ExpirationPeriod: nifcloud.Int32(int32(d.Get("session_stickiness_policy_expiration_period").(int))),
				},
			},
			RequestSorryPage: &types.RequestSorryPage{
				Enable:      nifcloud.Bool(d.Get("sorry_page_enable").(bool)),
				RedirectUrl: nifcloud.String(d.Get("sorry_page_redirect_url").(string)),
			},
			ListOfRequestAdditionalAttributes: &types.ListOfRequestAdditionalAttributes{
				Member: []types.RequestAdditionalAttributes{
					{
						Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesProtocol,
						Value: nifcloud.String(protocolAfter.(string)),
					},
					{
						Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesDescription,
						Value: nifcloud.String(d.Get("description").(string)),
					},
					{
						Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesElasticLoadBalancerPort,
						Value: nifcloud.String(strconv.Itoa(lbPortAfter.(int))),
					},
					{
						Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesInstancePort,
						Value: nifcloud.String(strconv.Itoa(instancePortAfter.(int))),
					},
					{
						Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesBalancingType,
						Value: nifcloud.String(strconv.Itoa(d.Get("balancing_type").(int))),
					},
				},
			},
		},
	}

	if d.Get("protocol").(string) == "HTTPS" {
		input.LoadBalancerAttributes.ListOfRequestAdditionalAttributes.Member = append(
			input.LoadBalancerAttributes.ListOfRequestAdditionalAttributes.Member, types.RequestAdditionalAttributes{
				Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesSslCertificateId,
				Value: nifcloud.String(d.Get("ssl_certificate_id").(string)),
			})
	}
	return input
}

func expandNiftyRegisterInstancesWithElasticLoadBalancerInput(
	d *schema.ResourceData,
	list []interface{},
) *computing.NiftyRegisterInstancesWithElasticLoadBalancerInput {
	var instances []types.RequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer
	for _, i := range list {
		instances = append(instances, types.RequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer{
			InstanceId: nifcloud.String(i.(string)),
		})
	}

	input := &computing.NiftyRegisterInstancesWithElasticLoadBalancerInput{
		ElasticLoadBalancerId:   nifcloud.String(getELBID(d)),
		ElasticLoadBalancerPort: nifcloud.Int32(int32(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int32(int32(d.Get("instance_port").(int))),
		Protocol:                types.ProtocolOfNiftyRegisterInstancesWithElasticLoadBalancerRequest(d.Get("protocol").(string)),
		Instances:               &types.ListOfRequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer{Member: instances},
	}
	return input
}

func expandNiftyDeregisterInstancesFromElasticLoadBalancerInput(
	d *schema.ResourceData,
	list []interface{},
) *computing.NiftyDeregisterInstancesFromElasticLoadBalancerInput {
	var instances []types.RequestInstancesOfNiftyDeregisterInstancesFromElasticLoadBalancer
	for _, i := range list {
		instances = append(instances, types.RequestInstancesOfNiftyDeregisterInstancesFromElasticLoadBalancer{
			InstanceId: nifcloud.String(i.(string)),
		})
	}

	input := &computing.NiftyDeregisterInstancesFromElasticLoadBalancerInput{
		ElasticLoadBalancerId:   nifcloud.String(getELBID(d)),
		ElasticLoadBalancerPort: nifcloud.Int32(int32(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int32(int32(d.Get("instance_port").(int))),
		Protocol:                types.ProtocolOfNiftyDeregisterInstancesFromElasticLoadBalancerRequest(d.Get("protocol").(string)),
		Instances:               &types.ListOfRequestInstancesOfNiftyDeregisterInstancesFromElasticLoadBalancer{Member: instances},
	}
	return input
}

func expandNiftyDeleteNiftyElasticLoadBalancerInput(d *schema.ResourceData) *computing.NiftyDeleteElasticLoadBalancerInput {
	return &computing.NiftyDeleteElasticLoadBalancerInput{
		ElasticLoadBalancerId:   nifcloud.String(getELBID(d)),
		ElasticLoadBalancerPort: nifcloud.Int32(int32(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int32(int32(d.Get("instance_port").(int))),
		Protocol:                types.ProtocolOfNiftyDeleteElasticLoadBalancerRequest(d.Get("protocol").(string)),
	}
}
