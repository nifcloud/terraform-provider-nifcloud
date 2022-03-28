package elblistener

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandNiftyRegisterPortWithElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":                             "test-elb-id_protocol_port_port",
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
	})
	rd.SetId("test-elb-id_protocol_port_port")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyRegisterPortWithElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyRegisterPortWithElasticLoadBalancerInput{
				ElasticLoadBalancerId: nifcloud.String("test-elb-id_protocol_port_port"),
				Listeners: &types.ListOfRequestListenersOfNiftyRegisterPortWithElasticLoadBalancer{
					Member: []types.RequestListenersOfNiftyRegisterPortWithElasticLoadBalancer{
						{
							Protocol:                types.ProtocolOfListenersForNiftyRegisterPortWithElasticLoadBalancer("test_protocol"),
							ElasticLoadBalancerPort: nifcloud.Int32(1),
							InstancePort:            nifcloud.Int32(1),
							BalancingType:           nifcloud.Int32(1),
							Description:             nifcloud.String("test_description"),
							SSLCertificateId:        nifcloud.String("test_ssl_certificate_id"),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyRegisterPortWithElasticLoadBalancerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDescribeElasticLoadBalancersInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":        "test-elb-id_protocol_port_port",
		"instance_port": 1,
		"protocol":      "test_protocol",
		"lb_port":       1,
	})
	rd.SetId("test-elb-id_protocol_port_port")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDescribeElasticLoadBalancersInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDescribeElasticLoadBalancersInput{
				ElasticLoadBalancers: &types.RequestElasticLoadBalancers{
					ListOfRequestElasticLoadBalancerId:   []string{"test-elb-id"},
					ListOfRequestElasticLoadBalancerPort: []int32{1},
					ListOfRequestInstancePort:            []int32{1},
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

func TestExpandNiftyConfigureElasticLoadBalancerHealthCheckInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":                "test-elb-id_protocol_port_port",
		"instance_port":         1,
		"protocol":              "test_protocol",
		"lb_port":               1,
		"unhealthy_threshold":   1,
		"health_check_target":   "test_health_check_target",
		"health_check_interval": 1,
	})
	rd.SetId("test-elb-id_protocol_port_port")

	rdWithHTTP := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":                             "test-elb-id",
		"instance_port":                      1,
		"protocol":                           "HTTP",
		"lb_port":                            1,
		"unhealthy_threshold":                1,
		"health_check_target":                "HTTPS:443",
		"health_check_interval":              1,
		"health_check_path":                  "test_health_check_path",
		"health_check_expectation_http_code": []interface{}{1},
	})
	rdWithHTTP.SetId("test-elb-id_protocol_port_port")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyConfigureElasticLoadBalancerHealthCheckInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyConfigureElasticLoadBalancerHealthCheckInput{
				ElasticLoadBalancerId:   nifcloud.String("test-elb-id"),
				ElasticLoadBalancerPort: nifcloud.Int32(1),
				InstancePort:            nifcloud.Int32(1),
				Protocol:                types.ProtocolOfNiftyConfigureElasticLoadBalancerHealthCheckRequest("test_protocol"),
				HealthCheck: &types.RequestHealthCheckOfNiftyConfigureElasticLoadBalancerHealthCheck{
					Target:             nifcloud.String("test_health_check_target"),
					Interval:           nifcloud.Int32(1),
					UnhealthyThreshold: nifcloud.Int32(1),
				},
			},
		},
		{
			name: "expands the resource data with http protocol",
			args: rdWithHTTP,
			want: &computing.NiftyConfigureElasticLoadBalancerHealthCheckInput{
				ElasticLoadBalancerId:   nifcloud.String("test-elb-id"),
				ElasticLoadBalancerPort: nifcloud.Int32(1),
				InstancePort:            nifcloud.Int32(1),
				Protocol:                types.ProtocolOfNiftyConfigureElasticLoadBalancerHealthCheckRequest("HTTP"),
				HealthCheck: &types.RequestHealthCheckOfNiftyConfigureElasticLoadBalancerHealthCheck{
					Target:             nifcloud.String("HTTPS:443"),
					Interval:           nifcloud.Int32(1),
					UnhealthyThreshold: nifcloud.Int32(1),
					Path:               nifcloud.String("test_health_check_path"),
					ListOfRequestExpectation: &types.ListOfRequestExpectation{
						Member: []types.RequestExpectation{
							{
								HttpCode: nifcloud.Int32(1),
							},
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
		"elb_id":                             "test-elb-id_protocol_port_port",
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
	})
	rd.SetId("test-elb-id_protocol_port_port")
	dn := r.Data(rd.State())

	rdWithHTTPS := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"availability_zone":                           "test_availability_zone",
		"accounting_type":                             "test_accounting_type",
		"network_volume":                              1,
		"description":                                 "test_description",
		"balancing_type":                              1,
		"instance_port":                               1,
		"protocol":                                    "HTTPS",
		"lb_port":                                     1,
		"ssl_certificate_id":                          "test_ssl_certificate_id",
		"unhealthy_threshold":                         1,
		"health_check_target":                         "test_health_check_target",
		"health_check_interval":                       1,
		"health_check_path":                           "test_health_check_path",
		"health_check_expectation_http_code":          []interface{}{1},
		"instances":                                   []interface{}{"test_instances"},
		"session_stickiness_policy_enable":            true,
		"session_stickiness_policy_method":            1,
		"session_stickiness_policy_expiration_period": 1,
		"sorry_page_enable":                           true,
		"sorry_page_redirect_url":                     "test_sorry_page_redirect_url",
		"route_table_id":                              "test_route_table_id",
		"route_table_association_id":                  "test_route_table_association_id",
		"dns_name":                                    "test_dns_name",
		"version":                                     "test_version",
		"elb_id":                                      "test-elb-id_protocol_port_port",
	})
	rdWithHTTPS.SetId("test-elb-id_protocol_port_port")
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
				ElasticLoadBalancerId:   nifcloud.String("test-elb-id"),
				ElasticLoadBalancerPort: nifcloud.Int32(1),
				InstancePort:            nifcloud.Int32(1),
				Protocol:                types.ProtocolOfNiftyModifyElasticLoadBalancerAttributesRequest("test_protocol"),
				LoadBalancerAttributes: &types.RequestLoadBalancerAttributes{
					RequestSession: &types.RequestSessionOfNiftyModifyElasticLoadBalancerAttributes{
						RequestStickinessPolicy: &types.RequestStickinessPolicyOfNiftyModifyElasticLoadBalancerAttributes{
							Enable: nifcloud.Bool(true),
							Method: types.MethodOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributes(
								"1",
							),
							ExpirationPeriod: nifcloud.Int32(1),
						},
					},
					RequestSorryPage: &types.RequestSorryPage{
						Enable:      nifcloud.Bool(true),
						RedirectUrl: nifcloud.String("test_sorry_page_redirect_url"),
					},
					ListOfRequestAdditionalAttributes: &types.ListOfRequestAdditionalAttributes{
						Member: []types.RequestAdditionalAttributes{
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesProtocol,
								Value: nifcloud.String("test_protocol"),
							},
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesDescription,
								Value: nifcloud.String("test_description"),
							},
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesElasticLoadBalancerPort,
								Value: nifcloud.String("1"),
							},
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesInstancePort,
								Value: nifcloud.String("1"),
							},
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesBalancingType,
								Value: nifcloud.String("1"),
							},
						},
					},
				},
			},
		},
		{
			name: "expands the resource data with https protocol",
			args: dnWithHTTPS,
			want: &computing.NiftyModifyElasticLoadBalancerAttributesInput{
				ElasticLoadBalancerId:   nifcloud.String("test-elb-id"),
				ElasticLoadBalancerPort: nifcloud.Int32(1),
				InstancePort:            nifcloud.Int32(1),
				Protocol:                types.ProtocolOfNiftyModifyElasticLoadBalancerAttributesRequest("HTTPS"),
				LoadBalancerAttributes: &types.RequestLoadBalancerAttributes{
					RequestSession: &types.RequestSessionOfNiftyModifyElasticLoadBalancerAttributes{
						RequestStickinessPolicy: &types.RequestStickinessPolicyOfNiftyModifyElasticLoadBalancerAttributes{
							Enable: nifcloud.Bool(true),
							Method: types.MethodOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributes(
								"1",
							),
							ExpirationPeriod: nifcloud.Int32(1),
						},
					},
					RequestSorryPage: &types.RequestSorryPage{
						Enable:      nifcloud.Bool(true),
						RedirectUrl: nifcloud.String("test_sorry_page_redirect_url"),
					},
					ListOfRequestAdditionalAttributes: &types.ListOfRequestAdditionalAttributes{
						Member: []types.RequestAdditionalAttributes{
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesProtocol,
								Value: nifcloud.String("HTTPS"),
							},
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesDescription,
								Value: nifcloud.String("test_description"),
							},
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesElasticLoadBalancerPort,
								Value: nifcloud.String("1"),
							},
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesInstancePort,
								Value: nifcloud.String("1"),
							},
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesBalancingType,
								Value: nifcloud.String("1"),
							},
							{
								Key:   types.KeyOfLoadBalancerAttributesForNiftyModifyElasticLoadBalancerAttributesSslCertificateId,
								Value: nifcloud.String("test_ssl_certificate_id"),
							},
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

func TestExpandNiftyDeregisterInstancesFromElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":        "test-elb-id_protocol_port_port",
		"instance_port": 1,
		"protocol":      "test_protocol",
		"lb_port":       1,
		"instances":     []interface{}{"test_instances"},
	})
	rd.SetId("test-elb-id_protocol_port_port")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeregisterInstancesFromElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeregisterInstancesFromElasticLoadBalancerInput{
				ElasticLoadBalancerId:   nifcloud.String("test-elb-id"),
				ElasticLoadBalancerPort: nifcloud.Int32(1),
				InstancePort:            nifcloud.Int32(1),
				Protocol:                types.ProtocolOfNiftyDeregisterInstancesFromElasticLoadBalancerRequest("test_protocol"),
				Instances: &types.ListOfRequestInstancesOfNiftyDeregisterInstancesFromElasticLoadBalancer{
					Member: []types.RequestInstancesOfNiftyDeregisterInstancesFromElasticLoadBalancer{
						{
							InstanceId: nifcloud.String("test_instances"),
						},
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
		"elb_id":        "test-elb-id_protocol_port_port",
		"instance_port": 1,
		"protocol":      "test_protocol",
		"lb_port":       1,
		"instances":     []interface{}{"test_instances"},
	})
	rd.SetId("test-elb-id_protocol_port_port")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyRegisterInstancesWithElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyRegisterInstancesWithElasticLoadBalancerInput{
				ElasticLoadBalancerId:   nifcloud.String("test-elb-id"),
				ElasticLoadBalancerPort: nifcloud.Int32(1),
				InstancePort:            nifcloud.Int32(1),
				Protocol:                types.ProtocolOfNiftyRegisterInstancesWithElasticLoadBalancerRequest("test_protocol"),
				Instances: &types.ListOfRequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer{
					Member: []types.RequestInstancesOfNiftyRegisterInstancesWithElasticLoadBalancer{
						{
							InstanceId: nifcloud.String("test_instances"),
						},
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

func TestExpandNiftyDeleteNiftyElasticLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"elb_id":        "test-elb-id_protocol_port_port",
		"instance_port": 1,
		"protocol":      "test_protocol",
		"lb_port":       1,
	})
	rd.SetId("test-elb-id_protocol_port_port")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteElasticLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteElasticLoadBalancerInput{
				ElasticLoadBalancerId:   nifcloud.String("test-elb-id"),
				ElasticLoadBalancerPort: nifcloud.Int32(1),
				InstancePort:            nifcloud.Int32(1),
				Protocol:                types.ProtocolOfNiftyDeleteElasticLoadBalancerRequest("test_protocol"),
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
