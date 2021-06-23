package loadbalancer

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                  "test_accounting_type",
		"dns_name":                         "test_dns_name",
		"load_balancer_name":               "test_load_balancer_name",
		"filter":                           []interface{}{"192.168.1.1"},
		"filter_type":                      "test_filter_type",
		"healthy_threshold":                1,
		"unhealthy_threshold":              1,
		"health_check_target":              "test_health_check_target",
		"health_check_interval":            5,
		"instances":                        []interface{}{"test_instance_id"},
		"instance_port":                    80,
		"load_balancer_port":               80,
		"balancing_type":                   1,
		"ssl_certificate_id":               "test_ssl_certificate_id",
		"ssl_policy_id":                    "test_ssl_policy_id",
		"ssl_policy_name":                  "test_ssl_policy_name",
		"session_stickiness_policy_enable": true,
		"session_stickiness_policy_expiration_period": 1,
		"sorry_page_enable":                           true,
		"sorry_page_status_code":                      503,
		"ip_version":                                  "test_ip_version",
		"network_volume":                              10,
		"policy_type":                                 "test_policy_type",
	})
	rd.SetId("test_load_balancer_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeLoadBalancersResponse
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
				res: &computing.DescribeLoadBalancersResponse{
					DescribeLoadBalancersOutput: &computing.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []computing.LoadBalancerDescriptions{
							{
								AccountingType: nifcloud.String("1"),
								DNSName:        nifcloud.String("test_dns_name"),
								Filter: &computing.Filter{
									FilterType: nifcloud.String("1"),
									IPAddresses: []computing.IPAddresses{
										{IPAddress: nifcloud.String("192.168.1.1")},
									},
								},
								HealthCheck: &computing.HealthCheckOfDescribeLoadBalancers{
									HealthyThreshold:   nifcloud.Int64(1),
									Interval:           nifcloud.Int64(5),
									Target:             nifcloud.String("test_target"),
									Timeout:            nifcloud.Int64(5),
									UnhealthyThreshold: nifcloud.Int64(1),
								},
								ListenerDescriptions: []computing.ListenerDescriptions{
									{
										Listener: &computing.Listener{
											InstancePort:     nifcloud.Int64(int64(80)),
											LoadBalancerPort: nifcloud.Int64(int64(80)),
											SSLPolicy: &computing.SSLPolicy{
												SSLPolicyId:   nifcloud.String("test_ssl_policy_id"),
												SSLPolicyName: nifcloud.String("test_ssl_policy_name"),
											},
										},
									},
								},
								Instances: []computing.Instances{
									{
										InstanceId: nifcloud.String("test_instance_id"),
									},
								},
								LoadBalancerName:        nifcloud.String("test_load_balancer_name"),
								NetworkVolume:           nifcloud.Int64(10),
								NextMonthAccountingType: nifcloud.String("test_accounting_type"),
								Policies: &computing.Policies{
									AppCookieStickinessPolicies: []computing.AppCookieStickinessPolicies{
										{
											CookieName: nifcloud.String("test_cookie_name"),
											PolicyName: nifcloud.String("test_app_policy_name"),
										},
									},
									LBCookieStickinessPolicies: []computing.LBCookieStickinessPolicies{
										{
											CookieExpirationPeriod: nifcloud.String("test_cookie_expiration_period"),
											PolicyName:             nifcloud.String("test_lb_policy_name"),
										},
									},
								},
								PolicyType: nifcloud.String("test_policy_type"),
								Option: &computing.Option{
									SessionStickinessPolicy: &computing.SessionStickinessPolicy{
										Enabled: nifcloud.Bool(false),
									},
									SorryPage: &computing.SorryPage{
										Enabled: nifcloud.Bool(false),
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
				res: &computing.DescribeLoadBalancersResponse{
					DescribeLoadBalancersOutput: &computing.DescribeLoadBalancersOutput{
						LoadBalancerDescriptions: []computing.LoadBalancerDescriptions{},
					},
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
