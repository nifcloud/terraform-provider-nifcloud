package elb

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandNiftyCreateElasticLoadBalancerInput(d *schema.ResourceData) *computing.NiftyCreateElasticLoadBalancerInput {
	var networkInterface []types.RequestNetworkInterfaceOfNiftyCreateElasticLoadBalancer
	for _, ni := range d.Get("network_interface").(*schema.Set).List() {
		if v, ok := ni.(map[string]interface{}); ok {
			n := types.RequestNetworkInterfaceOfNiftyCreateElasticLoadBalancer{}
			if row, ok := v["network_id"]; ok {
				n.NetworkId = nifcloud.String(row.(string))
			}
			if row, ok := v["network_name"]; ok {
				n.NetworkName = nifcloud.String(row.(string))
			}

			if row, ok := v["ip_address"]; ok {
				if row.(string) != "" {
					n.IpAddress = nifcloud.String(row.(string))
				}
			}

			if row, ok := v["is_vip_network"]; ok {
				n.IsVipNetwork = nifcloud.Bool(row.(bool))
			}
			if row, ok := v["system_ip_addresses"]; ok {
				var listOfRequestSystemIpAddresses []types.RequestSystemIpAddresses
				for _, sia := range row.(*schema.Set).List() {
					if vi, ok := sia.(map[string]interface{}); ok {
						s := types.RequestSystemIpAddresses{}
						if ipaddress, ok := vi["system_ip_address"]; ok {
							s.SystemIpAddress = nifcloud.String(ipaddress.(string))
						}
						listOfRequestSystemIpAddresses = append(listOfRequestSystemIpAddresses, s)
					}
				}
				n.ListOfRequestSystemIpAddresses = listOfRequestSystemIpAddresses
			}
			networkInterface = append(networkInterface, n)
		}
	}
	var expectations []types.RequestExpectationOfNiftyCreateElasticLoadBalancer
	for _, expectation := range d.Get("health_check_expectation_http_code").(*schema.Set).List() {
		expectations = append(expectations, types.RequestExpectationOfNiftyCreateElasticLoadBalancer{
			HttpCode: types.HttpCodeOfListenersForNiftyCreateElasticLoadBalancer(
				expectation.(string),
			),
		})
	}
	input := &computing.NiftyCreateElasticLoadBalancerInput{
		ElasticLoadBalancerName: nifcloud.String(d.Get("elb_name").(string)),
		Listeners: &types.ListOfRequestListenersOfNiftyCreateElasticLoadBalancer{
			Member: []types.RequestListenersOfNiftyCreateElasticLoadBalancer{
				{
					Protocol:                types.ProtocolOfListenersForNiftyCreateElasticLoadBalancer(d.Get("protocol").(string)),
					ElasticLoadBalancerPort: nifcloud.Int32(int32(d.Get("lb_port").(int))),
					InstancePort:            nifcloud.Int32(int32(d.Get("instance_port").(int))),
					BalancingType:           nifcloud.Int32(int32(d.Get("balancing_type").(int))),
					Description:             nifcloud.String(d.Get("description").(string)),
					RequestHealthCheck: &types.RequestHealthCheckOfNiftyCreateElasticLoadBalancer{
						Target:             nifcloud.String(d.Get("health_check_target").(string)),
						Interval:           nifcloud.Int32(int32(d.Get("health_check_interval").(int))),
						UnhealthyThreshold: nifcloud.Int32(int32(d.Get("unhealthy_threshold").(int))),
					},
					SSLCertificateId: nifcloud.String(d.Get("ssl_certificate_id").(string)),
					RequestSession: &types.RequestSession{
						RequestStickinessPolicy: &types.RequestStickinessPolicy{
							Enable:           nifcloud.Bool(d.Get("session_stickiness_policy_enable").(bool)),
							Method:           types.MethodOfListenersForNiftyCreateElasticLoadBalancer(d.Get("session_stickiness_policy_method").(string)),
							ExpirationPeriod: nifcloud.Int32(int32(d.Get("session_stickiness_policy_expiration_period").(int))),
						},
					},
					RequestSorryPage: &types.RequestSorryPage{
						Enable:      nifcloud.Bool(d.Get("sorry_page_enable").(bool)),
						RedirectUrl: nifcloud.String(d.Get("sorry_page_redirect_url").(string)),
					},
				},
			},
		},
		AvailabilityZones: &types.ListOfRequestAvailabilityZones{Member: []string{d.Get("availability_zone").(string)}},
		NetworkVolume:     nifcloud.Int32(int32(d.Get("network_volume").(int))),
		AccountingType:    types.AccountingTypeOfNiftyCreateElasticLoadBalancerRequest(d.Get("accounting_type").(string)),
		NetworkInterface:  networkInterface,
	}

	if strings.HasPrefix(nifcloud.ToString(input.Listeners.Member[0].RequestHealthCheck.Target), "HTTP") {
		input.Listeners.Member[0].RequestHealthCheck.ListOfRequestExpectation = &types.ListOfRequestExpectationOfNiftyCreateElasticLoadBalancer{Member: expectations}
		input.Listeners.Member[0].RequestHealthCheck.Path = nifcloud.String(d.Get("health_check_path").(string))
	}
	return input
}

