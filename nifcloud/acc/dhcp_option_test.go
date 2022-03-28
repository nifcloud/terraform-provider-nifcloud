package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_dhcp_option", &resource.Sweeper{
		Name: "nifcloud_dhcp_option",
		F:    testSweepDhcpOption,
		Dependencies: []string{
			"nifcloud_router",
		},
	})
}

func TestAcc_DhcpOption(t *testing.T) {
	var dhcpOption types.DhcpOptionsSet

	resourceName := "nifcloud_dhcp_option.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDhcpOptionResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpOption(t, "testdata/dhcp_option.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDhcpOptionExists(resourceName, &dhcpOption),
					testAccCheckDhcpOptionValues(&dhcpOption),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_option_id"),
					resource.TestCheckResourceAttr(resourceName, "default_router", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.0", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers.1", "192.168.0.2"),
					resource.TestCheckResourceAttr(resourceName, "ntp_servers.0", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "netbios_name_servers.0", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "netbios_name_servers.1", "192.168.0.2"),
					resource.TestCheckResourceAttr(resourceName, "netbios_node_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "lease_time", "600"),
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

func testAccDhcpOption(t *testing.T, fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func testAccCheckDhcpOptionExists(n string, dhcpOption *types.DhcpOptionsSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no dhcpOption resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no dhcpOption id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeDhcpOptions(context.Background(), &computing.DescribeDhcpOptionsInput{
			DhcpOptionsId: []string{saved.Primary.ID},
		})

		if err != nil {
			return err
		}

		if res == nil || len(res.DhcpOptionsSet) == 0 || len(res.DhcpOptionsSet[0].DhcpConfigurationSet) == 0 {
			return fmt.Errorf("dhcpOption does not found in cloud: %s", saved.Primary.ID)
		}

		foundDhcpOption := res.DhcpOptionsSet[0]

		if nifcloud.ToString(foundDhcpOption.DhcpOptionsId) != saved.Primary.ID {
			return fmt.Errorf("dhcpOption does not found in cloud: %s", saved.Primary.ID)
		}

		*dhcpOption = foundDhcpOption
		return nil
	}
}

func testAccCheckDhcpOptionValues(dhcpOption *types.DhcpOptionsSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(dhcpOption.DhcpOptionsId) == "" {
			return fmt.Errorf("bad dhcp_option_id state, expected not nil, got: nil")
		}

		if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].Key) == "default-router" {
			if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "192.168.0.1" {
				return fmt.Errorf("bad default_router state, expected \"192.168.0.1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].Key) == "domain-name" {
			if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "example.com" {
				return fmt.Errorf("bad domain_name state, expected \"example.com\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].Key) == "domain-name-server" {
			if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "192.168.0.1" {
				return fmt.Errorf("bad domain_name_servers state, expected \"192.168.0.1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].Key) == "domain-name-server" {
			if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].ValueSet[1].Value) != "192.168.0.2" {
				return fmt.Errorf("bad domain_name_servers state, expected \"192.168.0.2\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[1].Value)
			}
		}

		if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].Key) == "ntp_servers" {
			if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "192.168.0.1" {
				return fmt.Errorf("bad ntp_servers state, expected \"192.168.0.1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].Key) == "netbios_name_servers" {
			if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "192.168.0.1" {
				return fmt.Errorf("bad netbios_name_servers state, expected \"192.168.0.1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].Key) == "netbios_name_servers" {
			if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].ValueSet[1].Value) != "192.168.0.2" {
				return fmt.Errorf("bad netbios_name_servers state, expected \"192.168.0.2\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[1].Value)
			}
		}

		if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].Key) == "netbios_node_type" {
			if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "1" {
				return fmt.Errorf("bad netbios_node_type state, expected \"1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].Key) == "lease_time" {
			if nifcloud.ToString(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "600" {
				return fmt.Errorf("bad lease_time state, expected \"600\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}
		return nil
	}
}

func testAccDhcpOptionResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_dhcp_option" {
			continue
		}

		res, err := svc.DescribeDhcpOptions(context.Background(), &computing.DescribeDhcpOptionsInput{
			DhcpOptionsId: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.DhcpOptionsId" {
				return nil
			}
			return fmt.Errorf("failed DescribeDhcpOptionsRequest: %s", err)
		}

		if len(res.DhcpOptionsSet) > 0 {
			return fmt.Errorf("dhcpOption (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepDhcpOption(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeDhcpOptions(ctx, nil)
	if err != nil {
		return err
	}

	var sweepDhcpOptions []string
	for _, k := range res.DhcpOptionsSet {
		if strings.HasPrefix(nifcloud.ToString(k.DhcpOptionsId), prefix) {
			sweepDhcpOptions = append(sweepDhcpOptions, nifcloud.ToString(k.DhcpOptionsId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepDhcpOptions {
		dhcpOptionID := n
		eg.Go(func() error {
			_, err := svc.DeleteDhcpOptions(ctx, &computing.DeleteDhcpOptionsInput{
				DhcpOptionsId: nifcloud.String(dhcpOptionID),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
