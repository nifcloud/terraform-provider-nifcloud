package acc

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_elastic_ip", &resource.Sweeper{
		Name: "nifcloud_elastic_ip",
		F:    testSweepElasticIP,
		Dependencies: []string{
			"nifcloud_instance",
		},
	})
}

func TestAcc_ElasticIP(t *testing.T) {
	var elasticIP computing.AddressesSet

	resourceName := "nifcloud_elastic_ip.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccElasticIPResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticIP(t, "testdata/elastic_ip.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticIPExists(resourceName, &elasticIP),
					testAccCheckElasticIPValues(&elasticIP),
					resource.TestCheckResourceAttr(resourceName, "ip_type", "false"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
				),
			},
			{
				Config: testAccElasticIP(t, "testdata/elastic_ip_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticIPExists(resourceName, &elasticIP),
					testAccCheckElasticIPValuesUpdated(&elasticIP),
					resource.TestCheckResourceAttr(resourceName, "ip_type", "false"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo-upd"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
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

func testAccElasticIP(t *testing.T, fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func testAccCheckElasticIPExists(n string, elasticIP *computing.AddressesSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no elasticIP resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no elasticIP id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeAddressesRequest(&computing.DescribeAddressesInput{
			PublicIp: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.AddressesSet) == 0 {
			return fmt.Errorf("elasticIP does not found in cloud: %s", saved.Primary.ID)
		}

		foundElasticIP := res.AddressesSet[0]

		if nifcloud.StringValue(foundElasticIP.PublicIp) != saved.Primary.ID {
			return fmt.Errorf("elasticIP does not found in cloud: %s", saved.Primary.ID)
		}

		*elasticIP = foundElasticIP
		return nil
	}
}

func testAccCheckElasticIPValues(elasticIP *computing.AddressesSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(elasticIP.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", elasticIP.Description)
		}

		if nifcloud.StringValue(elasticIP.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", elasticIP.AvailabilityZone)
		}

		if nifcloud.StringValue(elasticIP.PublicIp) == "" {
			return fmt.Errorf("bad public_ip state, expected not nil, got: nil")
		}
		return nil
	}
}

func testAccCheckElasticIPValuesUpdated(elasticIP *computing.AddressesSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(elasticIP.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", elasticIP.Description)
		}

		if nifcloud.StringValue(elasticIP.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", elasticIP.AvailabilityZone)
		}

		if nifcloud.StringValue(elasticIP.PublicIp) == "" {
			return fmt.Errorf("bad public_ip state, expected not nil, got: nil")
		}
		return nil
	}

}

func testAccElasticIPResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_elastic_ip" {
			continue
		}

		res, err := svc.DescribeAddressesRequest(&computing.DescribeAddressesInput{
			PublicIp: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return fmt.Errorf("failed DescribeAddressesRequest: %s", err)
		}

		if len(res.AddressesSet) > 0 {
			return fmt.Errorf("elasticIP (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepElasticIP(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeAddressesRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepPrivateElasticIPs []string
	var sweepPublicElasticIPs []string
	for _, k := range res.AddressesSet {
		if strings.HasPrefix(nifcloud.StringValue(k.Description), prefix) {
			if nifcloud.StringValue(k.PrivateIpAddress) != "" {
				sweepPrivateElasticIPs = append(sweepPrivateElasticIPs, nifcloud.StringValue(k.PrivateIpAddress))
			} else if nifcloud.StringValue(k.PublicIp) != "" {
				sweepPublicElasticIPs = append(sweepPublicElasticIPs, nifcloud.StringValue(k.PublicIp))
			}
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepPrivateElasticIPs {
		privateIP := n
		eg.Go(func() error {
			_, err := svc.ReleaseAddressRequest(&computing.ReleaseAddressInput{
				PrivateIpAddress: nifcloud.String(privateIP),
			}).Send(ctx)
			return err
		})
	}
	for _, n := range sweepPublicElasticIPs {
		publicIP := n
		eg.Go(func() error {
			_, err := svc.ReleaseAddressRequest(&computing.ReleaseAddressInput{
				PublicIp: nifcloud.String(publicIP),
			}).Send(ctx)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
