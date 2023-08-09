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
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func TestAcc_ELBListener(t *testing.T) {
	var listener types.ListenerOfNiftyDescribeElasticLoadBalancers

	resourceName := "nifcloud_elb_listener.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccELBListenerResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccELBListener(t, "testdata/elb_listener.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBListenerExists(resourceName, randName, "HTTPS", 3000, 443, &listener),
					testAccCheckELBListenerValues(&listener, randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "balancing_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", "3000"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "lb_port", "443"),
					resource.TestCheckResourceAttrSet(resourceName, "ssl_certificate_id"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "health_check_target", "HTTP:3000"),
					resource.TestCheckResourceAttr(resourceName, "health_check_interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "health_check_path", "/health"),
					resource.TestCheckResourceAttr(resourceName, "health_check_expectation_http_code.0", "2xx"),
					resource.TestCheckResourceAttr(resourceName, "instances.0", randName),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_method", "1"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_expiration_period", "4"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_redirect_url", "https://example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "elb_id"),
				),
			},
			{
				Config: testAccELBListener(t, "testdata/elb_listener_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBListenerExists(resourceName, randName, "HTTP", 3001, 8080, &listener),
					testAccCheckELBListenerValuesUpdated(&listener, randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "balancing_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", "3001"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "ssl_certificate_id", ""),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "3"),
					resource.TestCheckResourceAttr(resourceName, "health_check_target", "HTTP:3001"),
					resource.TestCheckResourceAttr(resourceName, "health_check_interval", "11"),
					resource.TestCheckResourceAttr(resourceName, "health_check_path", "/health-upd"),
					resource.TestCheckResourceAttr(resourceName, "health_check_expectation_http_code.0", "3xx"),
					resource.TestCheckResourceAttr(resourceName, "instances.0", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_method", "2"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_expiration_period", "5"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_redirect_url", "http://example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "elb_id"),
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

func testAccELBListener(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
		rName,
		rName,
	)
}

func testAccCheckELBListenerExists(n, elbName, protocol string, InstancePort, LBPort int32, listener *types.ListenerOfNiftyDescribeElasticLoadBalancers) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no elb resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no elb id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.NiftyDescribeElasticLoadBalancers(context.Background(), &computing.NiftyDescribeElasticLoadBalancersInput{
			ElasticLoadBalancers: &types.RequestElasticLoadBalancers{
				ListOfRequestElasticLoadBalancerName: []string{elbName},
				ListOfRequestProtocol:                []string{protocol},
				ListOfRequestElasticLoadBalancerPort: []int32{LBPort},
				ListOfRequestInstancePort:            []int32{InstancePort},
			},
		})

		if err != nil {
			return err
		}

		if res == nil || len(res.NiftyDescribeElasticLoadBalancersResult.ElasticLoadBalancerDescriptions) == 0 {
			return fmt.Errorf("elb does not found in cloud: %s", saved.Primary.ID)
		}

		foundELB := res.NiftyDescribeElasticLoadBalancersResult.ElasticLoadBalancerDescriptions[0]
		foundListener := foundELB.ElasticLoadBalancerListenerDescriptions[0].Listener

		if nifcloud.ToString(foundELB.ElasticLoadBalancerId) != strings.Split(saved.Primary.ID, "_")[0] {
			return fmt.Errorf("elb does not found in cloud: %s", saved.Primary.ID)
		}

		*listener = *foundListener
		return nil
	}
}

