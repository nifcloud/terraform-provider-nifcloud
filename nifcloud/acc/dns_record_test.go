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

var dnsRecordName = os.Getenv("TF_VAR_dns_record_name")
var dnsRecordShorthandName = os.Getenv("TF_VAR_dns_record_shorthand_name")

func init() {
	resource.AddTestSweepers("nifcloud_dns_record", &resource.Sweeper{
		Name: "nifcloud_dns_record",
		F:    testSweepDnsRecord,
	})
}

func TestAcc_DnsRecord_AtSignAsName(t *testing.T) {
	var record types.ResourceRecordSets

	resourceName := "nifcloud_dns_record.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDnsRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(t, "testdata/dns_record_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists(resourceName, &record),
					testAccCheckDnsRecordZoneIdAsNameValues(&record),
					resource.TestCheckResourceAttr(resourceName, "zone_id", dnsZoneName),
					resource.TestCheckResourceAttr(resourceName, "name", "@"),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "60"),
					resource.TestCheckResourceAttr(resourceName, "record", "192.0.2.1"),
					resource.TestCheckResourceAttr(resourceName, "comment", "tfacc-memo"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDnsRecordImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_DnsRecord_ShorthandName(t *testing.T) {
	var record types.ResourceRecordSets

	resourceName := "nifcloud_dns_record.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDnsRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(t, "testdata/dns_record_shorthand_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists(resourceName, &record),
					testAccCheckDnsRecordNameValues(&record),
					resource.TestCheckResourceAttr(resourceName, "zone_id", dnsZoneName),
					resource.TestCheckResourceAttr(resourceName, "name", dnsRecordShorthandName),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "60"),
					resource.TestCheckResourceAttr(resourceName, "record", "192.0.2.1"),
					resource.TestCheckResourceAttr(resourceName, "comment", "tfacc-memo"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDnsRecordImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_DnsRecord_Weight(t *testing.T) {
	var record types.ResourceRecordSets

	resourceName := "nifcloud_dns_record.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDnsRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(t, "testdata/dns_record_weight.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists(resourceName, &record),
					testAccCheckDnsRecordWeightValues(&record),
					resource.TestCheckResourceAttr(resourceName, "zone_id", dnsZoneName),
					resource.TestCheckResourceAttr(resourceName, "name", dnsRecordName),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "60"),
					resource.TestCheckResourceAttr(resourceName, "record", "192.0.2.1"),
					resource.TestCheckResourceAttr(resourceName, "comment", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "weighted_routing_policy.0.weight", "90"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDnsRecordImportStateIDFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"weighted_routing_policy.",
				},
			},
		},
	})
}

func TestAcc_DnsRecord_Failover(t *testing.T) {
	var record types.ResourceRecordSets

	resourceName := "nifcloud_dns_record.basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDnsRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(t, "testdata/dns_record_failover.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists(resourceName, &record),
					testAccCheckDnsRecordFailoverValues(&record),
					resource.TestCheckResourceAttr(resourceName, "zone_id", dnsZoneName),
					resource.TestCheckResourceAttr(resourceName, "name", dnsRecordName),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "60"),
					resource.TestCheckResourceAttr(resourceName, "record", "192.0.2.1"),
					resource.TestCheckResourceAttr(resourceName, "comment", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "failover_routing_policy.0.type", "PRIMARY"),
					resource.TestCheckResourceAttr(resourceName, "failover_routing_policy.0.health_check.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "failover_routing_policy.0.health_check.0.ip_address", "192.0.2.2"),
					resource.TestCheckResourceAttr(resourceName, "failover_routing_policy.0.health_check.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "failover_routing_policy.0.health_check.0.resource_path", "test"),
					resource.TestCheckResourceAttr(resourceName, "failover_routing_policy.0.health_check.0.resource_domain", "example.test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDnsRecordImportStateIDFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"failover_routing_policy.",
				},
			},
		},
	})
}

func testAccDnsRecord(t *testing.T, fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func testAccCheckDnsRecordExists(n string, dnsRecord *types.ResourceRecordSets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no dnsZone resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no dnsZone id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DNS
		res, err := svc.ListResourceRecordSets(context.Background(), &dns.ListResourceRecordSetsInput{
			Identifier: nifcloud.String(saved.Primary.ID),
			ZoneID:     nifcloud.String(dnsZoneName),
		})

		if err != nil {
			return err
		}

		foundDnsRecord := res.ResourceRecordSets[0]

		if nifcloud.ToString(foundDnsRecord.SetIdentifier) != saved.Primary.ID {
			return fmt.Errorf("dnsRecord does not found in cloud: %s", saved.Primary.ID)
		}

		*dnsRecord = foundDnsRecord
		return nil
	}
}

