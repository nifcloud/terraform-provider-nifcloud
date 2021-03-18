package elb

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandNiftyCreateElasticLoadBalancerInput(d *schema.ResourceData) *computing.NiftyCreateElasticLoadBalancerInput {
	var networkInterface []computing.RequestNetworkInterfaceOfNiftyCreateElasticLoadBalancer
	for _, ni := range d.Get("network_interface").(*schema.Set).List() {
		if v, ok := ni.(map[string]interface{}); ok {
			n := computing.RequestNetworkInterfaceOfNiftyCreateElasticLoadBalancer{}
			if row, ok := v["network_id"]; ok {
				n.NetworkId = nifcloud.String(row.(string))
			}
			if row, ok := v["network_name"]; ok {
				n.NetworkName = nifcloud.String(row.(string))
			}
			if row, ok := v["ip_address"]; ok {
				n.IpAddress = nifcloud.String(row.(string))
			}
			if row, ok := v["is_vip_network"]; ok {
				n.IsVipNetwork = nifcloud.Bool(row.(bool))
			}
			networkInterface = append(networkInterface, n)
		}
	}

	var expectations []computing.RequestExpectationOfNiftyCreateElasticLoadBalancer
	for _, expectation := range d.Get("health_check_expectation_http_code").(*schema.Set).List() {
		expectations = append(expectations, computing.RequestExpectationOfNiftyCreateElasticLoadBalancer{
			HttpCode: nifcloud.Int64(int64(expectation.(int))),
		})
	}

	input := &computing.NiftyCreateElasticLoadBalancerInput{
		ElasticLoadBalancerName: nifcloud.String(d.Get("elb_name").(string)),
		Listeners: []computing.RequestListenersOfNiftyCreateElasticLoadBalancer{
			{
				Protocol:                computing.ProtocolOfListenersForNiftyCreateElasticLoadBalancer(d.Get("protocol").(string)),
				ElasticLoadBalancerPort: nifcloud.Int64(int64(d.Get("lb_port").(int))),
				InstancePort:            nifcloud.Int64(int64(d.Get("instance_port").(int))),
				BalancingType:           nifcloud.Int64(int64(d.Get("balancing_type").(int))),
				Description:             nifcloud.String(d.Get("description").(string)),
				RequestHealthCheck: &computing.RequestHealthCheckOfNiftyCreateElasticLoadBalancer{
					Target:             nifcloud.String(d.Get("health_check_target").(string)),
					Interval:           nifcloud.Int64(int64(d.Get("health_check_interval").(int))),
					UnhealthyThreshold: nifcloud.Int64(int64(d.Get("unhealthy_threshold").(int))),
				},
				SSLCertificateId: nifcloud.String(d.Get("ssl_certificate_id").(string)),
				RequestSession: &computing.RequestSession{
					RequestStickinessPolicy: &computing.RequestStickinessPolicy{
						Enable:           nifcloud.Bool(d.Get("session_stickiness_policy_enable").(bool)),
						Method:           computing.MethodOfListenersForNiftyCreateElasticLoadBalancer(d.Get("session_stickiness_policy_method").(string)),
						ExpirationPeriod: nifcloud.Int64(int64(d.Get("session_stickiness_policy_expiration_period").(int))),
					},
				},
				RequestSorryPage: &computing.RequestSorryPage{
					Enable:      nifcloud.Bool(d.Get("sorry_page_enable").(bool)),
					RedirectUrl: nifcloud.String(d.Get("sorry_page_redirect_url").(string)),
				},
			},
		},
		AvailabilityZones: []string{d.Get("availability_zone").(string)},
		NetworkVolume:     nifcloud.Int64(int64(d.Get("network_volume").(int))),
		AccountingType:    computing.AccountingTypeOfNiftyCreateElasticLoadBalancerRequest(d.Get("accounting_type").(string)),
		NetworkInterface:  networkInterface,
	}

	if strings.HasPrefix(nifcloud.StringValue(input.Listeners[0].RequestHealthCheck.Target), "HTTP") {
		input.Listeners[0].RequestHealthCheck.ListOfRequestExpectation = expectations
		input.Listeners[0].RequestHealthCheck.Path = nifcloud.String(d.Get("health_check_path").(string))
	}
	return input
}

func expandNiftyDescribeElasticLoadBalancersInput(d *schema.ResourceData) *computing.NiftyDescribeElasticLoadBalancersInput {
	lbPort := d.Get("lb_port")
	instancePort := d.Get("instance_port")
	protocol := d.Get("protocol")

	return &computing.NiftyDescribeElasticLoadBalancersInput{
		ElasticLoadBalancers: &computing.RequestElasticLoadBalancers{
			ListOfRequestElasticLoadBalancerId:   []string{d.Id()},
			ListOfRequestElasticLoadBalancerPort: []int64{int64(lbPort.(int))},
			ListOfRequestInstancePort:            []int64{int64(instancePort.(int))},
			ListOfRequestProtocol:                []string{protocol.(string)},
		},
	}
}

