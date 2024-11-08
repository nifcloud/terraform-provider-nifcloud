package acc

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

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
	resource.AddTestSweepers("nifcloud_multi_ip_address_group", &resource.Sweeper{
		Name:         "nifcloud_multi_ip_address_group",
		F:            testSweepMultiIPAddressGroup,
		Dependencies: []string{},
	})
}

func TestAcc_MultiIPAddressGroup(t *testing.T) {
	var multiIPAddressGroup types.MultiIpAddressGroupsSet

	resourceName := "nifcloud_multi_ip_address_group.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccMultiIPAddressGroupResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMultiIPAddressGroup(t, "testdata/multi_ip_address_group.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiIPAddressGroupExists(resourceName, &multiIPAddressGroup),
					testAccCheckMultiIPAddressGroupValues(&multiIPAddressGroup, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "ip_address_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "default_gateway"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_mask"),
				),
			},
			{
				Config: testAccMultiIPAddressGroup(t, "testdata/multi_ip_address_group_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiIPAddressGroupExists(resourceName, &multiIPAddressGroup),
					testAccCheckMultiIPAddressGroupValuesUpdated(&multiIPAddressGroup, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "ip_address_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "ip_addresses.#", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "default_gateway"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_mask"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccMultiIPAddressGroup(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccCheckMultiIPAddressGroupExists(n string, multiIPAddressGroup *types.MultiIpAddressGroupsSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no MultiIPAddressGroup resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no MultiIPAddressGroup id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeMultiIpAddressGroups(
			context.Background(),
			&computing.DescribeMultiIpAddressGroupsInput{
				MultiIpAddressGroupId: []string{saved.Primary.ID},
			},
		)
		if err != nil {
			return err
		}

		if len(res.MultiIpAddressGroupsSet) == 0 {
			return fmt.Errorf("MultiIPAddressGroup does not found in cloud: %s", saved.Primary.ID)
		}

		foundMultiIPAddressGroup := res.MultiIpAddressGroupsSet[0]

		if nifcloud.ToString(foundMultiIPAddressGroup.MultiIpAddressGroupId) != saved.Primary.ID {
			return fmt.Errorf("MultiIPAddressGroup does not found in cloud: %s", saved.Primary.ID)
		}

		*multiIPAddressGroup = foundMultiIPAddressGroup

		return nil
	}
}

func testAccCheckMultiIPAddressGroupValues(multiIPAddressGroup *types.MultiIpAddressGroupsSet, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(multiIPAddressGroup.MultiIpAddressGroupName) != name {
			return fmt.Errorf("bad name state, expected %#v, got: %#v", name, multiIPAddressGroup.MultiIpAddressGroupName)
		}

		if nifcloud.ToString(multiIPAddressGroup.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", multiIPAddressGroup.Description)
		}

		if nifcloud.ToString(multiIPAddressGroup.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", multiIPAddressGroup.AvailabilityZone)
		}

		if len(multiIPAddressGroup.MultiIpAddressNetwork.IpAddressesSet) != 1 {
			return fmt.Errorf("bad ip_address_count state, expected count: 1, got: %d", len(multiIPAddressGroup.MultiIpAddressNetwork.IpAddressesSet))
		}

		return nil
	}
}

func testAccCheckMultiIPAddressGroupValuesUpdated(multiIPAddressGroup *types.MultiIpAddressGroupsSet, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(multiIPAddressGroup.MultiIpAddressGroupName) != name+"upd" {
			return fmt.Errorf("bad name state, expected %#v, got: %#v", name, multiIPAddressGroup.MultiIpAddressGroupName)
		}

		if nifcloud.ToString(multiIPAddressGroup.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", multiIPAddressGroup.Description)
		}

		if nifcloud.ToString(multiIPAddressGroup.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", multiIPAddressGroup.AvailabilityZone)
		}

		if len(multiIPAddressGroup.MultiIpAddressNetwork.IpAddressesSet) != 3 {
			return fmt.Errorf("bad ip_address_count state, expected count: 3, got: %d", len(multiIPAddressGroup.MultiIpAddressNetwork.IpAddressesSet))
		}

		return nil
	}
}

func testAccMultiIPAddressGroupResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_multi_ip_address_group" {
			continue
		}

		res, err := svc.DescribeMultiIpAddressGroups(
			context.Background(),
			&computing.DescribeMultiIpAddressGroupsInput{
				MultiIpAddressGroupId: []string{rs.Primary.ID},
			},
		)
		if err != nil {
			return fmt.Errorf("failed DescribeMultiIpAddressGroupsRequest: %s", err)
		}

		if len(res.MultiIpAddressGroupsSet) > 0 {
			return fmt.Errorf("MultiIpAddressGroup (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testSweepMultiIPAddressGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeMultiIpAddressGroups(ctx, nil)
	if err != nil {
		return err
	}

	var sweepMultiIPAddressGroups []string
	for _, g := range res.MultiIpAddressGroupsSet {
		if strings.HasPrefix(nifcloud.ToString(g.MultiIpAddressGroupName), prefix) {
			sweepMultiIPAddressGroups = append(sweepMultiIPAddressGroups, nifcloud.ToString(g.MultiIpAddressGroupId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, id := range sweepMultiIPAddressGroups {
		groupID := id
		eg.Go(func() error {
			_, err := svc.DeleteMultiIpAddressGroup(ctx, &computing.DeleteMultiIpAddressGroupInput{
				MultiIpAddressGroupId: nifcloud.String(groupID),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
