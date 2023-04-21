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
	resource.AddTestSweepers("nifcloud_network_interface", &resource.Sweeper{
		Name: "nifcloud_network_interface",
		F:    testSweepNetworkInterface,
		Dependencies: []string{
			"nifcloud_instance",
		},
	})
}

func TestAcc_NetworkInterface(t *testing.T) {
	var networkInterface types.NetworkInterfaceSetOfDescribeNetworkInterfaces

	resourceName := "nifcloud_network_interface.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccNetworkInterfaceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterface(t, "testdata/network_interface.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
					testAccCheckNetworkInterfaceValues(&networkInterface, randName),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "static"),
					resource.TestCheckResourceAttr(resourceName, "private_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "description", randName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "network_id"),
				),
			},
			{
				Config: testAccNetworkInterface(t, "testdata/network_interface_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
					testAccCheckNetworkInterfaceValuesUpdated(&networkInterface, randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "192.168.100.100"),
					resource.TestCheckResourceAttr(resourceName, "private_ip", "192.168.100.100"),
					resource.TestCheckResourceAttr(resourceName, "description", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "network_id"),
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

func testAccNetworkInterface(t *testing.T, fileName, rName string) string {
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

func testAccCheckNetworkInterfaceExists(n string, networkInterface *types.NetworkInterfaceSetOfDescribeNetworkInterfaces) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no networkInterface resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no networkInterface id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeNetworkInterfaces(context.Background(), &computing.DescribeNetworkInterfacesInput{
			NetworkInterfaceId: []string{saved.Primary.ID},
		})

		if err != nil {
			return err
		}

		if len(res.NetworkInterfaceSet) == 0 {
			return fmt.Errorf("networkInterface does not found in cloud: %s", saved.Primary.ID)
		}

		foundNetworkInterface := res.NetworkInterfaceSet[0]

		if nifcloud.ToString(foundNetworkInterface.NetworkInterfaceId) != saved.Primary.ID {
			return fmt.Errorf("networkInterface does not found in cloud: %s", saved.Primary.ID)
		}

		*networkInterface = foundNetworkInterface
		return nil
	}
}

func testAccCheckNetworkInterfaceValues(networkInterface *types.NetworkInterfaceSetOfDescribeNetworkInterfaces, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(networkInterface.Description) != rName {
			return fmt.Errorf("bad description state, expected \"%s\", got: %#v", rName, networkInterface.Description)
		}

		if nifcloud.ToString(networkInterface.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", networkInterface.AvailabilityZone)
		}

		if nifcloud.ToString(networkInterface.PrivateIpAddress) != "" {
			return fmt.Errorf("bad private_ip state,  expected \"empty\", got: %#v", networkInterface.PrivateIpAddress)
		}

		if nifcloud.ToString(networkInterface.NetworkInterfaceId) == "" {
			return fmt.Errorf("bad network_interface_id state,  expected \"not empty\", got: %#v", networkInterface.NetworkInterfaceId)
		}

		if nifcloud.ToString(networkInterface.NiftyNetworkId) == "" {
			return fmt.Errorf("bad network_id state,  expected \"not empty\", got: %#v", networkInterface.NiftyNetworkId)
		}
		return nil
	}
}

func testAccCheckNetworkInterfaceValuesUpdated(networkInterface *types.NetworkInterfaceSetOfDescribeNetworkInterfaces, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(networkInterface.Description) != rName {
			return fmt.Errorf("bad description state, expected \"%s\", got: %#v", rName, networkInterface.Description)
		}

		if nifcloud.ToString(networkInterface.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", networkInterface.AvailabilityZone)
		}

		if nifcloud.ToString(networkInterface.PrivateIpAddress) != "192.168.100.100" {
			return fmt.Errorf("bad private_ip state,  expected \"192.168.100.100\", got: %#v", networkInterface.PrivateIpAddress)
		}

		if nifcloud.ToString(networkInterface.NetworkInterfaceId) == "" {
			return fmt.Errorf("bad network_interface_id state,  expected \"not empty\", got: %#v", networkInterface.NetworkInterfaceId)
		}

		if nifcloud.ToString(networkInterface.NiftyNetworkId) == "" {
			return fmt.Errorf("bad network_id state,  expected \"not empty\", got: %#v", networkInterface.NiftyNetworkId)
		}
		return nil
	}
}

func testAccNetworkInterfaceResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_network_interface" {
			continue
		}

		res, err := svc.DescribeNetworkInterfaces(context.Background(), &computing.DescribeNetworkInterfacesInput{
			NetworkInterfaceId: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.NetworkInterfaceId" {
				return nil
			}
			return fmt.Errorf("failed DescribeNetworkInterfacesRequest: %s", err)
		}

		if len(res.NetworkInterfaceSet) > 0 {
			return fmt.Errorf("networkInterface (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepNetworkInterface(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeNetworkInterfaces(ctx, nil)
	if err != nil {
		return err
	}

	var sweepNetworkInterfaces []string
	for _, k := range res.NetworkInterfaceSet {
		if strings.HasPrefix(nifcloud.ToString(k.Description), prefix) {
			sweepNetworkInterfaces = append(sweepNetworkInterfaces, nifcloud.ToString(k.NetworkInterfaceId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepNetworkInterfaces {
		id := n
		eg.Go(func() error {
			_, err := svc.DeleteNetworkInterface(ctx, &computing.DeleteNetworkInterfaceInput{
				NetworkInterfaceId: nifcloud.String(id),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
