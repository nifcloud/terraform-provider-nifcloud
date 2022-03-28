package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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
	resource.AddTestSweepers("nifcloud_elb", &resource.Sweeper{
		Name: "nifcloud_elb",
		F:    testSweepELB,
	})
}

func TestAcc_ELB(t *testing.T) {
	var elb types.ElasticLoadBalancerDescriptions

	resourceName := "nifcloud_elb.basic"
	randName := prefix + acctest.RandString(7)

	caKey := helper.GeneratePrivateKey(t, 2048)
	caCert := helper.GenerateSelfSignedCertificateAuthority(t, caKey)
	key := helper.GeneratePrivateKey(t, 2048)
	cert := helper.GenerateCertificate(t, caKey, caCert, key, randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccELBResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccELB(t, "testdata/elb.tf", randName, cert, key, caCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBExists(resourceName, &elb),
					testAccCheckELBValues(&elb, randName),
					resource.TestCheckResourceAttr(resourceName, "elb_name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_volume", "20"),
					resource.TestCheckResourceAttr(resourceName, "balancing_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", "3000"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "lb_port", "443"),
					resource.TestCheckResourceAttrSet(resourceName, "ssl_certificate_id"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "health_check_target", "HTTP:3000"),
					resource.TestCheckResourceAttr(resourceName, "health_check_interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "health_check_path", "/health"),
					resource.TestCheckResourceAttr(resourceName, "health_check_expectation_http_code.0", "200"),
					resource.TestCheckResourceAttr(resourceName, "instances.0", randName),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_method", "1"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_expiration_period", "4"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_redirect_url", "https://example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_association_id"),
					resource.TestCheckResourceAttrSet(resourceName, "dns_name"),
					resource.TestCheckResourceAttrSet(resourceName, "elb_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.1.network_name", randName),
					resource.TestCheckResourceAttr(resourceName, "network_interface.1.ip_address", "192.168.100.101"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.1.is_vip_network", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.network_id", "net-COMMON_GLOBAL"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.is_vip_network", "true"),
				),
			},
			{
				Config: testAccELB(t, "testdata/elb_update.tf", randName, cert, key, caCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBExists(resourceName, &elb),
					testAccCheckELBValuesUpdated(&elb, randName),
					resource.TestCheckResourceAttr(resourceName, "elb_name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_volume", "30"),
					resource.TestCheckResourceAttr(resourceName, "balancing_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", "3001"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "ssl_certificate_id", ""),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_threshold", "3"),
					resource.TestCheckResourceAttr(resourceName, "health_check_target", "HTTP:3001"),
					resource.TestCheckResourceAttr(resourceName, "health_check_interval", "11"),
					resource.TestCheckResourceAttr(resourceName, "health_check_path", "/health-upd"),
					resource.TestCheckResourceAttr(resourceName, "health_check_expectation_http_code.0", "302"),
					resource.TestCheckResourceAttr(resourceName, "instances.0", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_method", "2"),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_expiration_period", "5"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_redirect_url", "http://example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_association_id"),
					resource.TestCheckResourceAttrSet(resourceName, "dns_name"),
					resource.TestCheckResourceAttrSet(resourceName, "elb_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.1.network_name", randName),
					resource.TestCheckResourceAttr(resourceName, "network_interface.1.ip_address", "192.168.100.101"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.1.is_vip_network", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.network_id", "net-COMMON_GLOBAL"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.is_vip_network", "true"),
				),
			},
			{
				Config: testAccELB(t, "testdata/elb_reupdate.tf", randName, cert, key, caCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBExists(resourceName, &elb),
					resource.TestCheckResourceAttr(resourceName, "session_stickiness_policy_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "sorry_page_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "route_table_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccELBImportStateIDFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"network_interface",
				},
			},
		},
	})
}

func testAccELB(t *testing.T, fileName, rName, certificate, key, ca string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
		rName,
		rName,
		certificate,
		key,
		ca,
	)
}

func testAccCheckELBExists(n string, elb *types.ElasticLoadBalancerDescriptions) resource.TestCheckFunc {
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
				ListOfRequestElasticLoadBalancerId: []string{saved.Primary.ID},
			},
		})

		if err != nil {
			return err
		}

		if res == nil || len(res.NiftyDescribeElasticLoadBalancersResult.ElasticLoadBalancerDescriptions) == 0 {
			return fmt.Errorf("elb does not found in cloud: %s", saved.Primary.ID)
		}

		foundELB := res.NiftyDescribeElasticLoadBalancersResult.ElasticLoadBalancerDescriptions[0]

		if nifcloud.ToString(foundELB.ElasticLoadBalancerId) != saved.Primary.ID {
			return fmt.Errorf("elb does not found in cloud: %s", saved.Primary.ID)
		}

		*elb = foundELB
		return nil
	}
}

