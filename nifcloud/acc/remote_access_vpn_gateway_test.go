package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_ravgw", &resource.Sweeper{
		Name: "nifcloud_remote_access_vpn_gateway",
		F:    testSweepRemoteAccessVpnGateway,
	})
}

func TestAcc_RemoteAccessVpnGateway(t *testing.T) {
	var ravgw types.RemoteAccessVpnGatewaySetOfDescribeRemoteAccessVpnGateways

	resourceName := "nifcloud_remote_access_vpn_gateway.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccRemoteAccessVpnGatewayResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRemoteAccessVpnGateway(t, "testdata/remote_access_vpn_gateway.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemoteAccessVpnGatewayExists(resourceName, &ravgw),
					testAccCheckRemoteAccessVpnGatewayValues(&ravgw, randName),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.network_id"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.ip_address", "192.168.1.1"),
					resource.TestCheckResourceAttrSet(resourceName, "ssl_certificate_id"),
					resource.TestCheckResourceAttr(resourceName, "cipher_suite.0", "AES128-GCM-SHA256"),
					resource.TestCheckResourceAttr(resourceName, "type", "small"),
					resource.TestCheckResourceAttr(resourceName, "pool_network_cidr", "192.168.2.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "remote_access_vpn_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "user.0.name", "user1"),
					resource.TestCheckResourceAttr(resourceName, "user.0.description", "user1"),
					resource.TestCheckResourceAttrSet(resourceName, "user.0.password"),
					resource.TestCheckResourceAttr(resourceName, "user.1.name", "user2"),
					resource.TestCheckResourceAttr(resourceName, "user.1.description", "user2"),
					resource.TestCheckResourceAttrSet(resourceName, "user.1.password"),
				),
			},
			{
				Config: testAccRemoteAccessVpnGateway(t, "testdata/remote_access_vpn_gateway_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemoteAccessVpnGatewayExists(resourceName, &ravgw),
					testAccCheckRemoteAccessVpnGatewayValuesUpdated(&ravgw, randName),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"upd"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.network_id"),
					resource.TestCheckResourceAttrSet(resourceName, "remote_access_vpn_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.ip_address", "192.168.1.1"),
					resource.TestCheckResourceAttrSet(resourceName, "remote_access_vpn_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "user.0.name", "user3"),
					resource.TestCheckResourceAttr(resourceName, "user.0.description", "user3"),
					resource.TestCheckResourceAttrSet(resourceName, "user.0.password"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"user.0.password",
				},
			},
		},
	})
}

func testAccRemoteAccessVpnGateway(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
	)
}

func testAccCheckRemoteAccessVpnGatewayExists(n string, ravgw *types.RemoteAccessVpnGatewaySetOfDescribeRemoteAccessVpnGateways) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no ravgw resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no ravgw id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeRemoteAccessVpnGateways(context.Background(), &computing.DescribeRemoteAccessVpnGatewaysInput{
			RemoteAccessVpnGatewayId: []string{saved.Primary.ID},
		})
		if err != nil {
			return err
		}

		if res == nil || len(res.RemoteAccessVpnGatewaySet) == 0 {
			return fmt.Errorf("ravgw does not found in cloud: %s", saved.Primary.ID)
		}

		foundRemoteAccessVpnGateway := res.RemoteAccessVpnGatewaySet[0]

		if nifcloud.ToString(foundRemoteAccessVpnGateway.RemoteAccessVpnGatewayId) != saved.Primary.ID {
			return fmt.Errorf("ravgw does not found in cloud: %s", saved.Primary.ID)
		}

		*ravgw = foundRemoteAccessVpnGateway

		return nil
	}
}

func testAccCheckRemoteAccessVpnGatewayValues(ravgw *types.RemoteAccessVpnGatewaySetOfDescribeRemoteAccessVpnGateways, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(ravgw.AccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state,  expected \"2\", got: %#v", ravgw.AccountingType)
		}

		if nifcloud.ToString(ravgw.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", ravgw.AvailabilityZone)
		}

		if nifcloud.ToString(ravgw.PoolNetworkCidr) != "192.168.2.0/24" {
			return fmt.Errorf("bad pool_network_cidr state,  expected \"192.168.2.0/24\", got: %#v", ravgw.PoolNetworkCidr)
		}

		if nifcloud.ToString(ravgw.Description) != "memo" {
			return fmt.Errorf("bad description state,  expected \"memo\", got: %#v", ravgw.Description)
		}

		for _, ni := range ravgw.NetworkInterfaceSet {
			if nifcloud.ToString(ni.NiftyNetworkId) != "net-COMMON_GLOBAL" {
				if nifcloud.ToString(ni.NiftyNetworkId) == "" {
					return fmt.Errorf("bad network_interface.0.network_id state,  expected not empty.")
				}

				if nifcloud.ToString(ni.PrivateIpAddress) != "192.168.1.1" {
					return fmt.Errorf("bad network_interface.0.ip_address state,  expected \"192.168.1.1\", got: %#v", nifcloud.ToString(ni.PrivateIpAddress))
				}
			}
		}

		if len(ravgw.CipherSuiteSet) != 1 {
			return fmt.Errorf("bad cipher_suite length,  expected length 1, got %d", len(ravgw.CipherSuiteSet))
		}

		if nifcloud.ToString(ravgw.CipherSuiteSet[0].CipherSuite) != "AES128-GCM-SHA256" {
			return fmt.Errorf("bad cipher_suite state,  expected \"AES128-GCM-SHA256\", got: %#v", ravgw.CipherSuiteSet[0].CipherSuite)
		}

		if nifcloud.ToString(ravgw.SslCertificateId) == "" {
			return fmt.Errorf("bad ssl_certificate_id state,  expected not empty.")
		}

		if len(ravgw.RemoteUserSet) != 2 {
			return fmt.Errorf("bad user length,  expected length 2, got %d", len(ravgw.RemoteUserSet))
		}

		if nifcloud.ToString(ravgw.RemoteUserSet[0].UserName) != "user1" {
			return fmt.Errorf("bad user.0.name state,  expected \"user1\", got: %#v", ravgw.RemoteUserSet[0].UserName)
		}

		if nifcloud.ToString(ravgw.RemoteUserSet[0].Description) != "user1" {
			return fmt.Errorf("bad user.0.description state,  expected \"user1\", got: %#v", ravgw.RemoteUserSet[0].Description)
		}

		if nifcloud.ToString(ravgw.RemoteUserSet[1].UserName) != "user2" {
			return fmt.Errorf("bad user.1.name state,  expected \"user2\", got: %#v", ravgw.RemoteUserSet[1].UserName)
		}

		if nifcloud.ToString(ravgw.RemoteUserSet[1].Description) != "user2" {
			return fmt.Errorf("bad user.1.description state,  expected \"user2\", got: %#v", ravgw.RemoteUserSet[1].Description)
		}

		if nifcloud.ToString(ravgw.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad next_month_accounting_type state,  expected \"2\", got: %#v", ravgw.NextMonthAccountingType)
		}

		if nifcloud.ToString(ravgw.RemoteAccessVpnGatewayName) != rName {
			return fmt.Errorf("bad name,  expected \"%s\", got: %#v", rName, ravgw.RemoteAccessVpnGatewayName)
		}

		if nifcloud.ToString(ravgw.RemoteAccessVpnGatewayType) != "small" {
			return fmt.Errorf("bad type state,  expected \"small\", got: %#v", ravgw.RemoteAccessVpnGatewayType)
		}

		return nil
	}
}

