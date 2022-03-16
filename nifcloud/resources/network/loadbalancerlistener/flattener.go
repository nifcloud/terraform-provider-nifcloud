package loadbalancerlistener

import (
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeLoadBalancersOutput) error {

	if res == nil || len(res.DescribeLoadBalancersResult.LoadBalancerDescriptions) == 0 {
		d.SetId("")
		return nil
	}
	loadBalancer := res.DescribeLoadBalancersResult.LoadBalancerDescriptions[0]
	if nifcloud.ToString(loadBalancer.LoadBalancerName) != d.Get("load_balancer_name") {
		return fmt.Errorf("unable to find load balancer within: %#v", loadBalancer.LoadBalancerName)
	}
	if err := d.Set("load_balancer_name", loadBalancer.LoadBalancerName); err != nil {
		return err
	}
	instances := make([]string, len(loadBalancer.Instances))
	for i, instance := range loadBalancer.Instances {
		instances[i] = nifcloud.ToString(instance.InstanceId)
	}
	if err := d.Set("instances", instances); err != nil {
		return err
	}
	if d.Get("filter") != nil && len(loadBalancer.Filter.IPAddresses) > 0 && *loadBalancer.Filter.IPAddresses[0].IPAddress != "*.*.*.*" {
		filters := make([]string, len(loadBalancer.Filter.IPAddresses))
		for i, filter := range loadBalancer.Filter.IPAddresses {
			filters[i] = nifcloud.ToString(filter.IPAddress)
		}
		if err := d.Set("filter", filters); err != nil {
			return err
		}
	}

	if err := d.Set("filter_type", loadBalancer.Filter.FilterType); err != nil {
		return err
	}

	if err := d.Set("healthy_threshold", loadBalancer.HealthCheck.HealthyThreshold); err != nil {
		return err
	}

	if err := d.Set("unhealthy_threshold", loadBalancer.HealthCheck.UnhealthyThreshold); err != nil {
		return err
	}

	if err := d.Set("health_check_target", loadBalancer.HealthCheck.Target); err != nil {
		return err
	}

	if err := d.Set("health_check_interval", loadBalancer.HealthCheck.Interval); err != nil {
		return err
	}

	listener := loadBalancer.ListenerDescriptions[0]
	if err := d.Set("instance_port", listener.Listener.InstancePort); err != nil {
		return err
	}

	if err := d.Set("load_balancer_port", listener.Listener.LoadBalancerPort); err != nil {
		return err
	}

	if err := d.Set("balancing_type", listener.Listener.BalancingType); err != nil {
		return err
	}

	if err := d.Set("ssl_certificate_id", listener.Listener.SSLCertificateId); err != nil {
		return err
	}

	if err := d.Set("policy_type", loadBalancer.PolicyType); err != nil {
		return err
	}

	if listener.Listener.SSLPolicy != nil {
		if err := d.Set("ssl_policy_id", listener.Listener.SSLPolicy.SSLPolicyId); err != nil {
			return err
		}

		if err := d.Set("ssl_policy_name", listener.Listener.SSLPolicy.SSLPolicyName); err != nil {
			return err
		}
	}

	if loadBalancer.Option != nil {
		if loadBalancer.Option.SessionStickinessPolicy != nil {
			if err := d.Set("session_stickiness_policy_enable", loadBalancer.Option.SessionStickinessPolicy.Enabled); err != nil {
				return err
			}
			if err := d.Set("session_stickiness_policy_expiration_period", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod); err != nil {
				return err
			}
		}

		if loadBalancer.Option.SorryPage != nil {
			if err := d.Set("sorry_page_enable", loadBalancer.Option.SorryPage.Enabled); err != nil {
				return err
			}
			if err := d.Set("sorry_page_status_code", loadBalancer.Option.SorryPage.StatusCode); err != nil {
				return err
			}
		}
	}

	if ip := net.ParseIP(nifcloud.ToString(loadBalancer.DNSName)); ip != nil {
		if ip.To4() != nil {
			if err := d.Set("ip_version", "v4"); err != nil {
				return err
			}
		} else {
			if err := d.Set("ip_version", "v6"); err != nil {
				return err
			}
		}
	}
	return nil
}
