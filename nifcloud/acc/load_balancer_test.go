package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/acc/helper"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_load_balancer", &resource.Sweeper{
		Name: "nifcloud_load_balancer",
		F:    testSweepLoadBalancer,
	})
}

func TestAcc_LoadBalancer(t *testing.T) {
	var loadBalancer types.LoadBalancerDescriptions

	instanceName := prefix + acctest.RandString(7)

	resourceName := "nifcloud_load_balancer.basic"
	randName := prefix + acctest.RandString(7)
	sshKey := prefix + acctest.RandString(7)

	caKey := helper.GeneratePrivateKey(t, 2048)
	caCert := helper.GenerateSelfSignedCertificateAuthority(t, caKey)
	key := helper.GeneratePrivateKey(t, 2048)
	cert := helper.GenerateCertificate(t, caKey, caCert, key, randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccLoadBalancerResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancer(t, "testdata/load_balancer.tf", randName, instanceName, sshKey, cert, key, caCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(randName, 80, 443, &loadBalancer),
					testAccCheckLoadBalancerValues(&loadBalancer, randName, cert, instanceName),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_name", randName),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_volume", "10"),
					resource.TestCheckResourceAttr(resourceName, "balancing_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", "443"),
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
				Config: testAccLoadBalancer(t, "testdata/load_balancer_update.tf", randName, instanceName, sshKey, cert, key, caCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(randName, 80, 80, &loadBalancer),
					testAccCheckLoadBalancerValuesUpdated(&loadBalancer, randName, cert, instanceName),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_name", randName),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_volume", "20"),
					resource.TestCheckResourceAttr(resourceName, "balancing_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", "80"),
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
				ImportStateIdFunc: testAccLoadBalancerImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoadBalancer(t *testing.T, fileName, rName, instanceName, sshKey, certificate, key, ca string) string {
	b, err := os.ReadFile(fileName)
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

func testAccCheckLoadBalancerExists(lbName string, lbPort, instancePort int, loadBalancer *types.LoadBalancerDescriptions) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeLoadBalancers(context.Background(), &computing.DescribeLoadBalancersInput{
			LoadBalancerNames: &types.ListOfRequestLoadBalancerNames{
				Member: []types.RequestLoadBalancerNames{
					{
						LoadBalancerName: nifcloud.String(lbName),
						LoadBalancerPort: nifcloud.Int32(int32(lbPort)),
						InstancePort:     nifcloud.Int32(int32(instancePort)),
					},
				},
			},
		})

		if err != nil {
			return err
		}
		if res == nil || len(res.DescribeLoadBalancersResult.LoadBalancerDescriptions) == 0 {
			return fmt.Errorf("load_balancer does not found in cloud: %s", lbName)
		}

		foundLoadBalancer := res.DescribeLoadBalancersResult.LoadBalancerDescriptions[0]

		if nifcloud.ToString(foundLoadBalancer.LoadBalancerName) != lbName {
			return fmt.Errorf("load_balancer does not found in cloud: %s", lbName)
		}

		*loadBalancer = foundLoadBalancer
		return nil
	}
}

func testAccCheckLoadBalancerValues(loadBalancer *types.LoadBalancerDescriptions, rName, cert, iName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := loadBalancer.ListenerDescriptions[0].Listener

		if nifcloud.ToString(loadBalancer.NextMonthAccountingType) != "1" {
			return fmt.Errorf("bad accounting_type state, expected \"1\", got: %#v", loadBalancer.NextMonthAccountingType)
		}

		if nifcloud.ToString(loadBalancer.LoadBalancerName) != rName {
			return fmt.Errorf("bad load_balancer_name state,  expected \"%s\", got: %#v", rName, loadBalancer.LoadBalancerName)
		}

		if nifcloud.ToInt32(loadBalancer.NetworkVolume) != 10 {
			return fmt.Errorf("bad network_volume state,  expected \"10\", got: %#v", loadBalancer.NetworkVolume)
		}

		if nifcloud.ToInt32(listener.BalancingType) != 1 {
			return fmt.Errorf("bad balancing_type state, expected \"1\", got: %#v", listener.BalancingType)
		}

		if nifcloud.ToInt32(listener.InstancePort) != 443 {
			return fmt.Errorf("bad instance_port state, expected \"443\", got: %#v", listener.InstancePort)
		}

		if nifcloud.ToInt32(listener.LoadBalancerPort) != 80 {
			return fmt.Errorf("bad load_balancer_port state, expected \"80\", got: %#v", listener.LoadBalancerPort)
		}

		if nifcloud.ToInt32(loadBalancer.HealthCheck.HealthyThreshold) != 1 {
			return fmt.Errorf("bad healthy_threshold state, expected \"1\", got: %#v", loadBalancer.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.ToInt32(loadBalancer.HealthCheck.UnhealthyThreshold) != 2 {
			return fmt.Errorf("bad unhealthy_threshold state, expected \"2\", got: %#v", loadBalancer.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.ToString(loadBalancer.HealthCheck.Target) != "TCP:80" {
			return fmt.Errorf("bad health_check_target state, expected \"TCP:80\", got: %#v", loadBalancer.HealthCheck.Target)
		}

		if nifcloud.ToInt32(loadBalancer.HealthCheck.Interval) != 10 {
			return fmt.Errorf("bad health_check_interval state, expected \"10\", got: %#v", loadBalancer.HealthCheck.Interval)
		}

		if nifcloud.ToString(loadBalancer.Instances[0].InstanceId) != iName {
			return fmt.Errorf("bad instances state, expected \"%s\", got: %#v", rName, loadBalancer.Instances[0].InstanceId)
		}

		if nifcloud.ToBool(loadBalancer.Option.SessionStickinessPolicy.Enabled) != true {
			return fmt.Errorf("bad session_stickiness_policy_enable state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.Enabled)
		}

		if nifcloud.ToInt32(loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod) != 5 {
			return fmt.Errorf("bad session_stickiness_policy_expiration_period state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.ToBool(loadBalancer.Option.SorryPage.Enabled) != true {
			return fmt.Errorf("bad sorry_page_enable state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.ToString(listener.SSLCertificateId) == "" {
			return fmt.Errorf("bad ssl_certificate_id state, expected \"not null\", got: %#v", listener.SSLCertificateId)
		}

		if nifcloud.ToInt32(loadBalancer.Option.SorryPage.StatusCode) != 503 {
			return fmt.Errorf("bad sorry_page_status_code state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.ToString(loadBalancer.PolicyType) != "standard" {
			return fmt.Errorf("bad policy_type state,  expected \"standard\", got: %#v", loadBalancer.PolicyType)
		}

		if nifcloud.ToString(loadBalancer.Filter.FilterType) != "1" {
			return fmt.Errorf("bad filter_type state,  expected \"1\", got: %#v", loadBalancer.PolicyType)
		}

		if nifcloud.ToString(loadBalancer.Filter.IPAddresses[0].IPAddress) != "192.168.1.1" {
			return fmt.Errorf("bad filter state,  expected \"192.168.1.1\", got: %#v", loadBalancer.PolicyType)
		}
		return nil
	}
}

func testAccCheckLoadBalancerValuesUpdated(loadBalancer *types.LoadBalancerDescriptions, rName, cert, iName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := loadBalancer.ListenerDescriptions[0].Listener

		if nifcloud.ToString(loadBalancer.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state, expected \"2\", got: %#v", loadBalancer.NextMonthAccountingType)
		}

		if nifcloud.ToString(loadBalancer.LoadBalancerName) != rName {
			return fmt.Errorf("bad load_balancer_name state,  expected \"%s\", got: %#v", rName, loadBalancer.LoadBalancerName)
		}

		if nifcloud.ToInt32(loadBalancer.NetworkVolume) != 20 {
			return fmt.Errorf("bad network_volume state,  expected \"10\", got: %#v", loadBalancer.NetworkVolume)
		}

		if nifcloud.ToInt32(listener.BalancingType) != 2 {
			return fmt.Errorf("bad balancing_type state, expected \"2\", got: %#v", listener.BalancingType)
		}

		if nifcloud.ToInt32(listener.InstancePort) != 80 {
			return fmt.Errorf("bad instance_port state, expected \"80\", got: %#v", listener.InstancePort)
		}

		if nifcloud.ToInt32(listener.LoadBalancerPort) != 80 {
			return fmt.Errorf("bad load_balancer_port state, expected \"80\", got: %#v", listener.LoadBalancerPort)
		}

		if nifcloud.ToInt32(loadBalancer.HealthCheck.HealthyThreshold) != 1 {
			return fmt.Errorf("bad healthy_threshold state, expected \"1\", got: %#v", loadBalancer.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.ToInt32(loadBalancer.HealthCheck.UnhealthyThreshold) != 3 {
			return fmt.Errorf("bad unhealthy_threshold state, expected \"3\", got: %#v", loadBalancer.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.ToString(loadBalancer.HealthCheck.Target) != "ICMP" {
			return fmt.Errorf("bad health_check_target state, expected \"ICMP\", got: %#v", loadBalancer.HealthCheck.Target)
		}

		if nifcloud.ToInt32(loadBalancer.HealthCheck.Interval) != 11 {
			return fmt.Errorf("bad health_check_interval state, expected \"11\", got: %#v", loadBalancer.HealthCheck.Interval)
		}

		if nifcloud.ToString(loadBalancer.Instances[0].InstanceId) != iName {
			return fmt.Errorf("bad instances state, expected \"%s\", got: %#v", rName, loadBalancer.Instances[0].InstanceId)
		}

		if nifcloud.ToBool(loadBalancer.Option.SessionStickinessPolicy.Enabled) != true {
			return fmt.Errorf("bad session_stickiness_policy_enable state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.Enabled)
		}

		if nifcloud.ToInt32(loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod) != 5 {
			return fmt.Errorf("bad session_stickiness_policy_expiration_period state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.ToBool(loadBalancer.Option.SorryPage.Enabled) != true {
			return fmt.Errorf("bad sorry_page_enable state, expected \"true\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.ToString(listener.SSLCertificateId) == "" {
			return fmt.Errorf("bad ssl_certificate_id state, expected \"not null\", got: %#v", listener.SSLCertificateId)
		}

		if nifcloud.ToInt32(loadBalancer.Option.SorryPage.StatusCode) != 200 {
			return fmt.Errorf("bad sorry_page_status_code state, expected \"200\", got: %#v", loadBalancer.Option.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.ToString(loadBalancer.PolicyType) != "standard" {
			return fmt.Errorf("bad policy_type state,  expected \"standard\", got: %#v", loadBalancer.PolicyType)
		}

		if nifcloud.ToString(loadBalancer.Filter.FilterType) != "2" {
			return fmt.Errorf("bad filter_type state,  expected \"2\", got: %#v", loadBalancer.PolicyType)
		}

		if nifcloud.ToString(loadBalancer.Filter.IPAddresses[0].IPAddress) != "192.168.1.2" {
			return fmt.Errorf("bad filter state,  expected \"192.168.1.2\", got: %#v", loadBalancer.PolicyType)
		}
		return nil
	}
}

func testAccLoadBalancerResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_load_balancer" {
			continue
		}
		lbns := []types.RequestLoadBalancerNames{
			{
				LoadBalancerName: nifcloud.String(rs.Primary.ID),
			},
		}
		res, err := svc.DescribeLoadBalancers(context.Background(), &computing.DescribeLoadBalancersInput{
			LoadBalancerNames: &types.ListOfRequestLoadBalancerNames{Member: lbns},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.LoadBalancer" {
				return nil
			}
			return fmt.Errorf("failed DescribeLoadBalancersRequest: %s", err)
		}

		if len(res.DescribeLoadBalancersResult.LoadBalancerDescriptions) > 0 {
			return fmt.Errorf("load_balancer does not found in cloud: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepLoadBalancer(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeLoadBalancers(ctx, nil)
	if err != nil {
		return err
	}

	type lb struct {
		name         *string
		lbPort       *int32
		instancePort *int32
	}

	var sweepLBs []lb
	for _, b := range res.DescribeLoadBalancersResult.LoadBalancerDescriptions {
		for _, l := range b.ListenerDescriptions {
			if strings.HasPrefix(nifcloud.ToString(b.LoadBalancerName), prefix) {
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
			_, err := svc.DeleteLoadBalancer(ctx, &computing.DeleteLoadBalancerInput{
				LoadBalancerName: elb.name,
				LoadBalancerPort: elb.lbPort,
				InstancePort:     elb.instancePort,
			})
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

func testAccLoadBalancerImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		lbID := rs.Primary.Attributes["load_balancer_name"]
		lbPort := rs.Primary.Attributes["load_balancer_port"]
		instancePort := rs.Primary.Attributes["instance_port"]

		var parts []string
		parts = append(parts, lbID)
		parts = append(parts, lbPort)
		parts = append(parts, instancePort)

		id := strings.Join(parts, "_")
		return id, nil
	}
}
