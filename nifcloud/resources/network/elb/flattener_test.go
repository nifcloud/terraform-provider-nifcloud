package elb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
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
		"health_check_expectation_http_code": []interface{}{"test_health_check_expectation_http_code"},
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
			"system_ip_addresses": []interface{}{map[string]interface{}{
				"system_ip_address": "test_system_ip_address",
			}},
		}},
	})
	rd.SetId("test_elb_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.NiftyDescribeElasticLoadBalancersOutput
		d   *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want *schema.ResourceData
	}{
		{
			name: "flattens the response",
			args: args{
				d: rd,
				res: &computing.NiftyDescribeElasticLoadBalancersOutput{
					NiftyDescribeElasticLoadBalancersResult: &types.NiftyDescribeElasticLoadBalancersResult{
						ElasticLoadBalancerDescriptions: []types.ElasticLoadBalancerDescriptions{
							{
								ElasticLoadBalancerId:   nifcloud.String("test_elb_id"),
								ElasticLoadBalancerName: nifcloud.String("test_elb_name"),
								NextMonthAccountingType: nifcloud.String("test_accounting_type"),
								AvailabilityZones:       []string{"test_availability_zone"},
								NetworkVolume:           nifcloud.String("1"),
								ElasticLoadBalancerListenerDescriptions: []types.ElasticLoadBalancerListenerDescriptions{
									{
										Listener: &types.ListenerOfNiftyDescribeElasticLoadBalancers{
											Description:             nifcloud.String("test_description"),
											BalancingType:           nifcloud.Int32(1),
											InstancePort:            nifcloud.Int32(1),
											Protocol:                nifcloud.String("test_protocol"),
											ElasticLoadBalancerPort: nifcloud.Int32(1),
											SSLCertificateId:        nifcloud.String("test_ssl_certificate_id"),
											HealthCheck: &types.HealthCheckOfNiftyDescribeElasticLoadBalancers{
												UnhealthyThreshold: nifcloud.Int32(1),
												Target:             nifcloud.String("test_health_check_target"),
												Interval:           nifcloud.Int32(1),
												Path:               nifcloud.String("test_health_check_path"),
												Expectation: []types.Expectation{
													{
														HttpCode: nifcloud.String("test"),
													},
												},
											},
											Instances: []types.Instances{
												{
													InstanceId: nifcloud.String("test_instances"),
												},
											},
											SessionStickinessPolicy: &types.SessionStickinessPolicyOfNiftyDescribeElasticLoadBalancers{
												Enabled:          nifcloud.Bool(true),
												Method:           nifcloud.Int32(1),
												ExpirationPeriod: nifcloud.Int32(1),
											},
											SorryPage: &types.SorryPageOfNiftyDescribeElasticLoadBalancers{
												Enabled:     nifcloud.Bool(true),
												RedirectUrl: nifcloud.String("test_sorry_page_redirect_url"),
											},
										},
									},
								},
								NetworkInterfaces: []types.NetworkInterfaces{
									{
										IsVipNetwork: nifcloud.Bool(true),
										NetworkId:    nifcloud.String("test_network_id"),
										NetworkName:  nifcloud.String("test_network_name"),
										IpAddress:    nifcloud.String("test_ip_address"),
										SystemIpAddresses: []types.SystemIpAddresses{
											{
												SystemIpAddress: nifcloud.String("test_system_ip_address"),
											},
										},
									},
								},
								RouteTableId:            nifcloud.String("test_route_table_id"),
								RouteTableAssociationId: nifcloud.String("test_route_table_association_id"),
								DNSName:                 nifcloud.String("test_dns_name"),
								VersionInformation: &types.VersionInformation{
									Version: nifcloud.String("test_version"),
								},
							},
						},
					},
				},
			},
			want: rd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d: wantNotFoundRd,
				res: &computing.NiftyDescribeElasticLoadBalancersOutput{
					NiftyDescribeElasticLoadBalancersResult: &types.NiftyDescribeElasticLoadBalancersResult{},
				},
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			assert.NoError(t, err)

			if tt.args.res == nil {
				return
			}

			wantState := tt.want.State()
			if wantState == nil {
				tt.want.SetId("some")
				wantState = tt.want.State()
			}

			gotState := tt.args.d.State()
			if gotState == nil {
				tt.args.d.SetId("some")
				gotState = tt.args.d.State()
			}

			assert.Equal(t, wantState.Attributes, gotState.Attributes)
		})
	}
}