func testAccCheckELBListenerValues(listener *types.ListenerOfNiftyDescribeElasticLoadBalancers, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if nifcloud.ToString(listener.Description) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", listener.Description)
		}

		if nifcloud.ToInt32(listener.BalancingType) != 2 {
			return fmt.Errorf("bad balancing_type state, expected \"2\", got: %#v", listener.BalancingType)
		}

		if nifcloud.ToInt32(listener.InstancePort) != 3000 {
			return fmt.Errorf("bad instance_port state, expected \"3000\", got: %#v", listener.InstancePort)
		}

		if nifcloud.ToString(listener.Protocol) != "HTTPS" {
			return fmt.Errorf("bad protocol state, expected \"HTTPS\", got: %#v", listener.Protocol)
		}

		if nifcloud.ToInt32(listener.ElasticLoadBalancerPort) != 443 {
			return fmt.Errorf("bad lb_port state, expected \"443\", got: %#v", listener.ElasticLoadBalancerPort)
		}

		if nifcloud.ToString(listener.SSLCertificateId) == "" {
			return fmt.Errorf("bad ssl_certificate_id state, expected \"not null\", got: %#v", listener.SSLCertificateId)
		}

		if nifcloud.ToInt32(listener.HealthCheck.UnhealthyThreshold) != 2 {
			return fmt.Errorf("bad unhealthy_threshold state, expected \"2\", got: %#v", listener.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.ToString(listener.HealthCheck.Target) != "HTTP:3000" {
			return fmt.Errorf("bad health_check_target state, expected \"HTTP:3000\", got: %#v", listener.HealthCheck.Target)
		}

		if nifcloud.ToInt32(listener.HealthCheck.Interval) != 10 {
			return fmt.Errorf("bad health_check_interval state, expected \"10\", got: %#v", listener.HealthCheck.Interval)
		}

		if nifcloud.ToString(listener.HealthCheck.Path) != "/health" {
			return fmt.Errorf("bad health_check_path state, expected \"/health\", got: %#v", listener.HealthCheck.Path)
		}

		if nifcloud.ToString(listener.HealthCheck.Expectation[0].HttpCode) != "2xx" {
			return fmt.Errorf("bad health_check_expectation_http_code state, expected \"/2xx\", got: %#v", listener.HealthCheck.Expectation[0].HttpCode)
		}

		if nifcloud.ToString(listener.Instances[0].InstanceId) != rName {
			return fmt.Errorf("bad instances state, expected \"%s\", got: %#v", rName, listener.Instances[0].InstanceId)
		}

		if nifcloud.ToBool(listener.SessionStickinessPolicy.Enabled) != true {
			return fmt.Errorf("bad session_stickiness_policy_enable state, expected \"true\", got: %#v", listener.SessionStickinessPolicy.Enabled)
		}

		if nifcloud.ToInt32(listener.SessionStickinessPolicy.Method) != 1 {
			return fmt.Errorf("bad session_stickiness_policy_method state, expected \"1\", got: %#v", listener.SessionStickinessPolicy.Method)
		}

		if nifcloud.ToInt32(listener.SessionStickinessPolicy.ExpirationPeriod) != 4 {
			return fmt.Errorf("bad session_stickiness_policy_expiration_period state, expected \"4\", got: %#v", listener.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.ToBool(listener.SorryPage.Enabled) != true {
			return fmt.Errorf("bad sorry_page_enable state, expected \"true\", got: %#v", listener.SorryPage.Enabled)
		}

		if nifcloud.ToString(listener.SorryPage.RedirectUrl) != "https://example.com" {
			return fmt.Errorf("bad session_stickiness_policy_expiration_period state, expected \"https://example.com\", got: %#v", listener.SorryPage.RedirectUrl)
		}
		return nil
	}
}

func testAccCheckELBListenerValuesUpdated(listener *types.ListenerOfNiftyDescribeElasticLoadBalancers, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(listener.Description) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", listener.Description)
		}

		if nifcloud.ToInt32(listener.BalancingType) != 1 {
			return fmt.Errorf("bad balancing_type state, expected \"1\", got: %#v", listener.BalancingType)
		}

		if nifcloud.ToInt32(listener.InstancePort) != 3001 {
			return fmt.Errorf("bad instance_port state, expected \"3001\", got: %#v", listener.InstancePort)
		}

		if nifcloud.ToString(listener.Protocol) != "HTTP" {
			return fmt.Errorf("bad protocol state, expected \"HTTP\", got: %#v", listener.Protocol)
		}

		if nifcloud.ToInt32(listener.ElasticLoadBalancerPort) != 8080 {
			return fmt.Errorf("bad lb_port state, expected \"8080\", got: %#v", listener.ElasticLoadBalancerPort)
		}

		if nifcloud.ToString(listener.SSLCertificateId) != "" {
			return fmt.Errorf("bad ssl_certificate_id state, expected \"null\", got: %#v", listener.SSLCertificateId)
		}

		if nifcloud.ToInt32(listener.HealthCheck.UnhealthyThreshold) != 3 {
			return fmt.Errorf("bad unhealthy_threshold state, expected \"3\", got: %#v", listener.HealthCheck.UnhealthyThreshold)
		}

		if nifcloud.ToString(listener.HealthCheck.Target) != "HTTP:3001" {
			return fmt.Errorf("bad health_check_target state, expected \"HTTP:3001\", got: %#v", listener.HealthCheck.Target)
		}

		if nifcloud.ToInt32(listener.HealthCheck.Interval) != 11 {
			return fmt.Errorf("bad health_check_interval state, expected \"11\", got: %#v", listener.HealthCheck.Interval)
		}

		if nifcloud.ToString(listener.HealthCheck.Path) != "/health-upd" {
			return fmt.Errorf("bad health_check_path state, expected \"/health\", got: %#v", listener.HealthCheck.Path)
		}

		if nifcloud.ToString(listener.HealthCheck.Expectation[0].HttpCode) != "3xx" {
			return fmt.Errorf("bad health_check_expectation_http_code state, expected \"/3xx\", got: %#v", listener.HealthCheck.Expectation[0].HttpCode)
		}

		if nifcloud.ToString(listener.Instances[0].InstanceId) != rName+"upd" {
			return fmt.Errorf("bad instances state, expected \"%s\", got: %#v", rName+"upd", listener.Instances[0].InstanceId)
		}

		if nifcloud.ToBool(listener.SessionStickinessPolicy.Enabled) != true {
			return fmt.Errorf("bad session_stickiness_policy_enable state, expected \"true\", got: %#v", listener.SessionStickinessPolicy.Enabled)
		}

		if nifcloud.ToInt32(listener.SessionStickinessPolicy.Method) != 2 {
			return fmt.Errorf("bad session_stickiness_policy_method state, expected \"2\", got: %#v", listener.SessionStickinessPolicy.Method)
		}

		if nifcloud.ToInt32(listener.SessionStickinessPolicy.ExpirationPeriod) != 5 {
			return fmt.Errorf("bad session_stickiness_policy_expiration_period state, expected \"5\", got: %#v", listener.SessionStickinessPolicy.ExpirationPeriod)
		}

		if nifcloud.ToBool(listener.SorryPage.Enabled) != true {
			return fmt.Errorf("bad sorry_page_enable state, expected \"true\", got: %#v", listener.SorryPage.Enabled)
		}

		if nifcloud.ToString(listener.SorryPage.RedirectUrl) != "http://example.com" {
			return fmt.Errorf("bad session_stickiness_policy_expiration_period state, expected \"http://example.com\", got: %#v", listener.SorryPage.RedirectUrl)
		}

		return nil
	}
}

func testAccELBListenerResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_elb_listener" {
			continue
		}

		res, err := svc.NiftyDescribeElasticLoadBalancers(context.Background(), &computing.NiftyDescribeElasticLoadBalancersInput{
			ElasticLoadBalancers: &types.RequestElasticLoadBalancers{
				ListOfRequestElasticLoadBalancerId: []string{strings.Split(rs.Primary.ID, "_")[0]},
			},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.ElasticLoadBalancer" {
				return nil
			}
			return fmt.Errorf("failed NiftyDescribeElasticLoadBalancersRequest: %s", err)
		}

		if len(res.NiftyDescribeElasticLoadBalancersResult.ElasticLoadBalancerDescriptions) > 0 {
			return fmt.Errorf("elb listener (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}