func expandNiftyReplaceElasticLoadBalancerLatestVersionInput(d *schema.ResourceData) *computing.NiftyReplaceElasticLoadBalancerLatestVersionInput {
	var networkInterface []types.RequestNetworkInterfaceOfNiftyReplaceElasticLoadBalancerLatestVersion
	for _, ni := range d.Get("network_interface").(*schema.Set).List() {
		if v, ok := ni.(map[string]interface{}); ok {
			n := types.RequestNetworkInterfaceOfNiftyReplaceElasticLoadBalancerLatestVersion{}
			if row, ok := v["network_id"]; ok {
				n.NetworkId = nifcloud.String(row.(string))
			}
			if row, ok := v["system_ip_addresses"]; ok {
				var listOfRequestSystemIpAddresses []types.RequestSystemIpAddresses
				for _, sia := range row.(*schema.Set).List() {
					if vi, ok := sia.(map[string]interface{}); ok {
						s := types.RequestSystemIpAddresses{}
						if ipaddress, ok := vi["system_ip_address"]; ok {
							s.SystemIpAddress = nifcloud.String(ipaddress.(string))
						}
						listOfRequestSystemIpAddresses = append(listOfRequestSystemIpAddresses, s)
					}
				}
				n.ListOfRequestSystemIpAddresses = listOfRequestSystemIpAddresses
			}
			networkInterface = append(networkInterface, n)
		}
	}
	input := &computing.NiftyReplaceElasticLoadBalancerLatestVersionInput{
		ElasticLoadBalancerId:   nifcloud.String(d.Get("elb_id").(string)),
		ElasticLoadBalancerName: nifcloud.String(d.Get("elb_name").(string)),
		NetworkInterface:        networkInterface,
	}
	return input
}

func expandNiftyDescribeElasticLoadBalancersInput(d *schema.ResourceData) *computing.NiftyDescribeElasticLoadBalancersInput {
	lbPort := d.Get("lb_port")
	instancePort := d.Get("instance_port")
	protocol := d.Get("protocol")

	return &computing.NiftyDescribeElasticLoadBalancersInput{
		ElasticLoadBalancers: &types.RequestElasticLoadBalancers{
			ListOfRequestElasticLoadBalancerId:   []string{d.Id()},
			ListOfRequestElasticLoadBalancerPort: []int32{int32(lbPort.(int))},
			ListOfRequestInstancePort:            []int32{int32(instancePort.(int))},
			ListOfRequestProtocol:                []string{protocol.(string)},
		},
	}
}

func expandNiftyDescribeElasticLoadBalancersInputWithName(d *schema.ResourceData) *computing.NiftyDescribeElasticLoadBalancersInput {
	return &computing.NiftyDescribeElasticLoadBalancersInput{
		ElasticLoadBalancers: &types.RequestElasticLoadBalancers{
			ListOfRequestElasticLoadBalancerName: []string{d.Get("elb_name").(string)},
			ListOfRequestElasticLoadBalancerPort: []int32{int32(d.Get("lb_port").(int))},
			ListOfRequestInstancePort:            []int32{int32(d.Get("instance_port").(int))},
			ListOfRequestProtocol:                []string{d.Get("protocol").(string)},
		},
	}
}

func expandNiftyConfigureElasticLoadBalancerHealthCheckInput(d *schema.ResourceData) *computing.NiftyConfigureElasticLoadBalancerHealthCheckInput {
	var expectations []types.RequestExpectation
	for _, expectation := range d.Get("health_check_expectation_http_code").(*schema.Set).List() {
		expectations = append(expectations, types.RequestExpectation{
			HttpCode: nifcloud.String(string(expectation.(string))),
		})
	}

	input := &computing.NiftyConfigureElasticLoadBalancerHealthCheckInput{
		ElasticLoadBalancerId:   nifcloud.String(d.Id()),
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
		ElasticLoadBalancerId:   nifcloud.String(d.Id()),
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

func expandNiftyUpdateElasticLoadBalancerInput(d *schema.ResourceData) *computing.NiftyUpdateElasticLoadBalancerInput {
	accountingType, _ := strconv.ParseInt(d.Get("accounting_type").(string), 10, 32)

	input := &computing.NiftyUpdateElasticLoadBalancerInput{
		ElasticLoadBalancerId: nifcloud.String(d.Id()),
		AccountingTypeUpdate:  nifcloud.Int32(int32(accountingType)),
		NetworkVolumeUpdate:   nifcloud.Int32(int32(d.Get("network_volume").(int))),
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
	var instances []types.RequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer
	for _, i := range list {
		instances = append(instances, types.RequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer{
			InstanceId: nifcloud.String(i.(string)),
		})
	}

	input := &computing.NiftyRegisterInstancesWithElasticLoadBalancerInput{
		ElasticLoadBalancerId:   nifcloud.String(d.Id()),
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
		ElasticLoadBalancerId:   nifcloud.String(d.Id()),
		ElasticLoadBalancerPort: nifcloud.Int32(int32(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int32(int32(d.Get("instance_port").(int))),
		Protocol:                types.ProtocolOfNiftyDeregisterInstancesFromElasticLoadBalancerRequest(d.Get("protocol").(string)),
		Instances:               &types.ListOfRequestInstancesOfNiftyDeregisterInstancesFromElasticLoadBalancer{Member: instances},
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
		ElasticLoadBalancerPort: nifcloud.Int32(int32(d.Get("lb_port").(int))),
		InstancePort:            nifcloud.Int32(int32(d.Get("instance_port").(int))),
		Protocol:                types.ProtocolOfNiftyDeleteElasticLoadBalancerRequest(d.Get("protocol").(string)),
	}
}
