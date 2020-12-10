package loadbalancer

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeLoadBalancersResponse) error {

	if res == nil || len(res.LoadBalancerDescriptions) == 0 {
		d.SetId("")
		return nil
	}
	loadBalancer := res.LoadBalancerDescriptions[0]
	if nifcloud.StringValue(loadBalancer.LoadBalancerName) != d.Get("load_balancer_name") {
		return fmt.Errorf("unable to find load balancer within: %#v", loadBalancer.LoadBalancerName)
	}
	if err := d.Set("load_balancer_name", loadBalancer.LoadBalancerName); err != nil {
		return err
	}

	if err := d.Set("accounting_type", loadBalancer.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("availability_zones", loadBalancer.AvailabilityZones); err != nil {
		return err
	}

	if err := d.Set("dns_name", loadBalancer.DNSName); err != nil {
		return err
	}

	if err := d.Set("network_volume", loadBalancer.NetworkVolume); err != nil {
		return err
	}

	if err := d.Set("policy_type", loadBalancer.PolicyType); err != nil {
		return err
	}

	return nil
}
