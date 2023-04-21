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
	resource.AddTestSweepers("nifcloud_vpn_gateway", &resource.Sweeper{
		Name: "nifcloud_vpn_gateway",
		F:    testSweepVpnGateway,
		Dependencies: []string{
			"nifcloud_vpn_connection",
		},
	})
}

func TestAcc_VpnGateway(t *testing.T) {
	var vpnGateway types.VpnGatewaySetOfDescribeVpnGateways

	resourceName := "nifcloud_vpn_gateway.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccVpnGatewayResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnGateway(t, "testdata/vpn_gateway.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists(resourceName, &vpnGateway),
					testAccCheckVpnGatewayValues(&vpnGateway, randName),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "small"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_name", randName),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "192.168.3.1"),
					resource.TestCheckResourceAttr(resourceName, "security_group", randName),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_association_id"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_address"),
				),
			},
			{
				Config: testAccVpnGateway(t, "testdata/vpn_gateway_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists(resourceName, &vpnGateway),
					testAccCheckVpnGatewayValuesUpdated(&vpnGateway, randName),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "type", "medium"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_name", randName),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "192.168.3.2"),
					resource.TestCheckResourceAttr(resourceName, "security_group", randName),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"network_name",
					"network_id",
					"route_table_id",
				},
			},
		},
	})
}

func testAccVpnGateway(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
	)
}

func testAccCheckVpnGatewayExists(n string, vpnGateway *types.VpnGatewaySetOfDescribeVpnGateways) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no vpn gateway resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no vpn gateway id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeVpnGateways(context.Background(), &computing.DescribeVpnGatewaysInput{
			VpnGatewayId: []string{saved.Primary.ID},
		})
		if err != nil {
			return err
		}

		if res == nil || len(res.VpnGatewaySet) == 0 {
			return fmt.Errorf("vpn gateway does not found in cloud: %s", saved.Primary.ID)
		}

		foundVpnGateway := res.VpnGatewaySet[0]

		if nifcloud.ToString(foundVpnGateway.VpnGatewayId) != saved.Primary.ID {
			return fmt.Errorf("vpn gateway does not found in cloud: %s", saved.Primary.ID)
		}

		*vpnGateway = foundVpnGateway

		return nil
	}
}

func testAccCheckVpnGatewayValues(vpnGateway *types.VpnGatewaySetOfDescribeVpnGateways, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if nifcloud.ToString(vpnGateway.NiftyVpnGatewayDescription) != "memo" {
			return fmt.Errorf("bad description state,  expected \"memo\", got: %#v", vpnGateway.NiftyVpnGatewayDescription)
		}

		if nifcloud.ToString(vpnGateway.NiftyVpnGatewayName) != rName {
			return fmt.Errorf("bad name,  expected \"%s\", got: %#v", rName, vpnGateway.NiftyVpnGatewayName)
		}

		if nifcloud.ToString(vpnGateway.NiftyVpnGatewayType) != "small" {
			return fmt.Errorf("bad type state,  expected \"small\", got: %#v", vpnGateway.NiftyVpnGatewayType)
		}

		if nifcloud.ToString(vpnGateway.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", vpnGateway.AvailabilityZone)
		}

		if nifcloud.ToString(vpnGateway.AccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state,  expected \"2\", got: %#v", vpnGateway.AccountingType)
		}

		if nifcloud.ToString(vpnGateway.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad next_month_accounting_type state,  expected \"2\", got: %#v", vpnGateway.NextMonthAccountingType)
		}

		// swap network interfaces.
		for i, ni := range vpnGateway.NetworkInterfaceSet {
			if nifcloud.ToString(ni.NetworkId) == "net-COMMON_GLOBAL" {
				if i == 0 {
					break
				}
				vpnGateway.NetworkInterfaceSet[0], vpnGateway.NetworkInterfaceSet[1] = vpnGateway.NetworkInterfaceSet[1], vpnGateway.NetworkInterfaceSet[0]
			}
		}

		if nifcloud.ToString(vpnGateway.NetworkInterfaceSet[0].NetworkId) != "net-COMMON_GLOBAL" {
			return fmt.Errorf("bad network_id state,  expected net-COMMON_GLOBAL, got: %#v", vpnGateway.NetworkInterfaceSet[0].NetworkId)
		}

		if nifcloud.ToString(vpnGateway.NetworkInterfaceSet[1].NetworkName) != rName {
			return fmt.Errorf("bad network_name state,  expected \"%s\", got: %#v", rName, vpnGateway.NetworkInterfaceSet[1].NetworkName)
		}

		if nifcloud.ToString(vpnGateway.NetworkInterfaceSet[1].IpAddress) != "192.168.3.1" {
			return fmt.Errorf("bad ip_address state,  expected \"%s\", got: %#v", "192.168.3.1", vpnGateway.NetworkInterfaceSet[1].IpAddress)
		}

		if nifcloud.ToString(vpnGateway.GroupSet[0].GroupId) != rName {
			return fmt.Errorf("bad group_id state,  expected \"%s\", got: %#v", rName, vpnGateway.GroupSet[0].GroupId)
		}

		if nifcloud.ToString(vpnGateway.RouteTableAssociationId) == "" {
			return fmt.Errorf("bad route_table_association_id,  expected not empty string, got: %#v", vpnGateway.RouteTableAssociationId)
		}

		if nifcloud.ToString(vpnGateway.RouteTableId) == "" {
			return fmt.Errorf("bad route_table_id,  expected not empty string, got: %#v", vpnGateway.RouteTableId)
		}

		return nil
	}
}

