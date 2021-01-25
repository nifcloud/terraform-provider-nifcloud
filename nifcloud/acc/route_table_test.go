package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func init() {
	resource.AddTestSweepers("nifcloud_route_table", &resource.Sweeper{
		Name: "nifcloud_route_table",
		F:    testSweepRouteTable,
		Dependencies: []string{
			"nifcloud_elb",
			"nifcloud_router",
		},
	})
}

func TestAcc_RouteTable(t *testing.T) {
	var routeTable computing.RouteTableSet

	resourceName := "nifcloud_route_table.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccRouteTableResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTable(t, "testdata/route_table.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(resourceName, &routeTable),
					testAccCheckRouteTableValues(&routeTable, randName),
					resource.TestCheckResourceAttr(resourceName, "route.0.cidr_block", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "route.0.network_id"),
					resource.TestCheckResourceAttr(resourceName, "route.1.cidr_block", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route.1.ip_address", "1.1.1.1"),
				),
			},
			{
				Config: testAccRouteTable(t, "testdata/route_table_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(resourceName, &routeTable),
					testAccCheckRouteTableValuesUpdated(&routeTable, randName),
					resource.TestCheckResourceAttr(resourceName, "route.0.cidr_block", "192.168.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route.0.network_name", randName),
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

func testAccRouteTable(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccCheckRouteTableExists(n string, routeTable *computing.RouteTableSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no routeTable resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no routeTable id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeRouteTablesRequest(&computing.DescribeRouteTablesInput{
			RouteTableId: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.RouteTableSet) == 0 {
			return fmt.Errorf("routeTable does not found in cloud: %s", saved.Primary.ID)
		}

		foundRouteTable := res.RouteTableSet[0]

		if nifcloud.StringValue(foundRouteTable.RouteTableId) != saved.Primary.ID {
			return fmt.Errorf("routeTable does not found in cloud: %s", saved.Primary.ID)
		}

		*routeTable = foundRouteTable
		return nil
	}
}

func testAccCheckRouteTableValues(routeTable *computing.RouteTableSet, privateLanName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(routeTable.RouteSet) != 2 {
			return fmt.Errorf("bad routes: %#v", routeTable.RouteSet)
		}

		routes := make(map[string]computing.RouteSet)
		for _, r := range routeTable.RouteSet {
			routes[*r.DestinationCidrBlock] = r
		}

		if _, ok := routes["192.168.1.0/24"]; !ok {
			return fmt.Errorf("bad routes: %#v", routeTable.RouteSet)
		}

		if _, ok := routes["192.168.2.0/24"]; !ok {
			return fmt.Errorf("bad routes: %#v", routeTable.RouteSet)
		}

		if nifcloud.StringValue(routes["192.168.1.0/24"].NetworkId) == "" {
			return fmt.Errorf("bad routes network id, expected \"not null\", got: %#v", nifcloud.StringValue(routeTable.RouteSet[0].NetworkId))
		}

		if nifcloud.StringValue(routes["192.168.2.0/24"].IpAddress) != "1.1.1.1" {
			return fmt.Errorf("bad routes ipaddress, expected \"1.1.1.1\", got: %#v", nifcloud.StringValue(routeTable.RouteSet[0].IpAddress))
		}
		return nil
	}
}

func testAccCheckRouteTableValuesUpdated(routeTable *computing.RouteTableSet, privateLanName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(routeTable.RouteSet) != 1 {
			return fmt.Errorf("bad routes: %#v", routeTable.RouteSet)
		}

		routes := make(map[string]computing.RouteSet)
		for _, r := range routeTable.RouteSet {
			routes[*r.DestinationCidrBlock] = r
		}

		if _, ok := routes["192.168.3.0/24"]; !ok {
			return fmt.Errorf("bad routes: %#v", routeTable.RouteSet)
		}

		if nifcloud.StringValue(routes["192.168.3.0/24"].NetworkName) != privateLanName {
			return fmt.Errorf("bad routes network name, expected \"%s\", got: %#v", privateLanName, nifcloud.StringValue(routeTable.RouteSet[0].NetworkId))
		}
		return nil
	}
}

func testAccRouteTableResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_route_table" {
			continue
		}

		res, err := svc.DescribeRouteTablesRequest(&computing.DescribeRouteTablesInput{
			RouteTableId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.RouteTableId" {
				return fmt.Errorf("failed DescribeRouteTablesRequest: %s", err)
			}
		}

		if len(res.RouteTableSet) > 0 {
			return fmt.Errorf("routeTable (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepRouteTable(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeRouteTablesRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	for _, routeTable := range res.RouteTableSet {
		isMainRouteTableAssociation := false

		for _, routeTableAssociation := range routeTable.AssociationSet {
			if nifcloud.BoolValue(routeTableAssociation.Main) {
				isMainRouteTableAssociation = true
				break
			}

			input := &computing.DisassociateRouteTableInput{
				AssociationId: routeTableAssociation.RouteTableAssociationId,
			}

			_, err := svc.DisassociateRouteTableRequest(input).Send(ctx)
			if err != nil {
				return err
			}
		}

		if isMainRouteTableAssociation {
			continue
		}

		input := &computing.DeleteRouteTableInput{
			RouteTableId: routeTable.RouteTableId,
		}

		_, err := svc.DeleteRouteTableRequest(input).Send(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
