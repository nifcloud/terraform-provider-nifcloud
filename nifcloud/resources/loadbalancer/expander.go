package loadbalancer

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandCreateLoadBalancerInput(d *schema.ResourceData) *computing.CreateLoadBalancerInput {
	at := d.Get("accounting_type").(string)
	ipv := d.Get("ip_version").(string)
	pt := d.Get("policy_type").(string)
	var azs []string
	availabilityZones := d.Get("availability_zones").([]interface{})
	for _, az := range availabilityZones {
		azs = append(azs, az.(string))
	}

	return &computing.CreateLoadBalancerInput{
		AccountingType:    computing.AccountingTypeOfCreateLoadBalancerRequest(at),
		AvailabilityZones: azs,
		LoadBalancerName:  nifcloud.String(d.Get("load_balancer_name").(string)),
		Listeners: []computing.RequestListeners{{
			BalancingType:    nifcloud.Int64(int64(1)),
			InstancePort:     nifcloud.Int64(80),
			LoadBalancerPort: nifcloud.Int64(80),
		}},
		NetworkVolume: nifcloud.Int64(int64(d.Get("network_volume").(int))),
		IpVersion:     computing.IpVersionOfCreateLoadBalancerRequest(ipv),
		PolicyType:    computing.PolicyTypeOfCreateLoadBalancerRequest(pt),
	}
}
