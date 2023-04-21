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
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_private_lan", &resource.Sweeper{
		Name: "nifcloud_private_lan",
		F:    testSweepPrivateLan,
		Dependencies: []string{
			"nifcloud_router",
		},
	})
}

func TestAcc_PrivateLan(t *testing.T) {
	var privateLan types.PrivateLanSet

	resourceName := "nifcloud_private_lan.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccPrivateLanResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateLan(t, "testdata/private_lan.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateLanExists(resourceName, &privateLan),
					testAccCheckPrivateLanValues(&privateLan, randName),
					resource.TestCheckResourceAttr(resourceName, "private_lan_name", randName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
				),
			},
			{
				Config: testAccPrivateLan(t, "testdata/private_lan_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateLanExists(resourceName, &privateLan),
					testAccCheckPrivateLanValuesUpdated(&privateLan, randName),
					resource.TestCheckResourceAttr(resourceName, "private_lan_name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
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

func testAccPrivateLan(t *testing.T, fileName string, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccPrivateLanResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_private_lan" {
			continue
		}

		res, err := svc.NiftyDescribePrivateLans(context.Background(), &computing.NiftyDescribePrivateLansInput{
			NetworkId: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.NetworkId" {
				return nil
			}
			return fmt.Errorf("failed NiftyDescribePrivateLansRequest: %s", err)
		}

		if len(res.PrivateLanSet) > 0 {
			return fmt.Errorf("privateLan (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckPrivateLanValues(privateLan *types.PrivateLanSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(privateLan.PrivateLanName) != rName {
			return fmt.Errorf("bad name state, expected %#v, got: %#v", rName, privateLan.PrivateLanName)
		}

		if nifcloud.ToString(privateLan.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", privateLan.Description)
		}

		if nifcloud.ToString(privateLan.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", privateLan.AvailabilityZone)
		}

		if nifcloud.ToString(privateLan.CidrBlock) != "192.168.1.0/24" {
			return fmt.Errorf("bad cidr_block state,  expected \"192.168.1.0/24\", got: %#v", privateLan.CidrBlock)
		}

		if nifcloud.ToString(privateLan.NextMonthAccountingType) != "1" {
			return fmt.Errorf("bad accounting_type state, expected \"1\", got: %#v", privateLan.NextMonthAccountingType)
		}

		return nil
	}
}

func testAccCheckPrivateLanValuesUpdated(privateLan *types.PrivateLanSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(privateLan.PrivateLanName) != rName+"upd" {
			return fmt.Errorf("bad name state, expected %#v, got: %#v", rName, privateLan.PrivateLanName)
		}

		if nifcloud.ToString(privateLan.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", privateLan.Description)
		}

		if nifcloud.ToString(privateLan.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", privateLan.AvailabilityZone)
		}

		if nifcloud.ToString(privateLan.CidrBlock) != "192.168.2.0/24" {
			return fmt.Errorf("bad cidr_block state,  expected \"192.168.2.0/24\", got: %#v", privateLan.CidrBlock)
		}

		if nifcloud.ToString(privateLan.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state, expected \"1\", got: %#v", privateLan.NextMonthAccountingType)
		}

		return nil
	}
}

func testSweepPrivateLan(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.NiftyDescribePrivateLans(ctx, nil)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range res.PrivateLanSet {
		privateLan := n
		if strings.HasPrefix(nifcloud.ToString(privateLan.Description), prefix) {
			eg.Go(func() error {
				_, err := svc.NiftyDeletePrivateLan(ctx, &computing.NiftyDeletePrivateLanInput{
					NetworkId: privateLan.NetworkId,
				})
				return err
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func testAccCheckPrivateLanExists(n string, privateLan *types.PrivateLanSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no privateLan resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no privateLan id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.NiftyDescribePrivateLans(context.Background(), &computing.NiftyDescribePrivateLansInput{
			NetworkId: []string{saved.Primary.ID},
		})

		if err != nil {
			return err
		}

		if len(res.PrivateLanSet) == 0 {
			return fmt.Errorf("privateLan does not found in cloud: %s", saved.Primary.ID)
		}

		foundPrivateLan := res.PrivateLanSet[0]

		if nifcloud.ToString(foundPrivateLan.NetworkId) != saved.Primary.ID {
			return fmt.Errorf("privateLan does not found in cloud: %s", saved.Primary.ID)
		}

		*privateLan = foundPrivateLan
		return nil
	}
}
