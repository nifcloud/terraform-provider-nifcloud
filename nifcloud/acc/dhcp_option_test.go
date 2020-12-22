package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_dhcp_option", &resource.Sweeper{
		Name: "nifcloud_dhcp_option",
		F:    testSweepDhcpOption,
	})
}

func TestAcc_DhcpOption(t *testing.T) {
	var dhcpOption computing.DhcpOptionsSet

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

func testAccCheckDhcpOptionExists(n string, dhcpOption *computing.DhcpOptionsSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no dhcpOption resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no dhcpOption id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeDhcpOptionsRequest(&computing.DescribeDhcpOptionsInput{
			DhcpOptionsId: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if res == nil || len(res.DhcpOptionsSet) == 0 || len(res.DhcpOptionsSet[0].DhcpConfigurationSet) == 0 {
			return fmt.Errorf("dhcpOption does not found in cloud: %s", saved.Primary.ID)
		}

		foundDhcpOption := res.DhcpOptionsSet[0]

		if nifcloud.StringValue(foundDhcpOption.DhcpOptionsId) != saved.Primary.ID {
			return fmt.Errorf("dhcpOption does not found in cloud: %s", saved.Primary.ID)
		}

		*dhcpOption = foundDhcpOption
		return nil
	}
}

func testAccCheckDhcpOptionValues(dhcpOption *computing.DhcpOptionsSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(dhcpOption.DhcpOptionsId) == "" {
			return fmt.Errorf("bad dhcp_option_id state, expected not nil, got: nil")
		}

		if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].Key) == "default-router" {
			if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "192.168.0.1" {
				return fmt.Errorf("bad default_router state, expected \"192.168.0.1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].Key) == "domain-name" {
			if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "example.com" {
				return fmt.Errorf("bad domain_name state, expected \"example.com\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].Key) == "domain-name-server" {
			if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "192.168.0.1" {
				return fmt.Errorf("bad domain_name_servers state, expected \"192.168.0.1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].Key) == "domain-name-server" {
			if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].ValueSet[1].Value) != "192.168.0.2" {
				return fmt.Errorf("bad domain_name_servers state, expected \"192.168.0.2\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[1].Value)
			}
		}

		if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].Key) == "ntp_servers" {
			if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "192.168.0.1" {
				return fmt.Errorf("bad ntp_servers state, expected \"192.168.0.1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].Key) == "netbios_name_servers" {
			if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "192.168.0.1" {
				return fmt.Errorf("bad netbios_name_servers state, expected \"192.168.0.1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].Key) == "netbios_name_servers" {
			if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].ValueSet[1].Value) != "192.168.0.2" {
				return fmt.Errorf("bad netbios_name_servers state, expected \"192.168.0.2\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[1].Value)
			}
		}

		if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].Key) == "netbios_node_type" {
			if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "1" {
				return fmt.Errorf("bad netbios_node_type state, expected \"1\", got: %#v", dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value)
			}
		}

		if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].Key) == "lease_time" {
			if nifcloud.StringValue(dhcpOption.DhcpConfigurationSet[0].ValueSet[0].Value) != "600" {
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

		res, err := svc.DescribeDhcpOptionsRequest(&computing.DescribeDhcpOptionsInput{
			DhcpOptionsId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.DhcpOption" {
				return fmt.Errorf("failed DescribeDhcpOptionsRequest: %s", err)
			}
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

	res, err := svc.DescribeDhcpOptionsRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepDhcpOptions []string
	for _, k := range res.DhcpOptionsSet {
		if strings.HasPrefix(nifcloud.StringValue(k.DhcpOptionsId), prefix) {
			sweepDhcpOptions = append(sweepDhcpOptions, nifcloud.StringValue(k.DhcpOptionsId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepDhcpOptions {
		dhcpOptionID := n
		eg.Go(func() error {
			_, err := svc.DeleteDhcpOptionsRequest(&computing.DeleteDhcpOptionsInput{
				DhcpOptionsId: nifcloud.String(dhcpOptionID),
			}).Send(ctx)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
