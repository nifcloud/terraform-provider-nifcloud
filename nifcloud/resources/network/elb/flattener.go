package elb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.NiftyDescribeElasticLoadBalancersResponse) error {
	if res == nil || len(res.NiftyDescribeElasticLoadBalancersOutput.ElasticLoadBalancerDescriptions) == 0 {
		d.SetId("")
		return nil
	}

	elb := res.NiftyDescribeElasticLoadBalancersOutput.ElasticLoadBalancerDescriptions[0]

	if nifcloud.StringValue(elb.ElasticLoadBalancerId) != d.Id() {
		return fmt.Errorf(
			"unable to find elb within: %#v",
			res.NiftyDescribeElasticLoadBalancersOutput.ElasticLoadBalancerDescriptions,
		)
	}

	if len(elb.ElasticLoadBalancerListenerDescriptions) == 0 {
		d.SetId("")
		return nil
	}

	listener := elb.ElasticLoadBalancerListenerDescriptions[0].Listener

	if err := d.Set("elb_name", elb.ElasticLoadBalancerName); err != nil {
		return err
	}

	if err := d.Set("accounting_type", elb.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("availability_zone", elb.AvailabilityZones[0]); err != nil {
		return err
	}

	networkVolume, err := strconv.Atoi(nifcloud.StringValue(elb.NetworkVolume))
	if err != nil {
		return err
	}

	if err := d.Set("network_volume", networkVolume); err != nil {
		return err
	}

	if err := d.Set("description", listener.Description); err != nil {
		return err
	}

	if err := d.Set("balancing_type", listener.BalancingType); err != nil {
		return err
	}

	if err := d.Set("instance_port", listener.InstancePort); err != nil {
		return err
	}

	if err := d.Set("protocol", listener.Protocol); err != nil {
		return err
	}

	if err := d.Set("lb_port", listener.ElasticLoadBalancerPort); err != nil {
		return err
	}

	if err := d.Set("ssl_certificate_id", listener.SSLCertificateId); err != nil {
		return err
	}

	if err := d.Set("unhealthy_threshold", listener.HealthCheck.UnhealthyThreshold); err != nil {
		return err
	}

	if err := d.Set("health_check_target", listener.HealthCheck.Target); err != nil {
		return err
	}

	if err := d.Set("health_check_interval", listener.HealthCheck.Interval); err != nil {
		return err
	}

	if err := d.Set("health_check_path", listener.HealthCheck.Path); err != nil {
		return err
	}

	expectations := make([]int64, len(listener.HealthCheck.Expectation))
	for i, e := range listener.HealthCheck.Expectation {
		expectations[i] = nifcloud.Int64Value(e.HttpCode)
	}
	if err := d.Set("health_check_expectation_http_code", expectations); err != nil {
		return err
	}

	instances := make([]string, len(listener.Instances))
	for i, instance := range listener.Instances {
		instances[i] = nifcloud.StringValue(instance.InstanceId)
	}
	if err := d.Set("instances", instances); err != nil {
		return err
	}

	var networkInterfaces []map[string]interface{}
	for _, n := range elb.NetworkInterfaces {
		ni := map[string]interface{}{
			"is_vip_network": n.IsVipNetwork,
		}
		var findElm map[string]interface{}
		for _, dn := range d.Get("network_interface").(*schema.Set).List() {
			elm := dn.(map[string]interface{})

			if elm["network_id"] == nifcloud.StringValue(n.NetworkId) {
				findElm = elm
				break
			}

			if elm["network_name"] == nifcloud.StringValue(n.NetworkName) {
				findElm = elm
				break
			}
		}

		if findElm != nil {
			if findElm["ip_address"] != "" {
				ni["ip_address"] = n.IpAddress
			}

			if findElm["network_id"] != "" {
				ni["network_id"] = nifcloud.StringValue(n.NetworkId)
			} else {
				ni["network_name"] = nifcloud.StringValue(n.NetworkName)
			}

		}
		networkInterfaces = append(networkInterfaces, ni)
	}

	if err := d.Set("network_interface", networkInterfaces); err != nil {
		return err
	}

	if err := d.Set("session_stickiness_policy_enable", listener.SessionStickinessPolicy.Enabled); err != nil {
		return err
	}

	var sessionStickinessPolicyMethod *string
	if listener.SessionStickinessPolicy.Method != nil {
		sessionStickinessPolicyMethod = nifcloud.String(strconv.FormatInt(nifcloud.Int64Value(listener.SessionStickinessPolicy.Method), 10))
	}

	if err := d.Set("session_stickiness_policy_method", sessionStickinessPolicyMethod); err != nil {
		return err
	}

	if err := d.Set("session_stickiness_policy_expiration_period", listener.SessionStickinessPolicy.ExpirationPeriod); err != nil {
		return err
	}

	if err := d.Set("sorry_page_enable", listener.SorryPage.Enabled); err != nil {
		return err
	}

	if err := d.Set("sorry_page_redirect_url", listener.SorryPage.RedirectUrl); err != nil {
		return err
	}

	if err := d.Set("route_table_id", elb.RouteTableId); err != nil {
		return err
	}

	if err := d.Set("route_table_association_id", elb.RouteTableAssociationId); err != nil {
		return err
	}

	if err := d.Set("dns_name", elb.DNSName); err != nil {
		return err
	}

	if err := d.Set("elb_id", elb.ElasticLoadBalancerId); err != nil {
		return err
	}

	if err := d.Set("version", elb.VersionInformation.Version); err != nil {
		return err
	}

	if strings.Contains(d.Id(), "_") {
		// import so fix the id
		id := d.Get("elb_id").(string)
		d.SetId(id)
	}

	return nil
}