func testAccCheckELBValues(elb *types.ElasticLoadBalancerDescriptions, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := elb.ElasticLoadBalancerListenerDescriptions[0].Listener

		if nifcloud.ToString(elb.ElasticLoadBalancerName) != rName {
			return fmt.Errorf("bad elb_name state, expected \"%s\", got: %#v", rName, elb.ElasticLoadBalancerName)
		}

		if elb.AvailabilityZones[0] != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", elb.AvailabilityZones[0])
		}

		if nifcloud.ToString(elb.AccountingType) != "1" {
			return fmt.Errorf("bad accounting_type state, expected \"1\", got: %#v", elb.AccountingType)
		}

		if nifcloud.ToString(elb.NetworkVolume) != "20" {
			return fmt.Errorf("bad network_volume state, expected \"20\", got: %#v", elb.NetworkVolume)
		}

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

		if nifcloud.ToInt32(listener.HealthCheck.Expectation[0].HttpCode) != 200 {
			return fmt.Errorf("bad health_check_expectation_http_code state, expected \"/200\", got: %#v", listener.HealthCheck.Expectation[0].HttpCode)
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

		if nifcloud.ToString(elb.RouteTableId) == "" {
			return fmt.Errorf("bad route_table_id state, expected \"not null\", got: %#v", elb.RouteTableId)
		}
		return nil
	}
}

func testAccCheckELBValuesUpdated(elb *types.ElasticLoadBalancerDescriptions, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := elb.ElasticLoadBalancerListenerDescriptions[0].Listener

		if nifcloud.ToString(elb.ElasticLoadBalancerName) != rName+"upd" {
			return fmt.Errorf("bad elb_name state, expected \"%s\", got: %#v", rName+"upd", elb.ElasticLoadBalancerName)
		}

		if elb.AvailabilityZones[0] != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", elb.AvailabilityZones[0])
		}

		if nifcloud.ToString(elb.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state, expected \"2\", got: %#v", elb.NextMonthAccountingType)
		}

		if nifcloud.ToString(elb.NetworkVolume) != "30" {
			return fmt.Errorf("bad network_volume state, expected \"30\", got: %#v", elb.NetworkVolume)
		}

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

		if nifcloud.ToInt32(listener.ElasticLoadBalancerPort) != 80 {
			return fmt.Errorf("bad lb_port state, expected \"80\", got: %#v", listener.ElasticLoadBalancerPort)
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

		if nifcloud.ToInt32(listener.HealthCheck.Expectation[0].HttpCode) != 302 {
			return fmt.Errorf("bad health_check_expectation_http_code state, expected \"/302\", got: %#v", listener.HealthCheck.Expectation[0].HttpCode)
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

		if nifcloud.ToString(elb.RouteTableId) == "" {
			return fmt.Errorf("bad route_table_id state, expected \"not null\", got: %#v", elb.RouteTableId)
		}

		return nil
	}
}

func testAccELBResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_elb" {
			continue
		}

		res, err := svc.NiftyDescribeElasticLoadBalancers(context.Background(), &computing.NiftyDescribeElasticLoadBalancersInput{
			ElasticLoadBalancers: &types.RequestElasticLoadBalancers{
				ListOfRequestElasticLoadBalancerId: []string{rs.Primary.ID},
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
			return fmt.Errorf("elb (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepELB(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.NiftyDescribeElasticLoadBalancers(ctx, nil)
	if err != nil {
		return err
	}

	type elb struct {
		name         *string
		lbPort       *int32
		instancePort *int32
		protocol     types.ProtocolOfNiftyDeleteElasticLoadBalancerRequest
	}

	var sweepELBs []elb
	for _, e := range res.NiftyDescribeElasticLoadBalancersResult.ElasticLoadBalancerDescriptions {
		for _, l := range e.ElasticLoadBalancerListenerDescriptions {
			if strings.HasPrefix(nifcloud.ToString(e.ElasticLoadBalancerName), prefix) {
				sweepELBs = append(sweepELBs, elb{
					name:         e.ElasticLoadBalancerName,
					lbPort:       l.Listener.ElasticLoadBalancerPort,
					instancePort: l.Listener.InstancePort,
					protocol:     types.ProtocolOfNiftyDeleteElasticLoadBalancerRequest(nifcloud.ToString(l.Listener.Protocol)),
				})
			}
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, elb := range sweepELBs {
		elb := elb
		eg.Go(func() error {
			_, err := svc.NiftyDeleteElasticLoadBalancer(ctx, &computing.NiftyDeleteElasticLoadBalancerInput{
				ElasticLoadBalancerName: elb.name,
				ElasticLoadBalancerPort: elb.lbPort,
				InstancePort:            elb.instancePort,
				Protocol:                elb.protocol,
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func testAccELBImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		elbID := rs.Primary.Attributes["elb_id"]
		protocol := rs.Primary.Attributes["protocol"]
		lbPort := rs.Primary.Attributes["lb_port"]
		instancePort := rs.Primary.Attributes["instance_port"]

		var parts []string
		parts = append(parts, elbID)
		parts = append(parts, protocol)
		parts = append(parts, lbPort)
		parts = append(parts, instancePort)

		id := strings.Join(parts, "_")
		return id, nil
	}
}