func testAccCheckVpnGatewayValuesUpdated(vpnGateway *types.VpnGatewaySetOfDescribeVpnGateways, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(vpnGateway.NiftyVpnGatewayDescription) != "memo-upd" {
			return fmt.Errorf("bad description state,  expected \"memo-upd\", got: %#v", vpnGateway.NiftyVpnGatewayDescription)
		}

		if nifcloud.ToString(vpnGateway.NiftyVpnGatewayName) != rName+"upd" {
			return fmt.Errorf("bad name,  expected \"%s\", got: %#v", rName+"upd", vpnGateway.NiftyVpnGatewayName)
		}

		if nifcloud.ToString(vpnGateway.NiftyVpnGatewayType) != "medium" {
			return fmt.Errorf("bad type state,  expected \"medium\", got: %#v", vpnGateway.NiftyVpnGatewayType)
		}

		if nifcloud.ToString(vpnGateway.AccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state,  expected \"2\", got: %#v", vpnGateway.AccountingType)
		}

		if nifcloud.ToString(vpnGateway.NextMonthAccountingType) != "1" {
			return fmt.Errorf("bad next_month_accounting_type state,  expected \"1\", got: %#v", vpnGateway.NextMonthAccountingType)
		}

		if nifcloud.ToString(vpnGateway.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", vpnGateway.AvailabilityZone)
		}

		if len(vpnGateway.NetworkInterfaceSet) != 2 {
			return fmt.Errorf("bad network_interface length,  expected length 2, got %d", len(vpnGateway.NetworkInterfaceSet))
		}

		// swap network interfaces.
		for i, ni := range vpnGateway.NetworkInterfaceSet {
			if nifcloud.ToString(ni.NetworkId) == "net-COMMON_GLOBAL" {
				if i == 0 {
					break
				}
				vpnGateway.NetworkInterfaceSet[0], vpnGateway.NetworkInterfaceSet[1] = vpnGateway.NetworkInterfaceSet[1], vpnGateway.NetworkInterfaceSet[0]
			}
		}

		if nifcloud.ToString(vpnGateway.NetworkInterfaceSet[0].NetworkId) != "net-COMMON_GLOBAL" {
			return fmt.Errorf("bad network_interface.0.network_id state,  expected net-COMMON_GLOBAL, got: %#v", vpnGateway.NetworkInterfaceSet[0].NetworkId)
		}

		if nifcloud.ToString(vpnGateway.NetworkInterfaceSet[1].NetworkName) != rName {
			return fmt.Errorf("bad network_name state,  expected \"%s\", got: %#v", rName, vpnGateway.NetworkInterfaceSet[1].NetworkName)
		}

		if nifcloud.ToString(vpnGateway.NetworkInterfaceSet[1].IpAddress) != "192.168.3.2" {
			return fmt.Errorf("bad ip_address state,  expected \"%s\", got: %#v", "192.168.3.2", vpnGateway.NetworkInterfaceSet[1].IpAddress)
		}

		if nifcloud.ToString(vpnGateway.GroupSet[0].GroupId) != rName {
			return fmt.Errorf("bad group_id state,  expected \"%s\", got: %#v", rName, vpnGateway.GroupSet[0].GroupId)
		}

		if nifcloud.ToString(vpnGateway.RouteTableAssociationId) != "" {
			return fmt.Errorf("bad route_table_association_id,  expected empty string, got: empty string")
		}

		if nifcloud.ToString(vpnGateway.RouteTableId) != "" {
			return fmt.Errorf("bad route_table_id,  expected empty string, got: empty string")
		}
		return nil
	}
}

func testAccVpnGatewayResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_vpn_gateway" {
			continue
		}

		res, err := svc.DescribeVpnGateways(context.Background(), &computing.DescribeVpnGatewaysInput{
			VpnGatewayId: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.VpnGatewayId" {
				return nil
			}
			return fmt.Errorf("failed listing vpn gateways: %s", err)
		}

		if len(res.VpnGatewaySet) > 0 {
			return fmt.Errorf("vpn gateway (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testSweepVpnGateway(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeVpnGateways(ctx, nil)
	if err != nil {
		return err
	}

	var sweepVpnGateways []string
	for _, r := range res.VpnGatewaySet {
		if strings.HasPrefix(nifcloud.ToString(r.NiftyVpnGatewayName), prefix) {
			sweepVpnGateways = append(sweepVpnGateways, nifcloud.ToString(r.VpnGatewayId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepVpnGateways {
		vpnGatewayID := n
		eg.Go(func() error {
			_, err = svc.DeleteVpnGateway(ctx, &computing.DeleteVpnGatewayInput{
				VpnGatewayId: nifcloud.String(vpnGatewayID),
			})
			if err != nil {
				return err
			}

			err = computing.NewVpnGatewayDeletedWaiter(svc).Wait(ctx, &computing.DescribeVpnGatewaysInput{
				VpnGatewayId: []string{vpnGatewayID},
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