func testAccCheckDnsRecordZoneIdAsNameValues(dnsRecord *types.ResourceRecordSets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(dnsRecord.Name) != dnsZoneName {
			return fmt.Errorf("bad name state, expected %s, got: %#v", dnsZoneName, dnsRecord.Name)
		}

		if nifcloud.ToString(dnsRecord.Type) != "A" {
			return fmt.Errorf("bad type state, expected \"A\", got: %#v", dnsRecord.Type)
		}

		if nifcloud.ToInt32(dnsRecord.TTL) != 60 {
			return fmt.Errorf("bad ttl state, expected 60, got: %#v", dnsRecord.TTL)
		}

		if nifcloud.ToString(dnsRecord.ResourceRecords[0].Value) != "192.0.2.1" {
			return fmt.Errorf("bad resource_records.0.value state, expected \"192.0.2.1\", got: %#v", dnsRecord.ResourceRecords[0].Value)
		}

		if nifcloud.ToString(dnsRecord.XniftyComment) != "tfacc-memo" {
			return fmt.Errorf("bad x_nifty_comment state, expected \"tfacc-memo\", got: %#v", dnsRecord.XniftyComment)
		}

		return nil
	}
}

func testAccCheckDnsRecordNameValues(dnsRecord *types.ResourceRecordSets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(dnsRecord.Name) != dnsRecordName {
			return fmt.Errorf("bad name state, expected %s, got: %#v", dnsRecordName, dnsRecord.Name)
		}

		if nifcloud.ToString(dnsRecord.Type) != "A" {
			return fmt.Errorf("bad type state, expected \"A\", got: %#v", dnsRecord.Type)
		}

		if nifcloud.ToInt32(dnsRecord.TTL) != 60 {
			return fmt.Errorf("bad ttl state, expected 60, got: %#v", dnsRecord.TTL)
		}

		if nifcloud.ToString(dnsRecord.ResourceRecords[0].Value) != "192.0.2.1" {
			return fmt.Errorf("bad resource_records.0.value state, expected \"192.0.2.1\", got: %#v", dnsRecord.ResourceRecords[0].Value)
		}

		if nifcloud.ToString(dnsRecord.XniftyComment) != "tfacc-memo" {
			return fmt.Errorf("bad x_nifty_comment state, expected \"tfacc-memo\", got: %#v", dnsRecord.XniftyComment)
		}

		return nil
	}
}

func testAccCheckDnsRecordWeightValues(dnsRecord *types.ResourceRecordSets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(dnsRecord.Name) != dnsRecordName {
			return fmt.Errorf("bad name state, expected %s, got: %#v", dnsRecordName, dnsRecord.Name)
		}

		if nifcloud.ToString(dnsRecord.Type) != "A" {
			return fmt.Errorf("bad type state, expected \"A\", got: %#v", dnsRecord.Type)
		}

		if nifcloud.ToInt32(dnsRecord.TTL) != 60 {
			return fmt.Errorf("bad ttl state, expected 60, got: %#v", dnsRecord.TTL)
		}

		if nifcloud.ToString(dnsRecord.ResourceRecords[0].Value) != "192.0.2.1" {
			return fmt.Errorf("bad resource_records.0.value state, expected \"192.0.2.1\", got: %#v", dnsRecord.ResourceRecords[0].Value)
		}

		if nifcloud.ToString(dnsRecord.XniftyComment) != "tfacc-memo" {
			return fmt.Errorf("bad x_nifty_comment state, expected \"tfacc-memo\", got: %#v", dnsRecord.XniftyComment)
		}

		if nifcloud.ToInt32(dnsRecord.Weight) != 90 {
			return fmt.Errorf("bad weight state, expected \"90\", got: %#v", dnsRecord.Weight)
		}

		return nil
	}
}

