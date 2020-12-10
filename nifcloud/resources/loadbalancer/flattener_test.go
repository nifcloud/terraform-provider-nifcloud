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
		"accounting_type":    "test_accounting_type",
		"availability_zones": "test_availability_zones",
		"dns_name":           "test_dns_name",
		"load_balancer_name": "test_load_balancer_name",
		"network_volume":     10,
		"policy_type":        "test_policy_type",
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
								NextMonthAccountingType: nifcloud.String("test_accounting_type"),
								AvailabilityZones:       []string{"test_availability_zones"},
								DNSName:                 nifcloud.String("test_dns_name"),
								LoadBalancerName:        nifcloud.String("test_load_balancer_name"),
								NetworkVolume:           nifcloud.Int64(10),
								PolicyType:              nifcloud.String("test_policy_type"),
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
