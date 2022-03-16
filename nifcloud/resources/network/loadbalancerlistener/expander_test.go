package loadbalancerlistener

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandRegisterPortWithLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"load_balancer_name": "name",
		"instance_port":      80,
		"load_balancer_port": 80,
		"balancing_type":     1,
	})
	rd.SetId("name")

	var lis []types.RequestListenersOfRegisterPortWithLoadBalancer
	n := types.RequestListenersOfRegisterPortWithLoadBalancer{
		BalancingType:    nifcloud.Int32(int32(1)),
		InstancePort:     nifcloud.Int32(80),
		LoadBalancerPort: nifcloud.Int32(80),
	}
	lis = append(lis, n)

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.RegisterPortWithLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.RegisterPortWithLoadBalancerInput{
				LoadBalancerName: nifcloud.String("name"),
				Listeners:        &types.ListOfRequestListenersOfRegisterPortWithLoadBalancer{Member: lis},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRegisterPortWithLoadBalancerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeLoadBalancersInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_port":      80,
		"load_balancer_name": "name",
		"load_balancer_port": 80,
	})
	rd.SetId("name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeLoadBalancersInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeLoadBalancersInput{
				LoadBalancerNames: &types.ListOfRequestLoadBalancerNames{
					Member: []types.RequestLoadBalancerNames{
						{
							InstancePort:     nifcloud.Int32(int32(80)),
							LoadBalancerName: nifcloud.String("name"),
							LoadBalancerPort: nifcloud.Int32(int32(80)),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeLoadBalancersInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateLoadBalancer(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_port":      80,
		"load_balancer_name": "name",
		"load_balancer_port": 80,
	})
	rd.SetId("name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeLoadBalancersInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeLoadBalancersInput{
				LoadBalancerNames: &types.ListOfRequestLoadBalancerNames{
					Member: []types.RequestLoadBalancerNames{
						{
							InstancePort:     nifcloud.Int32(int32(80)),
							LoadBalancerName: nifcloud.String("name"),
							LoadBalancerPort: nifcloud.Int32(int32(80)),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeLoadBalancersInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
