package loadbalancer

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateLoadBalancerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type": "test_accounting_type",
		"availability_zones": []interface{}{
			"test_availability_zones",
		},
		"load_balancer_name": "test_load_balancer_name",
		"network_volume":     10,
		"ip_version":         "test_ip_version",
		"policy_type":        "test_policy_type",
	})
	rd.SetId("test_load_balancer_name")

	var lis []computing.RequestListeners
	n := computing.RequestListeners{
		BalancingType:    nifcloud.Int64(int64(1)),
		InstancePort:     nifcloud.Int64(80),
		LoadBalancerPort: nifcloud.Int64(80),
	}
	lis = append(lis, n)

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateLoadBalancerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateLoadBalancerInput{
				AccountingType: computing.AccountingTypeOfCreateLoadBalancerRequest("test_accounting_type"),
				AvailabilityZones: []string{
					"test_availability_zones",
				},
				LoadBalancerName: nifcloud.String("test_load_balancer_name"),
				Listeners:        lis,
				NetworkVolume:    nifcloud.Int64(int64(10)),
				IpVersion:        computing.IpVersionOfCreateLoadBalancerRequest("test_ip_version"),
				PolicyType:       computing.PolicyTypeOfCreateLoadBalancerRequest("test_policy_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateLoadBalancerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
