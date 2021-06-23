package elblistener

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"description":                                 "test_description",
		"balancing_type":                              1,
		"instance_port":                               1,
		"protocol":                                    "test_protocol",
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
		"elb_id":                                      "test-elb-id_protocol_port_port",
	})
	rd.SetId("test-elb-id_protocol_port_port")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.NiftyDescribeElasticLoadBalancersResponse
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
				res: &computing.NiftyDescribeElasticLoadBalancersResponse{
					NiftyDescribeElasticLoadBalancersOutput: &computing.NiftyDescribeElasticLoadBalancersOutput{
						ElasticLoadBalancerDescriptions: []computing.ElasticLoadBalancerDescriptions{
							{
								ElasticLoadBalancerId: nifcloud.String("test-elb-id"),
								ElasticLoadBalancerListenerDescriptions: []computing.ElasticLoadBalancerListenerDescriptions{
									{
										Listener: &computing.ListenerOfNiftyDescribeElasticLoadBalancers{
											Description:             nifcloud.String("test_description"),
											BalancingType:           nifcloud.Int64(1),
											InstancePort:            nifcloud.Int64(1),
											Protocol:                nifcloud.String("test_protocol"),
											ElasticLoadBalancerPort: nifcloud.Int64(1),
											SSLCertificateId:        nifcloud.String("test_ssl_certificate_id"),
											HealthCheck: &computing.HealthCheckOfNiftyDescribeElasticLoadBalancers{
												UnhealthyThreshold: nifcloud.Int64(1),
												Target:             nifcloud.String("test_health_check_target"),
												Interval:           nifcloud.Int64(1),
												Path:               nifcloud.String("test_health_check_path"),
												Expectation: []computing.Expectation{
													{
														HttpCode: nifcloud.Int64(1),
													},
												},
											},
											Instances: []computing.Instances{
												{
													InstanceId: nifcloud.String("test_instances"),
												},
											},
											SessionStickinessPolicy: &computing.SessionStickinessPolicyOfNiftyDescribeElasticLoadBalancers{
												Enabled:          nifcloud.Bool(true),
												Method:           nifcloud.Int64(1),
												ExpirationPeriod: nifcloud.Int64(1),
											},
											SorryPage: &computing.SorryPageOfNiftyDescribeElasticLoadBalancers{
												Enabled:     nifcloud.Bool(true),
												RedirectUrl: nifcloud.String("test_sorry_page_redirect_url"),
											},
										},
									},
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
				res: &computing.NiftyDescribeElasticLoadBalancersResponse{
					NiftyDescribeElasticLoadBalancersOutput: &computing.NiftyDescribeElasticLoadBalancersOutput{},
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