func testAccCheckDnsRecordFailoverValues(dnsRecord *types.ResourceRecordSets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(dnsRecord.Name) != dnsRecordName {
			return fmt.Errorf("bad name state, expected %s, got: %#v", dnsRecordName, dnsRecord.Name)
		}

		if nifcloud.ToString(dnsRecord.Type) != "A" {
			return fmt.Errorf("bad type state, expected \"A\", got: %#v", dnsRecord.Type)
		}

		if nifcloud.ToInt32(dnsRecord.TTL) != 60 {
			return fmt.Errorf("bad ttl state, expected 60, got: %#v", dnsRecord.TTL)
		}

		if nifcloud.ToString(dnsRecord.ResourceRecords[0].Value) != "192.0.2.1" {
			return fmt.Errorf("bad resource_records.0.value state, expected \"192.0.2.1\", got: %#v", dnsRecord.ResourceRecords[0].Value)
		}

		if nifcloud.ToString(dnsRecord.XniftyComment) != "tfacc-memo" {
			return fmt.Errorf("bad x_nifty_comment state, expected \"tfacc-memo\", got: %#v", dnsRecord.XniftyComment)
		}

		if nifcloud.ToString(dnsRecord.Failover) != "PRIMARY" {
			return fmt.Errorf("bad failover state, expected \"PRIMARY\", got: %#v", dnsRecord.Failover)
		}

		if nifcloud.ToString(dnsRecord.XniftyHealthCheckConfig.Protocol) != "HTTPS" {
			return fmt.Errorf("bad x_nifty_health_check_config.protocol state, expected \"HTTPS\", got: %#v", dnsRecord.XniftyHealthCheckConfig.Protocol)
		}

		if nifcloud.ToString(dnsRecord.XniftyHealthCheckConfig.IPAddress) != "192.0.2.2" {
			return fmt.Errorf("bad x_nifty_health_check_config.ipaddress state, expected \"192.0.2.2\", got: %#v", dnsRecord.XniftyHealthCheckConfig.IPAddress)
		}

		if nifcloud.ToInt32(dnsRecord.XniftyHealthCheckConfig.Port) != 443 {
			return fmt.Errorf("bad x_nifty_health_check_config.port state, expected \"443\", got: %#v", dnsRecord.XniftyHealthCheckConfig.Port)
		}

		if nifcloud.ToString(dnsRecord.XniftyHealthCheckConfig.ResourcePath) != "test" {
			return fmt.Errorf("bad x_nifty_health_check_config.resource_path state, expected \"test\", got: %#v", dnsRecord.XniftyHealthCheckConfig.ResourcePath)
		}

		if nifcloud.ToString(dnsRecord.XniftyHealthCheckConfig.FullyQualifiedDomainName) != "example.test" {
			return fmt.Errorf("bad x_nifty_health_check_config.fully_qualified_domain_name state, expected \"example.test\", got: %#v", dnsRecord.XniftyHealthCheckConfig.FullyQualifiedDomainName)
		}

		return nil
	}
}

func testAccDnsRecordResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).DNS

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_dns_record" {
			continue
		}

		res, err := svc.ListResourceRecordSets(context.Background(), &dns.ListResourceRecordSetsInput{
			Identifier: nifcloud.String(rs.Primary.ID),
			ZoneID:     nifcloud.String(dnsZoneName),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "NoSuchHostedZone" {
				return nil
			}
			return fmt.Errorf("failed ListResourceRecordSetsRequest: %s", err)
		}

		if res.ResourceRecordSets[0].Name == nifcloud.String(rs.Primary.ID) {
			return fmt.Errorf("dnsRecord (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepDnsRecord(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).DNS

	res, err := svc.ListResourceRecordSets(ctx, &dns.ListResourceRecordSetsInput{
		ZoneID: nifcloud.String(dnsZoneName),
	})
	if err != nil {
		return err
	}

	for _, resourceRecordSet := range res.ResourceRecordSets {
		if strings.HasPrefix(nifcloud.ToString(resourceRecordSet.XniftyComment), prefix) {
			input := &dns.ChangeResourceRecordSetsInput{
				ZoneID: nifcloud.String(dnsZoneName),
				RequestChangeBatch: &types.RequestChangeBatch{
					ListOfRequestChanges: []types.RequestChanges{{
						RequestChange: &types.RequestChange{
							Action: types.ActionOfChangeResourceRecordSetsRequestForChangeResourceRecordSetsDelete,
							RequestResourceRecordSet: &types.RequestResourceRecordSet{
								Name:              resourceRecordSet.Name,
								SetIdentifier:     resourceRecordSet.SetIdentifier,
								TTL:               resourceRecordSet.TTL,
								Type:              types.TypeOfChangeResourceRecordSetsRequestForChangeResourceRecordSets(nifcloud.ToString(resourceRecordSet.Type)),
								XniftyComment:     resourceRecordSet.XniftyComment,
								XniftyDefaultHost: resourceRecordSet.XniftyDefaultHost,
								ListOfRequestResourceRecords: []types.RequestResourceRecords{{
									RequestResourceRecord: &types.RequestResourceRecord{
										Value: resourceRecordSet.ResourceRecords[0].Value,
									},
								}},
							},
						},
					}},
				},
			}

			_, err := svc.ChangeResourceRecordSets(ctx, input)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func testAccDnsRecordImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		setIdentifier := rs.Primary.Attributes["set_identifier"]
		zoneId := rs.Primary.Attributes["zone_id"]
		name := rs.Primary.Attributes["name"]

		var parts []string
		parts = append(parts, setIdentifier)
		parts = append(parts, zoneId)
		parts = append(parts, name)

		id := strings.Join(parts, "_")
		return id, nil
	}
}