func expandNiftyDescribeElasticLoadBalancersInputWithName(d *schema.ResourceData) *computing.NiftyDescribeElasticLoadBalancersInput {
	return &computing.NiftyDescribeElasticLoadBalancersInput{
		ElasticLoadBalancers: &computing.RequestElasticLoadBalancers{
			ListOfRequestElasticLoadBalancerName: []string{d.Get("elb_name").(string)},
			ListOfRequestElasticLoadBalancerPort: []int64{int64(d.Get("lb_port").(int))},
			ListOfRequestInstancePort:            []int64{int64(d.Get("instance_port").(int))},
			ListOfRequestProtocol:                []string{d.Get("protocol").(string)},
		},
	}
}

func expandNiftyConfigureElasticLoadBalancerHealthCheckInput(d *schema.ResourceData) *computing.NiftyConfigureElasticLoadBalancerHealthCheckInput {
	var expectations []computing.RequestExpectation
	for _, expectation := range d.Get("health_check_expectation_http_code").(*schema.Set).List() {
		expectations = append(expectations, computing.RequestExpectation{
			HttpCode: nifcloud.Int64(int64(expectation.(int))),
		})
	}

	input := &computing.NiftyConfigureElasticLoadBalancerHealthCheckInput{
		ElasticLoadBalancerId:   nifcloud.String(d.Id()),
		ElasticLoadBalancerPort: nifcloud.Int64(int64(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int64(int64(d.Get("instance_port").(int))),
		Protocol:                computing.ProtocolOfNiftyConfigureElasticLoadBalancerHealthCheckRequest(d.Get("protocol").(string)),
		HealthCheck: &computing.RequestHealthCheckOfNiftyConfigureElasticLoadBalancerHealthCheck{
			Target:             nifcloud.String(d.Get("health_check_target").(string)),
			Interval:           nifcloud.Int64(int64(d.Get("health_check_interval").(int))),
			UnhealthyThreshold: nifcloud.Int64(int64(d.Get("unhealthy_threshold").(int))),
		},
	}

	if strings.HasPrefix(nifcloud.StringValue(input.HealthCheck.Target), "HTTP") {
		input.HealthCheck.ListOfRequestExpectation = expectations
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
		ElasticLoadBalancerId:   nifcloud.String(d.Id()),
		ElasticLoadBalancerPort: nifcloud.Int64(int64(lbPortBefore.(int))),
		InstancePort:            nifcloud.Int64(int64(instancePortBefore.(int))),
		Protocol:                computing.ProtocolOfNiftyModifyElasticLoadBalancerAttributesRequest(protocolBefore.(string)),
		LoadBalancerAttributes: &computing.RequestLoadBalancerAttributes{
			RequestSession: &computing.RequestSessionOfNiftyModifyElasticLoadBalancerAttributes{
				RequestStickinessPolicy: &computing.RequestStickinessPolicyOfNiftyModifyElasticLoadBalancerAttributes{
					Enable: nifcloud.Bool(d.Get("session_stickiness_policy_enable").(bool)),
					Method: computing.MethodOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributes(
						d.Get("session_stickiness_policy_method").(string),
					),
					ExpirationPeriod: nifcloud.Int64(int64(d.Get("session_stickiness_policy_expiration_period").(int))),
				},
			},
			RequestSorryPage: &computing.RequestSorryPageOfNiftyModifyElasticLoadBalancerAttributes{
				Enable:      nifcloud.Bool(d.Get("sorry_page_enable").(bool)),
				RedirectUrl: nifcloud.String(d.Get("sorry_page_redirect_url").(string)),
			},
			ListOfRequestAdditionalAttributes: []computing.RequestAdditionalAttributes{
				{
					Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesProtocol,
					Value: nifcloud.String(protocolAfter.(string)),
				},
				{
					Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesDescription,
					Value: nifcloud.String(d.Get("description").(string)),
				},
				{
					Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesElasticLoadBalancerPort,
					Value: nifcloud.String(strconv.Itoa(lbPortAfter.(int))),
				},
				{
					Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesInstancePort,
					Value: nifcloud.String(strconv.Itoa(instancePortAfter.(int))),
				},
				{
					Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesBalancingType,
					Value: nifcloud.String(strconv.Itoa(d.Get("balancing_type").(int))),
				},
			},
		},
	}

	if d.Get("protocol").(string) == "HTTPS" {
		input.LoadBalancerAttributes.ListOfRequestAdditionalAttributes = append(
			input.LoadBalancerAttributes.ListOfRequestAdditionalAttributes, computing.RequestAdditionalAttributes{
				Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesSslCertificateId,
				Value: nifcloud.String(d.Get("ssl_certificate_id").(string)),
			})
	}
	return input
}

func expandNiftyUpdateElasticLoadBalancerInput(d *schema.ResourceData) *computing.NiftyUpdateElasticLoadBalancerInput {
	accountingType, _ := strconv.Atoi(d.Get("accounting_type").(string))

	input := &computing.NiftyUpdateElasticLoadBalancerInput{
		ElasticLoadBalancerId: nifcloud.String(d.Id()),
		AccountingTypeUpdate:  nifcloud.Int64(int64(accountingType)),
		NetworkVolumeUpdate:   nifcloud.Int64(int64(d.Get("network_volume").(int))),
	}

	if d.HasChange("elb_name") && !d.IsNewResource() {
		input.ElasticLoadBalancerNameUpdate = nifcloud.String(d.Get("elb_name").(string))
	}
	return input
}

func expandNiftyRegisterInstancesWithElasticLoadBalancerInput(
	d *schema.ResourceData,
	list []interface{},
) *computing.NiftyRegisterInstancesWithElasticLoadBalancerInput {
	var instances []computing.RequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer
	for _, i := range list {
		instances = append(instances, computing.RequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer{
			InstanceId: nifcloud.String(i.(string)),
		})
	}

	input := &computing.NiftyRegisterInstancesWithElasticLoadBalancerInput{
		ElasticLoadBalancerId:   nifcloud.String(d.Id()),
		ElasticLoadBalancerPort: nifcloud.Int64(int64(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int64(int64(d.Get("instance_port").(int))),
		Protocol:                computing.ProtocolOfNiftyRegisterInstancesWithElasticLoadBalancerRequest(d.Get("protocol").(string)),
		Instances:               instances,
	}
	return input
}

func expandNiftyDeregisterInstancesFromElasticLoadBalancerInput(
	d *schema.ResourceData,
	list []interface{},
) *computing.NiftyDeregisterInstancesFromElasticLoadBalancerInput {
	var instances []computing.RequestInstancesOfNiftyDeregisterInstancesFromElasticLoadBalancer
	for _, i := range list {
		instances = append(instances, computing.RequestInstancesOfNiftyDeregisterInstancesFromElasticLoadBalancer{
			InstanceId: nifcloud.String(i.(string)),
		})
	}

	input := &computing.NiftyDeregisterInstancesFromElasticLoadBalancerInput{
		ElasticLoadBalancerId:   nifcloud.String(d.Id()),
		ElasticLoadBalancerPort: nifcloud.Int64(int64(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int64(int64(d.Get("instance_port").(int))),
		Protocol:                computing.ProtocolOfNiftyDeregisterInstancesFromElasticLoadBalancerRequest(d.Get("protocol").(string)),
		Instances:               instances,
	}
	return input
}

func expandNiftyReplaceRouteTableAssociationWithElasticLoadBalancerInput(
	d *schema.ResourceData,
) *computing.NiftyReplaceRouteTableAssociationWithElasticLoadBalancerInput {
	input := &computing.NiftyReplaceRouteTableAssociationWithElasticLoadBalancerInput{
		RouteTableId:  nifcloud.String(d.Get("route_table_id").(string)),
		AssociationId: nifcloud.String(d.Get("route_table_association_id").(string)),
	}
	return input
}

func expandNiftyDisassociateRouteTableFromElasticLoadBalancerInput(
	d *schema.ResourceData,
) *computing.NiftyDisassociateRouteTableFromElasticLoadBalancerInput {
	input := &computing.NiftyDisassociateRouteTableFromElasticLoadBalancerInput{
		AssociationId: nifcloud.String(d.Get("route_table_association_id").(string)),
	}
	return input
}

func expandNiftyAssociateRouteTableWithElasticLoadBalancerInput(
	d *schema.ResourceData,
) *computing.NiftyAssociateRouteTableWithElasticLoadBalancerInput {
	input := &computing.NiftyAssociateRouteTableWithElasticLoadBalancerInput{
		RouteTableId:          nifcloud.String(d.Get("route_table_id").(string)),
		ElasticLoadBalancerId: nifcloud.String(d.Id()),
	}
	return input
}

func expandNiftyDeleteNiftyElasticLoadBalancerInput(d *schema.ResourceData) *computing.NiftyDeleteElasticLoadBalancerInput {
	return &computing.NiftyDeleteElasticLoadBalancerInput{
		ElasticLoadBalancerId:   nifcloud.String(d.Id()),
		ElasticLoadBalancerPort: nifcloud.Int64(int64(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int64(int64(d.Get("instance_port").(int))),
		Protocol:                computing.ProtocolOfNiftyDeleteElasticLoadBalancerRequest(d.Get("protocol").(string)),
	}
}
