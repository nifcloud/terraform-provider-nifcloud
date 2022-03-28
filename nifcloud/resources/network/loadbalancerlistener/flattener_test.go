package loadbalancerlistener

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
		"accounting_type":    "test_accounting_type",
		"dns_name":           "test_dns_name",
		"load_balancer_name": "test_load_balancer_name",
		"filter":             []interface{}{"192.168.1.1"},
		"filter_type":        "test_filter_type",
		"instances":          []interface{}{"test_instance_id"},
		"instance_port":      80,
		"load_balancer_port": 80,
		"balancing_type":     1,
		"ssl_certificate_id": "test_ssl_certificate_id",
		"ssl_policy_id":      "test_ssl_policy_id",
		"ssl_policy_name":    "test_ssl_policy_name",
	})
	rd.SetId("test_load_balancer_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeLoadBalancersOutput
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
				res: &computing.DescribeLoadBalancersOutput{
					DescribeLoadBalancersResult: &types.DescribeLoadBalancersResult{
						LoadBalancerDescriptions: []types.LoadBalancerDescriptions{
							{
								AccountingType: nifcloud.String("1"),
								DNSName:        nifcloud.String("test_dns_name"),
								Filter: &types.Filter{
									FilterType: nifcloud.String("1"),
									IPAddresses: []types.IPAddresses{
										{IPAddress: nifcloud.String("192.168.1.1")},
									},
								},
								HealthCheck: &types.HealthCheckOfDescribeLoadBalancers{
									HealthyThreshold:   nifcloud.Int32(1),
									Interval:           nifcloud.Int32(5),
									Target:             nifcloud.String("test_target"),
									Timeout:            nifcloud.Int32(5),
									UnhealthyThreshold: nifcloud.Int32(1),
								},
								ListenerDescriptions: []types.ListenerDescriptions{
									{
										Listener: &types.Listener{
											InstancePort:     nifcloud.Int32(int32(80)),
											LoadBalancerPort: nifcloud.Int32(int32(80)),
											SSLPolicy: &types.SSLPolicy{
												SSLPolicyId:   nifcloud.String("test_ssl_policy_id"),
												SSLPolicyName: nifcloud.String("test_ssl_policy_name"),
											},
										},
									},
								},
								Instances: []types.Instances{
									{
										InstanceId: nifcloud.String("test_instance_id"),
									},
								},
								LoadBalancerName:        nifcloud.String("test_load_balancer_name"),
								NetworkVolume:           nifcloud.Int32(10),
								NextMonthAccountingType: nifcloud.String("test_accounting_type"),
								Policies: &types.Policies{
									AppCookieStickinessPolicies: []types.AppCookieStickinessPolicies{
										{
											CookieName: nifcloud.String("test_cookie_name"),
											PolicyName: nifcloud.String("test_app_policy_name"),
										},
									},
									LBCookieStickinessPolicies: []types.LBCookieStickinessPolicies{
										{
											CookieExpirationPeriod: nifcloud.String("test_cookie_expiration_period"),
											PolicyName:             nifcloud.String("test_lb_policy_name"),
										},
									},
								},
								PolicyType: nifcloud.String("test_policy_type"),
								Option: &types.Option{
									SessionStickinessPolicy: &types.SessionStickinessPolicy{
										Enabled: nifcloud.Bool(false),
									},
									SorryPage: &types.SorryPage{
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
				res: &computing.DescribeLoadBalancersOutput{
					DescribeLoadBalancersResult: &types.DescribeLoadBalancersResult{
						LoadBalancerDescriptions: []types.LoadBalancerDescriptions{},
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
