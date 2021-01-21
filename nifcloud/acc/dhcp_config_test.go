package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func init() {
	resource.AddTestSweepers("nifcloud_dhcp_config", &resource.Sweeper{
		Name: "nifcloud_dhcp_config",
		F:    testSweepDhcpConfig,
		Dependencies: []string{
			"nifcloud_router",
		},
	})
}

func TestAcc_DhcpConfig(t *testing.T) {
	var dhcpConfig computing.DhcpConfigsSet

	resourceName := "nifcloud_dhcp_config.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDhcpConfigResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpConfig(t, "testdata/dhcp_config.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDhcpConfigExists(resourceName, &dhcpConfig),
					testAccCheckDhcpConfigValues(&dhcpConfig),
					resource.TestCheckResourceAttr(resourceName, "static_mapping.0.static_mapping_ipaddress", "192.168.1.10"),
					resource.TestCheckResourceAttr(resourceName, "static_mapping.0.static_mapping_macaddress", "00:00:5e:00:53:00"),
					resource.TestCheckResourceAttr(resourceName, "static_mapping.0.static_mapping_description", "static-mapping-memo"),
					resource.TestCheckResourceAttr(resourceName, "ipaddress_pool.0.ipaddress_pool_start", "192.168.2.1"),
					resource.TestCheckResourceAttr(resourceName, "ipaddress_pool.0.ipaddress_pool_stop", "192.168.2.100"),
					resource.TestCheckResourceAttr(resourceName, "ipaddress_pool.0.ipaddress_pool_description", "ipaddress-pool-memo"),
				),
			},
			{
				Config: testAccDhcpConfig(t, "testdata/dhcp_config_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDhcpConfigExists(resourceName, &dhcpConfig),
					testAccCheckDhcpConfigValuesUpdated(&dhcpConfig),
					resource.TestCheckResourceAttr(resourceName, "static_mapping.0.static_mapping_ipaddress", "192.168.2.10"),
					resource.TestCheckResourceAttr(resourceName, "static_mapping.0.static_mapping_macaddress", "00:00:5e:00:53:FF"),
					resource.TestCheckResourceAttr(resourceName, "static_mapping.0.static_mapping_description", "static-mapping-memo-upd"),
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

func testAccDhcpConfig(t *testing.T, fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func testAccCheckDhcpConfigExists(n string, dhcpConfig *computing.DhcpConfigsSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no dhcpConfig resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no dhcpConfig id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.NiftyDescribeDhcpConfigsRequest(&computing.NiftyDescribeDhcpConfigsInput{
			DhcpConfigId: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.DhcpConfigsSet) == 0 {
			return fmt.Errorf("dhcpConfig does not found in cloud: %s", saved.Primary.ID)
		}

		foundDhcpConfig := res.DhcpConfigsSet[0]

		if nifcloud.StringValue(foundDhcpConfig.DhcpConfigId) != saved.Primary.ID {
			return fmt.Errorf("dhcpConfig does not found in cloud: %s", saved.Primary.ID)
		}

		*dhcpConfig = foundDhcpConfig
		return nil
	}
}

func testAccCheckDhcpConfigValues(dhcpConfig *computing.DhcpConfigsSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(dhcpConfig.StaticMappingsSet) != 1 {
			return fmt.Errorf("bad static mappings: %#v", dhcpConfig.StaticMappingsSet)
		}

		if len(dhcpConfig.IpAddressPoolsSet) != 1 {
			return fmt.Errorf("bad ipaddress pools: %#v", dhcpConfig.IpAddressPoolsSet)
		}

		staticmappings := make(map[string]computing.StaticMappingsSet)
		for _, r := range dhcpConfig.StaticMappingsSet {
			staticmappings[*r.Description] = r
		}

		ipaddresspools := make(map[string]computing.IpAddressPoolsSet)
		for _, r := range dhcpConfig.IpAddressPoolsSet {
			ipaddresspools[*r.Description] = r
		}

		if _, ok := staticmappings["static-mapping-memo"]; !ok {
			return fmt.Errorf("bad static mapping: %#v", dhcpConfig.StaticMappingsSet)
		}

		if _, ok := ipaddresspools["ipaddress-pool-memo"]; !ok {
			return fmt.Errorf("bad ipaddress pool: %#v", dhcpConfig.IpAddressPoolsSet)
		}

		if nifcloud.StringValue(staticmappings["static-mapping-memo"].IpAddress) != "192.168.1.10" {
			return fmt.Errorf("bad static mapping IP address, expected \"192.168.1.10\", got: %#v", nifcloud.StringValue(staticmappings["static-mapping-memo"].IpAddress))
		}

		if nifcloud.StringValue(staticmappings["static-mapping-memo"].MacAddress) != "00:00:5e:00:53:00" {
			return fmt.Errorf("bad static mapping MAC address, expected \"00:00:5e:00:53:00\", got: %#v", nifcloud.StringValue(staticmappings["static-mapping-memo"].MacAddress))
		}

		if nifcloud.StringValue(ipaddresspools["ipaddress-pool-memo"].StartIpAddress) != "192.168.2.1" {
			return fmt.Errorf("bad ipaddress pool start IP address, expected \"192.168.2.1\", got: %#v", nifcloud.StringValue(ipaddresspools["ipaddress-pool-memo"].StartIpAddress))
		}

		if nifcloud.StringValue(ipaddresspools["ipaddress-pool-memo"].StopIpAddress) != "192.168.2.100" {
			return fmt.Errorf("bad ipaddress pool stop IP address, expected \"192.168.2.100\", got: %#v", nifcloud.StringValue(ipaddresspools["ipaddress-pool-memo"].StopIpAddress))
		}

		return nil
	}
}

func testAccCheckDhcpConfigValuesUpdated(dhcpConfig *computing.DhcpConfigsSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(dhcpConfig.StaticMappingsSet) != 1 {
			return fmt.Errorf("bad static mappings: %#v", dhcpConfig.StaticMappingsSet)
		}

		staticmappings := make(map[string]computing.StaticMappingsSet)
		for _, r := range dhcpConfig.StaticMappingsSet {
			staticmappings[*r.Description] = r
		}

		if _, ok := staticmappings["static-mapping-memo-upd"]; !ok {
			return fmt.Errorf("bad static mapping: %#v", dhcpConfig.StaticMappingsSet)
		}

		if nifcloud.StringValue(staticmappings["static-mapping-memo-upd"].IpAddress) != "192.168.2.10" {
			return fmt.Errorf("bad static mapping IP address, expected \"192.168.2.10\", got: %#v", nifcloud.StringValue(staticmappings["static-mapping-memo-upd"].IpAddress))
		}

		if nifcloud.StringValue(staticmappings["static-mapping-memo-upd"].MacAddress) != "00:00:5e:00:53:FF" {
			return fmt.Errorf("bad static mapping MAC address, expected \"00:00:5e:00:53:FF\", got: %#v", nifcloud.StringValue(staticmappings["static-mapping-memo-upd"].MacAddress))
		}

		return nil
	}
}

func testAccDhcpConfigResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_dhcp_config" {
			continue
		}

		res, err := svc.NiftyDescribeDhcpConfigsRequest(&computing.NiftyDescribeDhcpConfigsInput{
			DhcpConfigId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.DhcpConfigId" {
				return fmt.Errorf("failed NiftyDescribeDhcpConfigsRequest: %s", err)
			}
		}

		if len(res.DhcpConfigsSet) > 0 {
			return fmt.Errorf("dhcpConfig (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepDhcpConfig(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.NiftyDescribeDhcpConfigsRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	for _, dhcpConfig := range res.DhcpConfigsSet {

		input := &computing.NiftyDeleteDhcpConfigInput{
			DhcpConfigId: dhcpConfig.DhcpConfigId,
		}

		_, err := svc.NiftyDeleteDhcpConfigRequest(input).Send(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