func testAccCheckRemoteAccessVpnGatewayValuesUpdated(ravgw *types.RemoteAccessVpnGatewaySetOfDescribeRemoteAccessVpnGateways, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(ravgw.AccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state,  expected \"2\", got: %#v", ravgw.AccountingType)
		}

		if nifcloud.ToString(ravgw.Description) != "memo-upd" {
			return fmt.Errorf("bad description state,  expected \"memo-upd\", got: %#v", ravgw.Description)
		}

		if nifcloud.ToString(ravgw.NextMonthAccountingType) != "1" {
			return fmt.Errorf("bad next_month_accounting_type state,  expected \"1\", got: %#v", ravgw.NextMonthAccountingType)
		}

		if nifcloud.ToString(ravgw.RemoteAccessVpnGatewayName) != rName+"upd" {
			return fmt.Errorf("bad name,  expected \"%s\", got: %#v", rName+"upd", ravgw.RemoteAccessVpnGatewayName)
		}

		if nifcloud.ToString(ravgw.RemoteAccessVpnGatewayType) != "medium" {
			return fmt.Errorf("bad type state,  expected \"medium\", got: %#v", ravgw.RemoteAccessVpnGatewayType)
		}

		if nifcloud.ToString(ravgw.SslCertificateId) == "" {
			return fmt.Errorf("bad ssl_certificate_id state,  expected not empty.")
		}

		if len(ravgw.RemoteUserSet) != 1 {
			return fmt.Errorf("bad user length,  expected length 1, got %d", len(ravgw.RemoteUserSet))
		}

		if nifcloud.ToString(ravgw.RemoteUserSet[0].UserName) != "user3" {
			return fmt.Errorf("bad user.0.name state,  expected \"user3\", got: %#v", ravgw.RemoteUserSet[0].UserName)
		}

		if nifcloud.ToString(ravgw.RemoteUserSet[0].Description) != "user3" {
			return fmt.Errorf("bad user.3.description state,  expected \"user3\", got: %#v", ravgw.RemoteUserSet[0].Description)
		}

		return nil
	}
}

func testAccRemoteAccessVpnGatewayResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_remote_access_vpn_gateway" {
			continue
		}

		res, err := svc.DescribeRemoteAccessVpnGateways(context.Background(), &computing.DescribeRemoteAccessVpnGatewaysInput{
			RemoteAccessVpnGatewayId: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.RemoteAccessVpnGatewayId" {
				return nil
			}
			return fmt.Errorf("failed listing ravgws: %s", err)
		}

		if len(res.RemoteAccessVpnGatewaySet) > 0 {
			return fmt.Errorf("ravgw (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testSweepRemoteAccessVpnGateway(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeRemoteAccessVpnGateways(ctx, nil)
	if err != nil {
		return err
	}

	var sweepRemoteAccessVpnGateways []string
	for _, r := range res.RemoteAccessVpnGatewaySet {
		if strings.HasPrefix(nifcloud.ToString(r.RemoteAccessVpnGatewayName), prefix) {
			sweepRemoteAccessVpnGateways = append(sweepRemoteAccessVpnGateways, nifcloud.ToString(r.RemoteAccessVpnGatewayId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepRemoteAccessVpnGateways {
		ravgwID := n
		eg.Go(func() error {
			_, err = svc.DeleteRemoteAccessVpnGateway(ctx, &computing.DeleteRemoteAccessVpnGatewayInput{
				RemoteAccessVpnGatewayId: nifcloud.String(ravgwID),
			})
			if err != nil {
				return err
			}

			err = computing.NewRemoteAccessVpnGatewayDeletedWaiter(svc).Wait(ctx, &computing.DescribeRemoteAccessVpnGatewaysInput{
				RemoteAccessVpnGatewayId: []string{ravgwID},
			}, 600*time.Second)
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
