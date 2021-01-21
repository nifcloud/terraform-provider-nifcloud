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
	resource.AddTestSweepers("nifcloud_router", &resource.Sweeper{
		Name: "nifcloud_router",
		F:    testSweepRouter,
	})
}

func TestAcc_Router(t *testing.T) {
	var router computing.RouterSetOfNiftyDescribeRouters

	resourceName := "nifcloud_router.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccRouterResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouter(t, "testdata/router.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterExists(resourceName, &router),
					testAccCheckRouterValues(&router, randName),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.network_id"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.ip_address", "192.168.1.1"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.dhcp", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.dhcp_options_id"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.dhcp_config_id"),
					resource.TestCheckResourceAttr(resourceName, "security_group", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "small"),
					resource.TestCheckResourceAttrSet(resourceName, "router_id"),
				),
			},
			{
				Config: testAccRouter(t, "testdata/router_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterExists(resourceName, &router),
					testAccCheckRouterValuesUpdated(&router, randName),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"upd"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.network_id"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.1.network_id"),
					resource.TestCheckResourceAttr(resourceName, "security_group", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "small"),
					resource.TestCheckResourceAttrSet(resourceName, "router_id"),
					resource.TestCheckResourceAttrSet(resourceName, "nat_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "nat_table_association_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_association_id"),
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

func testAccRouter(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
	)
}

func testAccCheckRouterExists(n string, router *computing.RouterSetOfNiftyDescribeRouters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no router resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no router id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.NiftyDescribeRoutersRequest(&computing.NiftyDescribeRoutersInput{
			RouterId: []string{saved.Primary.ID},
		}).Send(context.Background())
		if err != nil {
			return err
		}

		if res == nil || len(res.RouterSet) == 0 {
			return fmt.Errorf("router does not found in cloud: %s", saved.Primary.ID)
		}

		foundRouter := res.RouterSet[0]

		if nifcloud.StringValue(foundRouter.RouterId) != saved.Primary.ID {
			return fmt.Errorf("router does not found in cloud: %s", saved.Primary.ID)
		}

		*router = foundRouter

		return nil
	}
}

func testAccCheckRouterValues(router *computing.RouterSetOfNiftyDescribeRouters, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(router.AccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state,  expected \"2\", got: %#v", router.NextMonthAccountingType)
		}

		if nifcloud.StringValue(router.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", router.AvailabilityZone)
		}

		if nifcloud.StringValue(router.Description) != "memo" {
			return fmt.Errorf("bad description state,  expected \"memo\", got: %#v", router.Description)
		}

		if nifcloud.StringValue(router.GroupSet[0].GroupId) != rName {
			return fmt.Errorf("bad group_id state,  expected \"%s\", got: %#v", rName, router.GroupSet[0].GroupId)
		}

		if nifcloud.StringValue(router.NatTableAssociationId) != "" {
			return fmt.Errorf("bad nat_table_association_id,  expected empty string, got: %#v", router.NatTableAssociationId)
		}

		if nifcloud.StringValue(router.NatTableId) != "" {
			return fmt.Errorf("bad nat_table_id,  expected empty string, got: %#v", router.NatTableId)
		}

		if len(router.NetworkInterfaceSet) != 1 {
			return fmt.Errorf("bad network_interface length,  expected length 1, got %d", len(router.NetworkInterfaceSet))
		}

		if nifcloud.StringValue(router.NetworkInterfaceSet[0].NetworkName) != rName {
			return fmt.Errorf("bad network_interface.0.network_name state,  expected %s, got: %#v", rName, router.NetworkInterfaceSet[0].NetworkName)
		}

		if nifcloud.StringValue(router.NetworkInterfaceSet[0].IpAddress) != "192.168.1.1" {
			return fmt.Errorf("bad network_interface.0.ip_address state,  expected \"192.168.1.1\", got: %#v", router.NetworkInterfaceSet[0].IpAddress)
		}

		if !nifcloud.BoolValue(router.NetworkInterfaceSet[0].Dhcp) {
			return fmt.Errorf("bad network_interface.0.dhcp state,  expected true, got: false")
		}

		if nifcloud.StringValue(router.NetworkInterfaceSet[0].DhcpConfigId) == "" {
			return fmt.Errorf("bad network_interface.0.dhcp_config_id state,  expected not empty string, got: empty string")
		}

		if nifcloud.StringValue(router.NetworkInterfaceSet[0].DhcpOptionsId) == "" {
			return fmt.Errorf("bad network_interface.0.dhcp_options_id state,  expected not empty string, got: empty string")
		}

		if nifcloud.StringValue(router.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad next_month_accounting_type state,  expected \"2\", got: %#v", router.NextMonthAccountingType)
		}

		if nifcloud.StringValue(router.RouteTableAssociationId) != "" {
			return fmt.Errorf("bad route_table_association_id,  expected empty string, got: %#v", router.RouteTableAssociationId)
		}

		if nifcloud.StringValue(router.RouteTableId) != "" {
			return fmt.Errorf("bad route_table_id,  expected empty string, got: %#v", router.RouteTableId)
		}

		if nifcloud.StringValue(router.RouterName) != rName {
			return fmt.Errorf("bad name,  expected \"%s\", got: %#v", rName, router.RouterName)
		}

		if nifcloud.StringValue(router.Type) != "small" {
			return fmt.Errorf("bad type state,  expected \"small\", got: %#v", router.Type)
		}

		return nil
	}
}

