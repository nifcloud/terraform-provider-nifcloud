package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/acc/helper"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_load_balancer_listener", &resource.Sweeper{
		Name: "nifcloud_load_balancer",
		F:    testSweepLoadBalancerListener,
	})
}

func TestAcc_LoadBalancerListener(t *testing.T) {
	var loadBalancer computing.LoadBalancerDescriptions

	instanceName := prefix + acctest.RandString(7)

	resourceName := "nifcloud_load_balancer_listener.basic"
	randName := prefix + acctest.RandString(7)
	sshKey := prefix + acctest.RandString(7)

	caKey := helper.GeneratePrivateKey(t, 2048)
	caCert := helper.GenerateSelfSignedCertificateAuthority(t, caKey)
	key := helper.GeneratePrivateKey(t, 2048)
	cert := helper.GenerateCertificate(t, caKey, caCert, key, randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccLoadBalancerListenerResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerListener(t, "testdata/load_balancer_listener.tf", randName, instanceName, sshKey, cert, key, caCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerListenerExists(randName, 80, 8081, &loadBalancer),
					testAccCheckLoadBalancerListenerValues(&loadBalancer, randName, cert, instanceName),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_name", randName),
					resource.TestCheckResourceAttr(resourceName, "balancing_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", "8081"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "v4"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "policy_type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "filter_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "filter.0", "192.168.1.1"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "health_check_target", "TCP:80"),
					resource.TestCheckResourceAttr(resourceName, "health_check_interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "instances.0", instanceName),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_expiration_period", "5"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_status_code", "503"),
				),
			},
			{
				Config: testAccLoadBalancerListener(t, "testdata/load_balancer_listener_update.tf", randName, instanceName, sshKey, cert, key, caCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerListenerExists(randName, 80, 8083, &loadBalancer),
					testAccCheckLoadBalancerListenerValuesUpdated(&loadBalancer, randName, cert, instanceName),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_name", randName),
					resource.TestCheckResourceAttr(resourceName, "balancing_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", "8083"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "v4"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "policy_type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "filter_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "filter.0", "192.168.1.2"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "3"),
					resource.TestCheckResourceAttr(resourceName, "health_check_target", "ICMP"),
					resource.TestCheckResourceAttr(resourceName, "health_check_interval", "11"),
					resource.TestCheckResourceAttr(resourceName, "instances.0", instanceName),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_expiration_period", "5"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_status_code", "200"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoadBalancerListener(t *testing.T, fileName, rName, instanceName, sshKey, certificate, key, ca string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		instanceName,
		sshKey,
		certificate,
		key,
		ca,
	)
}

func testAccCheckLoadBalancerListenerExists(lbName string, lbPort, instancePort int, loadBalancer *computing.LoadBalancerDescriptions) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeLoadBalancersRequest(&computing.DescribeLoadBalancersInput{
			LoadBalancerNames: []computing.RequestLoadBalancerNames{
				{
					LoadBalancerName: nifcloud.String(lbName),
					LoadBalancerPort: nifcloud.Int64(int64(lbPort)),
					InstancePort:     nifcloud.Int64(int64(instancePort)),
				},
			},
		}).Send(context.Background())

		if err != nil {
			return err
		}
		if res == nil || len(res.DescribeLoadBalancersOutput.DescribeLoadBalancersResult.LoadBalancerDescriptions) == 0 {
			return fmt.Errorf("load_balancer does not found in cloud: %s", lbName)
		}

		foundLoadBalancer := res.DescribeLoadBalancersOutput.DescribeLoadBalancersResult.LoadBalancerDescriptions[0]

		if nifcloud.StringValue(foundLoadBalancer.LoadBalancerName) != lbName {
			return fmt.Errorf("load_balancer does not found in cloud: %s", lbName)
		}

		*loadBalancer = foundLoadBalancer
		return nil
	}
}

func testAccCheckLoadBalancerListenerValues(loadBalancer *computing.LoadBalancerDescriptions, rName, cert, iName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := loadBalancer.ListenerDescriptions[0].Listener
		log.Print("TFloadBalancer")
		log.Print(loadBalancer)
		if nifcloud.StringValue(loadBalancer.LoadBalancerName) != rName {
			return fmt.Errorf("bad load_balancer_name state,  expected \"%s\", got: %#v", rName, loadBalancer.LoadBalancerName)
		}

		if nifcloud.Int64Value(listener.BalancingType) != 1 {
			return fmt.Errorf("bad balancing_type state, expected \"1\", got: %#v", listener.BalancingType)
		}

		if nifcloud.Int64Value(listener.InstancePort) != 8081 {
			return fmt.Errorf("bad instance_port state, expected \"443\", got: %#v", listener.InstancePort)
		}

		if nifcloud.Int64Value(listener.LoadBalancerPort) != 80 {
			return fmt.Errorf("bad load_balancer_port state, expected \"80\", got: %#v", listener.LoadBalancerPort)
		}

		if nifcloud.Int64Value(loadBalancer.HealthCheck.HealthyThreshold) != 1 {
			return fmt.Errorf("bad healthy_threshold state, expected \"1\", got: %#v", loadBalancer.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.Int64Value(loadBalancer.HealthCheck.UnhealthyThreshold) != 2 {
			return fmt.Errorf("bad unhealthy_threshold state, expected \"2\", got: %#v", loadBalancer.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.StringValue(loadBalancer.HealthCheck.Target) != "TCP:80" {
			return fmt.Errorf("bad health_check_target state, expected \"TCP:80\", got: %#v", loadBalancer.HealthCheck.Target)
		}

		if nifcloud.Int64Value(loadBalancer.HealthCheck.Interval) != 10 {
			return fmt.Errorf("bad health_check_interval state, expected \"10\", got: %#v", loadBalancer.HealthCheck.Interval)
		}

		if nifcloud.StringValue(loadBalancer.Instances[0].InstanceId) != iName {
			return fmt.Errorf("bad instances state, expected \"%s\", got: %#v", rName, loadBalancer.Instances[0].InstanceId)
		}

		if nifcloud.BoolValue(loadBalancer.Option.SessionStickinessPolicy.Enabled) != true {
			return fmt.Errorf("bad session_stickiness_policy_enable state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.Enabled)
		}

		if nifcloud.Int64Value(loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod) != 5 {
			return fmt.Errorf("bad session_stickiness_policy_expiration_period state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.BoolValue(loadBalancer.Option.SorryPage.Enabled) != true {
			return fmt.Errorf("bad sorry_page_enable state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.StringValue(listener.SSLCertificateId) == "" {
			return fmt.Errorf("bad ssl_certificate_id state, expected \"not null\", got: %#v", listener.SSLCertificateId)
		}

		if nifcloud.Int64Value(loadBalancer.Option.SorryPage.StatusCode) != 503 {
			return fmt.Errorf("bad sorry_page_status_code state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.StringValue(loadBalancer.PolicyType) != "standard" {
			return fmt.Errorf("bad policy_type state,  expected \"standard\", got: %#v", loadBalancer.PolicyType)
		}

		if nifcloud.StringValue(loadBalancer.Filter.FilterType) != "1" {
			return fmt.Errorf("bad filter_type state,  expected \"1\", got: %#v", loadBalancer.PolicyType)
		}

		if nifcloud.StringValue(loadBalancer.Filter.IPAddresses[0].IPAddress) != "192.168.1.1" {
			return fmt.Errorf("bad filter state,  expected \"192.168.1.1\", got: %#v", loadBalancer.PolicyType)
		}
		return nil
	}
}

func testAccCheckLoadBalancerListenerValuesUpdated(loadBalancer *computing.LoadBalancerDescriptions, rName, cert, iName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := loadBalancer.ListenerDescriptions[0].Listener

		if nifcloud.StringValue(loadBalancer.LoadBalancerName) != rName {
			return fmt.Errorf("bad load_balancer_name state,  expected \"%s\", got: %#v", rName, loadBalancer.LoadBalancerName)
		}

		if nifcloud.Int64Value(listener.BalancingType) != 2 {
			return fmt.Errorf("bad balancing_type state, expected \"2\", got: %#v", listener.BalancingType)
		}

		if nifcloud.Int64Value(listener.InstancePort) != 8083 {
			return fmt.Errorf("bad instance_port state, expected \"80\", got: %#v", listener.InstancePort)
		}

		if nifcloud.Int64Value(listener.LoadBalancerPort) != 80 {
			return fmt.Errorf("bad load_balancer_port state, expected \"80\", got: %#v", listener.LoadBalancerPort)
		}

		if nifcloud.Int64Value(loadBalancer.HealthCheck.HealthyThreshold) != 1 {
			return fmt.Errorf("bad healthy_threshold state, expected \"1\", got: %#v", loadBalancer.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.Int64Value(loadBalancer.HealthCheck.UnhealthyThreshold) != 3 {
			return fmt.Errorf("bad unhealthy_threshold state, expected \"3\", got: %#v", loadBalancer.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.StringValue(loadBalancer.HealthCheck.Target) != "ICMP" {
			return fmt.Errorf("bad health_check_target state, expected \"ICMP\", got: %#v", loadBalancer.HealthCheck.Target)
		}

		if nifcloud.Int64Value(loadBalancer.HealthCheck.Interval) != 11 {
			return fmt.Errorf("bad health_check_interval state, expected \"11\", got: %#v", loadBalancer.HealthCheck.Interval)
		}

		if nifcloud.StringValue(loadBalancer.Instances[0].InstanceId) != iName {
			return fmt.Errorf("bad instances state, expected \"%s\", got: %#v", rName, loadBalancer.Instances[0].InstanceId)
		}

		if nifcloud.BoolValue(loadBalancer.Option.SessionStickinessPolicy.Enabled) != true {
			return fmt.Errorf("bad session_stickiness_policy_enable state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.Enabled)
		}

		if nifcloud.Int64Value(loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod) != 5 {
			return fmt.Errorf("bad session_stickiness_policy_expiration_period state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.BoolValue(loadBalancer.Option.SorryPage.Enabled) != true {
			return fmt.Errorf("bad session_stickiness_policy_enable state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.StringValue(listener.SSLCertificateId) == "" {
			return fmt.Errorf("bad ssl_certificate_id state, expected \"not null\", got: %#v", listener.SSLCertificateId)
		}

		if nifcloud.Int64Value(loadBalancer.Option.SorryPage.StatusCode) != 200 {
			return fmt.Errorf("bad sorry_page_status_code state, expected \"200\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.StringValue(loadBalancer.PolicyType) != "standard" {
			return fmt.Errorf("bad policy_type state,  expected \"standard\", got: %#v", loadBalancer.PolicyType)
		}

		if nifcloud.StringValue(loadBalancer.Filter.FilterType) != "2" {
			return fmt.Errorf("bad filter_type state,  expected \"2\", got: %#v", loadBalancer.PolicyType)
		}

		if nifcloud.StringValue(loadBalancer.Filter.IPAddresses[0].IPAddress) != "192.168.1.2" {
			return fmt.Errorf("bad filter state,  expected \"192.168.1.2\", got: %#v", loadBalancer.PolicyType)
		}
		return nil
	}
}

func testAccLoadBalancerListenerResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_load_balancer" {
			continue
		}
		lbns := []computing.RequestLoadBalancerNames{
			{
				LoadBalancerName: nifcloud.String(rs.Primary.ID),
			},
		}
		res, err := svc.DescribeLoadBalancersRequest(&computing.DescribeLoadBalancersInput{
			LoadBalancerNames: lbns,
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.LoadBalancer" {
				return nil
			}
			return fmt.Errorf("failed DescribeLoadBalancersRequest: %s", err)
		}

		if len(res.LoadBalancerDescriptions) > 0 {
			return fmt.Errorf("load_balancer does not found in cloud: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepLoadBalancerListener(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeLoadBalancersRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	type lb struct {
		name         *string
		lbPort       *int64
		instancePort *int64
	}

	var sweepLBs []lb
	for _, b := range res.DescribeLoadBalancersOutput.DescribeLoadBalancersResult.LoadBalancerDescriptions {
		for _, l := range b.ListenerDescriptions {
			if strings.HasPrefix(nifcloud.StringValue(b.LoadBalancerName), prefix) {
				sweepLBs = append(sweepLBs, lb{
					name:         b.LoadBalancerName,
					lbPort:       l.Listener.LoadBalancerPort,
					instancePort: l.Listener.InstancePort,
				})
			}
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, elb := range sweepLBs {
		elb := elb
		eg.Go(func() error {
			_, err := svc.DeleteLoadBalancerRequest(&computing.DeleteLoadBalancerInput{
				LoadBalancerName: elb.name,
				LoadBalancerPort: elb.lbPort,
				InstancePort:     elb.instancePort,
			}).Send(ctx)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
