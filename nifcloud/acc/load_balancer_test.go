package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_load_balancer", &resource.Sweeper{
		Name: "nifcloud_load_balancer",
		F:    testSweepLoadBalancer,
		Dependencies: []string{
			"nifcloud_volume",
		},
	})
}

func TestAcc_LoadBalancer(t *testing.T) {
	var loadBalancer computing.LoadBalancerDescriptions

	resourceName := "nifcloud_load_balancer.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccLoadBalancerResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancer(t, "testdata/load_balancer.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &loadBalancer),
					testAccCheckLoadBalancerValues(&loadBalancer, randName),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_name", randName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "v4"),
					resource.TestCheckResourceAttr(resourceName, "network_volume", "10"),
					resource.TestCheckResourceAttr(resourceName, "policy_type", "standard"),
				),
			},
			{
				Config: testAccLoadBalancer(t, "testdata/load_balancer_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName, &loadBalancer),
					testAccCheckLoadBalancerValuesUpdated(&loadBalancer, randName),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_name", randName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "v4"),
					resource.TestCheckResourceAttr(resourceName, "network_volume", "20"),
					resource.TestCheckResourceAttr(resourceName, "policy_type", "standard"),
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

func testAccLoadBalancer(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccCheckLoadBalancerExists(n string, loadBalancer *computing.LoadBalancerDescriptions) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no load_balancer resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no load_balancer id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		lbns := []computing.RequestLoadBalancerNames{
			{
				LoadBalancerName: nifcloud.String(saved.Primary.ID),
			},
		}
		res, err := svc.DescribeLoadBalancersRequest(&computing.DescribeLoadBalancersInput{
			LoadBalancerNames: lbns,
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if res == nil || len(res.LoadBalancerDescriptions) == 0 || len(res.LoadBalancerDescriptions) == 0 {
			return fmt.Errorf("load_balancer does not found in cloud: %s", saved.Primary.ID)
		}

		foundLoadBalancer := res.LoadBalancerDescriptions[0]

		if nifcloud.StringValue(foundLoadBalancer.LoadBalancerName) != saved.Primary.ID {
			return fmt.Errorf("load_balancer does not found in cloud: %s", saved.Primary.ID)
		}

		*loadBalancer = foundLoadBalancer
		return nil
	}
}

func testAccCheckLoadBalancerValues(loadBalancer *computing.LoadBalancerDescriptions, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(loadBalancer.NextMonthAccountingType) != "1" {
			return fmt.Errorf("bad accounting_type state, expected \"1\", got: %#v", loadBalancer.NextMonthAccountingType)
		}

		if nifcloud.StringValue(&loadBalancer.AvailabilityZones[0]) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", loadBalancer.AvailabilityZones[0])
		}

		if nifcloud.StringValue(loadBalancer.LoadBalancerName) != rName {
			return fmt.Errorf("bad load_balancer_name state,  expected \"%s\", got: %#v", rName, loadBalancer.LoadBalancerName)
		}

		if nifcloud.Int64Value(loadBalancer.NetworkVolume) != 10 {
			return fmt.Errorf("bad network_volume state,  expected \"10\", got: %#v", loadBalancer.NetworkVolume)
		}

		if nifcloud.StringValue(loadBalancer.PolicyType) != "standard" {
			return fmt.Errorf("bad policy_type state,  expected \"standard\", got: %#v", loadBalancer.PolicyType)
		}
		return nil
	}
}

func testAccCheckLoadBalancerValuesUpdated(loadBalancer *computing.LoadBalancerDescriptions, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(loadBalancer.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state, expected \"1\", got: %#v", loadBalancer.NextMonthAccountingType)
		}

		if nifcloud.StringValue(&loadBalancer.AvailabilityZones[0]) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", loadBalancer.AvailabilityZones[0])
		}

		if nifcloud.StringValue(loadBalancer.LoadBalancerName) != rName {
			return fmt.Errorf("bad load_balancer_name state,  expected \"%s\", got: %#v", rName, loadBalancer.LoadBalancerName)
		}

		if nifcloud.Int64Value(loadBalancer.NetworkVolume) != 20 {
			return fmt.Errorf("bad availability_zone state,  expected \"20\", got: %#v", loadBalancer.NetworkVolume)
		}

		if nifcloud.StringValue(loadBalancer.PolicyType) != "standard" {
			return fmt.Errorf("bad policy_type state,  expected \"standard\", got: %#v", loadBalancer.PolicyType)
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
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.LoadBalancerName" {
				return fmt.Errorf("failed DescribeLoadBalancersRequest: %s", err)
			}
		}

		if len(res.LoadBalancerDescriptions) > 0 {
			return fmt.Errorf("load_balancer does not found in cloud: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepLoadBalancer(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeLoadBalancersRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range res.LoadBalancerDescriptions {
		if strings.HasPrefix(nifcloud.StringValue(n.LoadBalancerName), prefix) {
			eg.Go(func() error {
				_, err := svc.DeleteLoadBalancerRequest(&computing.DeleteLoadBalancerInput{
					LoadBalancerName: n.LoadBalancerName,
				}).Send(ctx)
				return err
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