func testAccCheckRouterValuesUpdated(router *computing.RouterSetOfNiftyDescribeRouters, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(router.AccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state,  expected \"2\", got: %#v", router.NextMonthAccountingType)
		}

		if nifcloud.StringValue(router.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", router.AvailabilityZone)
		}

		if nifcloud.StringValue(router.Description) != "memo-upd" {
			return fmt.Errorf("bad description state,  expected \"memo-upd\", got: %#v", router.Description)
		}

		if nifcloud.StringValue(router.GroupSet[0].GroupId) != rName {
			return fmt.Errorf("bad group_id state,  expected \"%s\", got: %#v", rName, router.GroupSet[0].GroupId)
		}

		if nifcloud.StringValue(router.NatTableAssociationId) == "" {
			return fmt.Errorf("bad nat_table_association_id,  expected not empty string, got: empty string")
		}

		if nifcloud.StringValue(router.NatTableId) == "" {
			return fmt.Errorf("bad nat_table_id,  expected not empty string, got: empty string")
		}

		if len(router.NetworkInterfaceSet) != 2 {
			return fmt.Errorf("bad network_interface length,  expected length 2, got %d", len(router.NetworkInterfaceSet))
		}

		// swap network interfaces.
		for i, ni := range router.NetworkInterfaceSet {
			if nifcloud.StringValue(ni.NetworkId) == "net-COMMON_GLOBAL" {
				if i == 0 {
					break
				}
				router.NetworkInterfaceSet[0], router.NetworkInterfaceSet[1] = router.NetworkInterfaceSet[1], router.NetworkInterfaceSet[0]
			}
		}

		if nifcloud.StringValue(router.NetworkInterfaceSet[0].NetworkId) != "net-COMMON_GLOBAL" {
			return fmt.Errorf("bad network_interface.0.network_id state,  expected net-COMMON_GLOBAL, got: %#v", router.NetworkInterfaceSet[0].NetworkId)
		}

		if nifcloud.StringValue(router.NetworkInterfaceSet[1].NetworkName) != rName {
			return fmt.Errorf("bad network_interface.1.network_name state,  expected %s, got: %#v", rName, router.NetworkInterfaceSet[1].NetworkName)
		}

		if nifcloud.StringValue(router.NetworkInterfaceSet[1].IpAddress) != "192.168.1.1" {
			return fmt.Errorf("bad network_interface.1.ip_address state,  expected \"192.168.1.1\", got: %#v", router.NetworkInterfaceSet[1].IpAddress)
		}

		if !nifcloud.BoolValue(router.NetworkInterfaceSet[1].Dhcp) {
			return fmt.Errorf("bad network_interface.1.dhcp state,  expected true, got: false")
		}

		if nifcloud.StringValue(router.NetworkInterfaceSet[1].DhcpConfigId) == "" {
			return fmt.Errorf("bad network_interface.1.dhcp_config_id state,  expected not empty string, got: empty string")
		}

		if nifcloud.StringValue(router.NetworkInterfaceSet[1].DhcpOptionsId) == "" {
			return fmt.Errorf("bad network_interface.1.dhcp_options_id state,  expected not empty string, got: empty string")
		}

		if nifcloud.StringValue(router.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad next_month_accounting_type state,  expected \"2\", got: %#v", router.NextMonthAccountingType)
		}

		if nifcloud.StringValue(router.RouteTableAssociationId) == "" {
			return fmt.Errorf("bad route_table_association_id,  expected not empty string, got: empty string")
		}

		if nifcloud.StringValue(router.RouteTableId) == "" {
			return fmt.Errorf("bad route_table_id,  expected not empty string, got: empty string")
		}

		if nifcloud.StringValue(router.RouterName) != rName+"upd" {
			return fmt.Errorf("bad name,  expected \"%s\", got: %#v", rName+"upd", router.RouterName)
		}

		if nifcloud.StringValue(router.Type) != "small" {
			return fmt.Errorf("bad type state,  expected \"small\", got: %#v", router.Type)
		}

		return nil
	}
}

func testAccRouterResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_router" {
			continue
		}

		res, err := svc.NiftyDescribeRoutersRequest(&computing.NiftyDescribeRoutersInput{
			RouterId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.RouterId" {
				return fmt.Errorf("failed listing routers: %s", err)
			}
		}

		if len(res.RouterSet) > 0 {
			return fmt.Errorf("router (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testSweepRouter(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.NiftyDescribeRoutersRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepRouters []string
	for _, r := range res.RouterSet {
		if strings.HasPrefix(nifcloud.StringValue(r.RouterName), prefix) {
			sweepRouters = append(sweepRouters, nifcloud.StringValue(r.RouterId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepRouters {
		routerID := n
		eg.Go(func() error {
			_, err = svc.NiftyDeleteRouterRequest(&computing.NiftyDeleteRouterInput{
				RouterId: nifcloud.String(routerID),
			}).Send(ctx)
			if err != nil {
				return err
			}

			err = svc.WaitUntilRouterDeleted(ctx, &computing.NiftyDescribeRoutersInput{
				RouterId: []string{routerID},
			})
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
