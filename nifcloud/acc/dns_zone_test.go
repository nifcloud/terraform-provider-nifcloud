package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

var dnsZoneName = os.Getenv("TF_VAR_dns_zone_name")

func init() {
	resource.AddTestSweepers("nifcloud_dns_zone", &resource.Sweeper{
		Name: "nifcloud_dns_zone",
		F:    testSweepDnsZone,
		Dependencies: []string{
			"nifcloud_dns_record",
		},
	})
}

func TestAcc_DnsZone(t *testing.T) {
	var zone types.HostedZone

	resourceName := "nifcloud_dns_zone.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDnsZoneResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsZone(t, "testdata/dns_zone.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsZoneExists(resourceName, &zone),
					testAccCheckDnsZoneValues(&zone),
					resource.TestCheckResourceAttr(resourceName, "name", dnsZoneName),
					resource.TestCheckResourceAttr(resourceName, "comment", "tfacc-memo"),
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

func testAccDnsZone(t *testing.T, fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func testAccCheckDnsZoneExists(n string, dnsZone *types.HostedZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no dnsZone resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no dnsZone id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DNS
		res, err := svc.GetHostedZone(context.Background(), &dns.GetHostedZoneInput{
			ZoneID: nifcloud.String(saved.Primary.ID),
		})

		if err != nil {
			return err
		}

		foundDnsZone := res.HostedZone

		if nifcloud.ToString(foundDnsZone.Name) != saved.Primary.ID {
			return fmt.Errorf("dnsZone does not found in cloud: %s", saved.Primary.ID)
		}

		*dnsZone = *foundDnsZone
		return nil
	}
}

func testAccCheckDnsZoneValues(dnsZone *types.HostedZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(dnsZone.Name) != dnsZoneName {
			return fmt.Errorf("bad name state, expected %s, got: %#v", dnsZoneName, dnsZone.Name)
		}

		if nifcloud.ToString(dnsZone.Config.Comment) != "tfacc-memo" {
			return fmt.Errorf("bad comment state, expected \"tfacc-memo\", got: %#v", dnsZone.Config.Comment)
		}

		return nil
	}
}

func testAccDnsZoneResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).DNS

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_dns_zone" {
			continue
		}

		res, err := svc.GetHostedZone(context.Background(), &dns.GetHostedZoneInput{
			ZoneID: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "NoSuchHostedZone" {
				return nil
			}
			return fmt.Errorf("failed GetHostedZoneRequest: %s", err)
		}

		if res.HostedZone.Name == nifcloud.String(rs.Primary.ID) {
			return fmt.Errorf("dnsZone (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepDnsZone(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).DNS

	res, err := svc.ListHostedZones(ctx, nil)
	if err != nil {
		return err
	}

	for _, dnsZone := range res.HostedZones {
		if strings.HasPrefix(nifcloud.ToString(dnsZone.Config.Comment), prefix) {
			input := &dns.DeleteHostedZoneInput{
				ZoneID: dnsZone.Name,
			}

			_, err := svc.DeleteHostedZone(ctx, input)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
