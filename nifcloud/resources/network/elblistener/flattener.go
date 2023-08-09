package elblistener

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.NiftyDescribeElasticLoadBalancersOutput) error {
	if res == nil || len(res.NiftyDescribeElasticLoadBalancersResult.ElasticLoadBalancerDescriptions) == 0 {
		d.SetId("")
		return nil
	}

	elb := res.NiftyDescribeElasticLoadBalancersResult.ElasticLoadBalancerDescriptions[0]

	if nifcloud.ToString(elb.ElasticLoadBalancerId) != getELBID(d) {
		return fmt.Errorf(
			"unable to find elb within: %#v",
			res.NiftyDescribeElasticLoadBalancersResult.ElasticLoadBalancerDescriptions,
		)
	}

	if len(elb.ElasticLoadBalancerListenerDescriptions) == 0 {
		d.SetId("")
		return nil
	}

	listener := elb.ElasticLoadBalancerListenerDescriptions[0].Listener

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

	expectations := make([]string, len(listener.HealthCheck.Expectation))
	for i, e := range listener.HealthCheck.Expectation {
		expectations[i] = nifcloud.ToString(e.HttpCode)
	}
	if err := d.Set("health_check_expectation_http_code", expectations); err != nil {
		return err
	}

	instances := make([]string, len(listener.Instances))
	for i, instance := range listener.Instances {
		instances[i] = nifcloud.ToString(instance.InstanceId)
	}
	if err := d.Set("instances", instances); err != nil {
		return err
	}

	if err := d.Set("session_stickiness_policy_enable", listener.SessionStickinessPolicy.Enabled); err != nil {
		return err
	}

	var sessionStickinessPolicyMethod *string
	if listener.SessionStickinessPolicy.Method != nil {
		sessionStickinessPolicyMethod = nifcloud.String(strconv.Itoa(int(nifcloud.ToInt32(listener.SessionStickinessPolicy.Method))))
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

	if err := d.Set("elb_id", elb.ElasticLoadBalancerId); err != nil {
		return err
	}
	return nil
}
