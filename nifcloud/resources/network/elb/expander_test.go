package elb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandNiftyCreateElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_name":                           "test_elb_name",
		"availability_zone":                  "test_availability_zone",
		"accounting_type":                    "test_accounting_type",
		"network_volume":                     1,
		"description":                        "test_description",
		"balancing_type":                     1,
		"instance_port":                      1,
		"protocol":                           "test_protocol",
		"lb_port":                            1,
		"ssl_certificate_id":                 "test_ssl_certificate_id",
		"unhealthy_threshold":                1,
		"health_check_target":                "test_health_check_target",
		"health_check_interval":              1,
		"health_check_path":                  "test_health_check_path",
		"health_check_expectation_http_code": []interface{}{1},
		"instances":                          []interface{}{"test_instances"},
		"session_stickiness_policy_enable":   true,
		"session_stickiness_policy_method":   1,
		"session_stickiness_policy_expiration_period": 1,
		"sorry_page_enable":                           true,
		"sorry_page_redirect_url":                     "test_sorry_page_redirect_url",
		"route_table_id":                              "test_route_table_id",
		"route_table_association_id":                  "test_route_table_association_id",
		"dns_name":                                    "test_dns_name",
		"version":                                     "test_version",
		"elb_id":                                      "test_elb_id",
		"network_interface": []interface{}{map[string]interface{}{
			"network_id":     "test_network_id",
			"network_name":   "test_network_name",
			"ip_address":     "test_ip_address",
			"is_vip_network": true,
		}},
	})
	rd.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreateElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreateElasticLoadBalancerInput{
				ElasticLoadBalancerName: nifcloud.String("test_elb_name"),
				Listeners: []computing.RequestListenersOfNiftyCreateElasticLoadBalancer{
					{
						Protocol:                computing.ProtocolOfListenersForNiftyCreateElasticLoadBalancer("test_protocol"),
						ElasticLoadBalancerPort: nifcloud.Int64(1),
						InstancePort:            nifcloud.Int64(1),
						BalancingType:           nifcloud.Int64(1),
						Description:             nifcloud.String("test_description"),
						RequestHealthCheck: &computing.RequestHealthCheckOfNiftyCreateElasticLoadBalancer{
							Target:             nifcloud.String("test_health_check_target"),
							Interval:           nifcloud.Int64(1),
							UnhealthyThreshold: nifcloud.Int64(1),
						},
						SSLCertificateId: nifcloud.String("test_ssl_certificate_id"),
						RequestSession: &computing.RequestSession{
							RequestStickinessPolicy: &computing.RequestStickinessPolicy{
								Enable:           nifcloud.Bool(true),
								Method:           computing.MethodOfListenersForNiftyCreateElasticLoadBalancer("1"),
								ExpirationPeriod: nifcloud.Int64(1),
							},
						},
						RequestSorryPage: &computing.RequestSorryPage{
							Enable:      nifcloud.Bool(true),
							RedirectUrl: nifcloud.String("test_sorry_page_redirect_url"),
						},
					},
				},
				AvailabilityZones: []string{"test_availability_zone"},
				NetworkVolume:     nifcloud.Int64(1),
				AccountingType:    computing.AccountingTypeOfNiftyCreateElasticLoadBalancerRequest("test_accounting_type"),
				NetworkInterface: []computing.RequestNetworkInterfaceOfNiftyCreateElasticLoadBalancer{
					{
						NetworkId:    nifcloud.String("test_network_id"),
						NetworkName:  nifcloud.String("test_network_name"),
						IpAddress:    nifcloud.String("test_ip_address"),
						IsVipNetwork: nifcloud.Bool(true),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreateElasticLoadBalancerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDescribeElasticLoadBalancersInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_port": 1,
		"protocol":      "test_protocol",
		"lb_port":       1,
	})
	rd.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDescribeElasticLoadBalancersInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDescribeElasticLoadBalancersInput{
				ElasticLoadBalancers: &computing.RequestElasticLoadBalancers{
					ListOfRequestElasticLoadBalancerId:   []string{"test_elb_id"},
					ListOfRequestElasticLoadBalancerPort: []int64{1},
					ListOfRequestInstancePort:            []int64{1},
					ListOfRequestProtocol:                []string{"test_protocol"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDescribeElasticLoadBalancersInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDescribeElasticLoadBalancersInputWithName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_port": 1,
		"protocol":      "test_protocol",
		"elb_name":      "test_elb_name",
		"lb_port":       1,
	})
	rd.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDescribeElasticLoadBalancersInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDescribeElasticLoadBalancersInput{
				ElasticLoadBalancers: &computing.RequestElasticLoadBalancers{
					ListOfRequestElasticLoadBalancerName: []string{"test_elb_name"},
					ListOfRequestElasticLoadBalancerPort: []int64{1},
					ListOfRequestInstancePort:            []int64{1},
					ListOfRequestProtocol:                []string{"test_protocol"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDescribeElasticLoadBalancersInputWithName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyConfigureElasticLoadBalancerHealthCheckInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_port":         1,
		"protocol":              "test_protocol",
		"elb_name":              "test_elb_name",
		"lb_port":               1,
		"unhealthy_threshold":   1,
		"health_check_target":   "test_health_check_target",
		"health_check_interval": 1,
	})
	rd.SetId("test_elb_id")

	rdWithHTTP := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_port":                      1,
		"protocol":                           "HTTP",
		"elb_name":                           "test_elb_name",
		"lb_port":                            1,
		"unhealthy_threshold":                1,
		"health_check_target":                "HTTPS:443",
		"health_check_interval":              1,
		"health_check_path":                  "test_health_check_path",
		"health_check_expectation_http_code": []interface{}{1},
	})
	rdWithHTTP.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyConfigureElasticLoadBalancerHealthCheckInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyConfigureElasticLoadBalancerHealthCheckInput{
				ElasticLoadBalancerId:   nifcloud.String("test_elb_id"),
				ElasticLoadBalancerPort: nifcloud.Int64(1),
				InstancePort:            nifcloud.Int64(1),
				Protocol:                computing.ProtocolOfNiftyConfigureElasticLoadBalancerHealthCheckRequest("test_protocol"),
				HealthCheck: &computing.RequestHealthCheckOfNiftyConfigureElasticLoadBalancerHealthCheck{
					Target:             nifcloud.String("test_health_check_target"),
					Interval:           nifcloud.Int64(1),
					UnhealthyThreshold: nifcloud.Int64(1),
				},
			},
		},
		{
			name: "expands the resource data with http protocol",
			args: rdWithHTTP,
			want: &computing.NiftyConfigureElasticLoadBalancerHealthCheckInput{
				ElasticLoadBalancerId:   nifcloud.String("test_elb_id"),
				ElasticLoadBalancerPort: nifcloud.Int64(1),
				InstancePort:            nifcloud.Int64(1),
				Protocol:                computing.ProtocolOfNiftyConfigureElasticLoadBalancerHealthCheckRequest("HTTP"),
				HealthCheck: &computing.RequestHealthCheckOfNiftyConfigureElasticLoadBalancerHealthCheck{
					Target:             nifcloud.String("HTTPS:443"),
					Interval:           nifcloud.Int64(1),
					UnhealthyThreshold: nifcloud.Int64(1),
					Path:               nifcloud.String("test_health_check_path"),
					ListOfRequestExpectation: []computing.RequestExpectation{
						{
							HttpCode: nifcloud.Int64(1),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyConfigureElasticLoadBalancerHealthCheckInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyElasticLoadBalancerAttributesInput(t *testing.T) {
	r := New()
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_name":                           "test_elb_name",
		"availability_zone":                  "test_availability_zone",
		"accounting_type":                    "test_accounting_type",
		"network_volume":                     1,
		"description":                        "test_description",
		"balancing_type":                     1,
		"instance_port":                      1,
		"protocol":                           "test_protocol",
		"lb_port":                            1,
		"unhealthy_threshold":                1,
		"health_check_target":                "test_health_check_target",
		"health_check_interval":              1,
		"health_check_path":                  "test_health_check_path",
		"health_check_expectation_http_code": []interface{}{1},
		"instances":                          []interface{}{"test_instances"},
		"session_stickiness_policy_enable":   true,
		"session_stickiness_policy_method":   1,
		"session_stickiness_policy_expiration_period": 1,
		"sorry_page_enable":                           true,
		"sorry_page_redirect_url":                     "test_sorry_page_redirect_url",
		"route_table_id":                              "test_route_table_id",
		"route_table_association_id":                  "test_route_table_association_id",
		"dns_name":                                    "test_dns_name",
		"version":                                     "test_version",
		"elb_id":                                      "test_elb_id",
	})
	rd.SetId("test_elb_id")
	dn := r.Data(rd.State())

	rdWithHTTPS := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_name":                           "test_elb_name",
		"availability_zone":                  "test_availability_zone",
		"accounting_type":                    "test_accounting_type",
		"network_volume":                     1,
		"description":                        "test_description",
		"balancing_type":                     1,
		"instance_port":                      1,
		"protocol":                           "HTTPS",
		"lb_port":                            1,
		"ssl_certificate_id":                 "test_ssl_certificate_id",
		"unhealthy_threshold":                1,
		"health_check_target":                "test_health_check_target",
		"health_check_interval":              1,
		"health_check_path":                  "test_health_check_path",
		"health_check_expectation_http_code": []interface{}{1},
		"instances":                          []interface{}{"test_instances"},
		"session_stickiness_policy_enable":   true,
		"session_stickiness_policy_method":   1,
		"session_stickiness_policy_expiration_period": 1,
		"sorry_page_enable":                           true,
		"sorry_page_redirect_url":                     "test_sorry_page_redirect_url",
		"route_table_id":                              "test_route_table_id",
		"route_table_association_id":                  "test_route_table_association_id",
		"dns_name":                                    "test_dns_name",
		"version":                                     "test_version",
		"elb_id":                                      "test_elb_id",
	})
	rdWithHTTPS.SetId("test_elb_id")
	dnWithHTTPS := r.Data(rdWithHTTPS.State())

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyElasticLoadBalancerAttributesInput
	}{
		{
			name: "expands the resource data",
			args: dn,
			want: &computing.NiftyModifyElasticLoadBalancerAttributesInput{
				ElasticLoadBalancerId:   nifcloud.String("test_elb_id"),
				ElasticLoadBalancerPort: nifcloud.Int64(1),
				InstancePort:            nifcloud.Int64(1),
				Protocol:                computing.ProtocolOfNiftyModifyElasticLoadBalancerAttributesRequest("test_protocol"),
				LoadBalancerAttributes: &computing.RequestLoadBalancerAttributes{
					RequestSession: &computing.RequestSessionOfNiftyModifyElasticLoadBalancerAttributes{
						RequestStickinessPolicy: &computing.RequestStickinessPolicyOfNiftyModifyElasticLoadBalancerAttributes{
							Enable: nifcloud.Bool(true),
							Method: computing.MethodOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributes(
								"1",
							),
							ExpirationPeriod: nifcloud.Int64(1),
						},
					},
					RequestSorryPage: &computing.RequestSorryPageOfNiftyModifyElasticLoadBalancerAttributes{
						Enable:      nifcloud.Bool(true),
						RedirectUrl: nifcloud.String("test_sorry_page_redirect_url"),
					},
					ListOfRequestAdditionalAttributes: []computing.RequestAdditionalAttributes{
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesProtocol,
							Value: nifcloud.String("test_protocol"),
						},
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesDescription,
							Value: nifcloud.String("test_description"),
						},
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesElasticLoadBalancerPort,
							Value: nifcloud.String("1"),
						},
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesInstancePort,
							Value: nifcloud.String("1"),
						},
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesBalancingType,
							Value: nifcloud.String("1"),
						},
					},
				},
			},
		},
		{
			name: "expands the resource data with https protocol",
			args: dnWithHTTPS,
			want: &computing.NiftyModifyElasticLoadBalancerAttributesInput{
				ElasticLoadBalancerId:   nifcloud.String("test_elb_id"),
				ElasticLoadBalancerPort: nifcloud.Int64(1),
				InstancePort:            nifcloud.Int64(1),
				Protocol:                computing.ProtocolOfNiftyModifyElasticLoadBalancerAttributesRequest("HTTPS"),
				LoadBalancerAttributes: &computing.RequestLoadBalancerAttributes{
					RequestSession: &computing.RequestSessionOfNiftyModifyElasticLoadBalancerAttributes{
						RequestStickinessPolicy: &computing.RequestStickinessPolicyOfNiftyModifyElasticLoadBalancerAttributes{
							Enable: nifcloud.Bool(true),
							Method: computing.MethodOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributes(
								"1",
							),
							ExpirationPeriod: nifcloud.Int64(1),
						},
					},
					RequestSorryPage: &computing.RequestSorryPageOfNiftyModifyElasticLoadBalancerAttributes{
						Enable:      nifcloud.Bool(true),
						RedirectUrl: nifcloud.String("test_sorry_page_redirect_url"),
					},
					ListOfRequestAdditionalAttributes: []computing.RequestAdditionalAttributes{
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesProtocol,
							Value: nifcloud.String("HTTPS"),
						},
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesDescription,
							Value: nifcloud.String("test_description"),
						},
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesElasticLoadBalancerPort,
							Value: nifcloud.String("1"),
						},
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesInstancePort,
							Value: nifcloud.String("1"),
						},
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesBalancingType,
							Value: nifcloud.String("1"),
						},
						{
							Key:   computing.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesSslCertificateId,
							Value: nifcloud.String("test_ssl_certificate_id"),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyElasticLoadBalancerAttributesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyUpdateElasticLoadBalancerInput(t *testing.T) {
	r := New()
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_name":        "test_elb_name",
		"accounting_type": "1",
		"network_volume":  1,
		"elb_id":          "test_elb_id",
	})
	rd.SetId("test_elb_id")
	dn := r.Data(rd.State())

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyUpdateElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: dn,
			want: &computing.NiftyUpdateElasticLoadBalancerInput{
				ElasticLoadBalancerId: nifcloud.String("test_elb_id"),
				AccountingTypeUpdate:  nifcloud.Int64(1),
				NetworkVolumeUpdate:   nifcloud.Int64(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyUpdateElasticLoadBalancerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeregisterInstancesFromElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":        "test_elb_id",
		"instance_port": 1,
		"protocol":      "test_protocol",
		"lb_port":       1,
		"instances":     []interface{}{"test_instances"},
	})
	rd.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeregisterInstancesFromElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeregisterInstancesFromElasticLoadBalancerInput{
				ElasticLoadBalancerId:   nifcloud.String("test_elb_id"),
				ElasticLoadBalancerPort: nifcloud.Int64(1),
				InstancePort:            nifcloud.Int64(1),
				Protocol:                computing.ProtocolOfNiftyDeregisterInstancesFromElasticLoadBalancerRequest("test_protocol"),
				Instances: []computing.RequestInstancesOfNiftyDeregisterInstancesFromElasticLoadBalancer{
					{
						InstanceId: nifcloud.String("test_instances"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeregisterInstancesFromElasticLoadBalancerInput(tt.args, []interface{}{"test_instances"})
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyRegisterInstancesWithElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":        "test_elb_id",
		"instance_port": 1,
		"protocol":      "test_protocol",
		"lb_port":       1,
		"instances":     []interface{}{"test_instances"},
	})
	rd.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyRegisterInstancesWithElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyRegisterInstancesWithElasticLoadBalancerInput{
				ElasticLoadBalancerId:   nifcloud.String("test_elb_id"),
				ElasticLoadBalancerPort: nifcloud.Int64(1),
				InstancePort:            nifcloud.Int64(1),
				Protocol:                computing.ProtocolOfNiftyRegisterInstancesWithElasticLoadBalancerRequest("test_protocol"),
				Instances: []computing.RequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer{
					{
						InstanceId: nifcloud.String("test_instances"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyRegisterInstancesWithElasticLoadBalancerInput(tt.args, []interface{}{"test_instances"})
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyReplaceRouteTableAssociationWithElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":                     "test_elb_id",
		"route_table_id":             "test_route_table_id",
		"route_table_association_id": "test_route_table_association_id",
	})
	rd.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyReplaceRouteTableAssociationWithElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyReplaceRouteTableAssociationWithElasticLoadBalancerInput{
				RouteTableId:  nifcloud.String("test_route_table_id"),
				AssociationId: nifcloud.String("test_route_table_association_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyReplaceRouteTableAssociationWithElasticLoadBalancerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDisassociateRouteTableFromElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":                     "test_elb_id",
		"route_table_id":             "test_route_table_id",
		"route_table_association_id": "test_route_table_association_id",
	})
	rd.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDisassociateRouteTableFromElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDisassociateRouteTableFromElasticLoadBalancerInput{
				AssociationId: nifcloud.String("test_route_table_association_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDisassociateRouteTableFromElasticLoadBalancerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyAssociateRouteTableWithElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":                     "test_elb_id",
		"route_table_id":             "test_route_table_id",
		"route_table_association_id": "test_route_table_association_id",
	})
	rd.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyAssociateRouteTableWithElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyAssociateRouteTableWithElasticLoadBalancerInput{
				RouteTableId:          nifcloud.String("test_route_table_id"),
				ElasticLoadBalancerId: nifcloud.String("test_elb_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyAssociateRouteTableWithElasticLoadBalancerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestExpandNiftyDeleteNiftyElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":        "test_elb_id",
		"instance_port": 1,
		"protocol":      "test_protocol",
		"lb_port":       1,
	})
	rd.SetId("test_elb_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteElasticLoadBalancerInput{
				ElasticLoadBalancerId:   nifcloud.String("test_elb_id"),
				ElasticLoadBalancerPort: nifcloud.Int64(1),
				InstancePort:            nifcloud.Int64(1),
				Protocol:                computing.ProtocolOfNiftyDeleteElasticLoadBalancerRequest("test_protocol"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteNiftyElasticLoadBalancerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
