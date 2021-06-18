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
	resource.AddTestSweepers("nifcloud_customer_gateway", &resource.Sweeper{
		Name: "nifcloud_customer_gateway",
		F:    testSweepCustomerGateway,
		Dependencies: []string{
			"nifcloud_vpn_connection",
		},
	})
}

func TestAcc_CustomerGateway(t *testing.T) {
	var customerGateway computing.CustomerGatewaySet

	resourceName := "nifcloud_customer_gateway.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCustomerGatewayResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomerGateway(t, "testdata/customer_gateway.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomerGatewayExists(resourceName, &customerGateway),
					testAccCheckCustomerGatewayValues(&customerGateway, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "lan_side_ip_address", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "lan_side_cidr_block", "192.168.0.0/28"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateway_id"),
				),
			},
			{
				Config: testAccCustomerGateway(t, "testdata/customer_gateway_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomerGatewayExists(resourceName, &customerGateway),
					testAccCheckCustomerGatewayValuesUpdated(&customerGateway, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memoupdated"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "lan_side_ip_address", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "lan_side_cidr_block", "192.168.0.0/28"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateway_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"type",
				},
			},
		},
	})
}

func testAccCustomerGateway(t *testing.T, fileName string, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), rName)
}

func testAccCheckCustomerGatewayExists(n string, customerGateway *computing.CustomerGatewaySet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no customerGateway resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no customerGateway id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeCustomerGatewaysRequest(&computing.DescribeCustomerGatewaysInput{
			CustomerGatewayId: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.CustomerGatewaySet) == 0 {
			return fmt.Errorf("customerGateway does not found in cloud: %s", saved.Primary.ID)
		}

		foundCustomerGateway := res.CustomerGatewaySet[0]

		if nifcloud.StringValue(foundCustomerGateway.CustomerGatewayId) != saved.Primary.ID {
			return fmt.Errorf("customerGateway does not found in cloud: %s", saved.Primary.ID)
		}

		*customerGateway = foundCustomerGateway
		return nil
	}
}

func testAccCheckCustomerGatewayValues(customerGateway *computing.CustomerGatewaySet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(customerGateway.NiftyCustomerGatewayName) != rName {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", rName, customerGateway.NiftyCustomerGatewayName)
		}

		if nifcloud.StringValue(customerGateway.NiftyCustomerGatewayDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", customerGateway.NiftyCustomerGatewayDescription)
		}

		if nifcloud.StringValue(customerGateway.IpAddress) != "192.168.0.1" {
			return fmt.Errorf("bad ip_address state, expected \"192.168.0.1\", got: %#v", customerGateway.IpAddress)
		}

		if nifcloud.StringValue(customerGateway.NiftyLanSideIpAddress) != "192.168.0.1" {
			return fmt.Errorf("bad lan_side_ip_address state, expected \"192.168.0.1\", got: %#v", customerGateway.NiftyLanSideIpAddress)
		}

		if nifcloud.StringValue(customerGateway.NiftyLanSideCidrBlock) != "192.168.0.0/28" {
			return fmt.Errorf("bad lan_side_cidr_block state, expected \"192.168.0.0/28\", got: %#v", customerGateway.NiftyLanSideCidrBlock)
		}
		return nil
	}
}

func testAccCheckCustomerGatewayValuesUpdated(customerGateway *computing.CustomerGatewaySet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(customerGateway.NiftyCustomerGatewayName) != rName+"upd" {
			return fmt.Errorf("bad name state, expected \"%supd\", got: %#v", rName, customerGateway.NiftyCustomerGatewayName)
		}

		if nifcloud.StringValue(customerGateway.NiftyCustomerGatewayDescription) != "memoupdated" {
			return fmt.Errorf("bad description state, expected \"memoupdated\", got: %#v", customerGateway.NiftyCustomerGatewayDescription)
		}

		if nifcloud.StringValue(customerGateway.IpAddress) != "192.168.0.1" {
			return fmt.Errorf("bad ip_address state, expected \"192.168.0.1\", got: %#v", customerGateway.IpAddress)
		}

		if nifcloud.StringValue(customerGateway.NiftyLanSideIpAddress) != "192.168.0.1" {
			return fmt.Errorf("bad lan_side_ip_address state, expected \"192.168.0.1\", got: %#v", customerGateway.NiftyLanSideIpAddress)
		}

		if nifcloud.StringValue(customerGateway.NiftyLanSideCidrBlock) != "192.168.0.0/28" {
			return fmt.Errorf("bad lan_side_cidr_block state, expected \"192.168.0.0/28\", got: %#v", customerGateway.NiftyLanSideCidrBlock)
		}
		return nil
	}

}

func testAccCustomerGatewayResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_customer_gateway" {
			continue
		}

		res, err := svc.DescribeCustomerGatewaysRequest(&computing.DescribeCustomerGatewaysInput{
			CustomerGatewayId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.CustomerGatewayId" {
				return nil
			}
			return fmt.Errorf("failed DescribeCustomerGatewaysRequest: %s", err)
		}

		if len(res.CustomerGatewaySet) > 0 {
			return fmt.Errorf("customerGateway (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepCustomerGateway(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeCustomerGatewaysRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepCustomerGateways []string
	for _, k := range res.CustomerGatewaySet {
		if strings.HasPrefix(nifcloud.StringValue(k.NiftyCustomerGatewayName), prefix) {
			sweepCustomerGateways = append(sweepCustomerGateways, nifcloud.StringValue(k.CustomerGatewayId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepCustomerGateways {
		customerGatewayID := n
		eg.Go(func() error {
			_, err := svc.DeleteCustomerGatewayRequest(&computing.DeleteCustomerGatewayInput{
				CustomerGatewayId: nifcloud.String(customerGatewayID),
			}).Send(ctx)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
