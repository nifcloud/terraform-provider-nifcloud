package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

var dnsRecordName = os.Getenv("TF_VAR_dns_record_name")

func init() {
	resource.AddTestSweepers("nifcloud_dns_record", &resource.Sweeper{
		Name: "nifcloud_dns_record",
		F:    testSweepDnsRecord,
	})
}

func TestAcc_DnsRecord_Weight(t *testing.T) {
	var record dns.ResourceRecordSets

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
	var record dns.ResourceRecordSets

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

func testAccCheckDnsRecordExists(n string, dnsRecord *dns.ResourceRecordSets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no dnsZone resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no dnsZone id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DNS
		res, err := svc.ListResourceRecordSetsRequest(&dns.ListResourceRecordSetsInput{
			Identifier: nifcloud.String(saved.Primary.ID),
			ZoneID:     nifcloud.String(dnsZoneName),
		}).Send(context.Background())

		if err != nil {
			return err
		}

		foundDnsRecord := res.ResourceRecordSets[0]

		if nifcloud.StringValue(foundDnsRecord.SetIdentifier) != saved.Primary.ID {
			return fmt.Errorf("dnsRecord does not found in cloud: %s", saved.Primary.ID)
		}

		*dnsRecord = foundDnsRecord
		return nil
	}
}

func testAccCheckDnsRecordWeightValues(dnsRecord *dns.ResourceRecordSets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(dnsRecord.Name) != dnsRecordName {
			return fmt.Errorf("bad name state, expected %s, got: %#v", dnsRecordName, dnsRecord.Name)
		}

		if nifcloud.StringValue(dnsRecord.Type) != "A" {
			return fmt.Errorf("bad comment state, expected \"A\", got: %#v", dnsRecord.Type)
		}

		if nifcloud.Int64Value(dnsRecord.TTL) != 60 {
			return fmt.Errorf("bad comment state, expected 60, got: %#v", dnsRecord.TTL)
		}

		if nifcloud.StringValue(dnsRecord.ResourceRecords[0].Value) != "192.0.2.1" {
			return fmt.Errorf("bad comment state, expected \"192.0.2.1\", got: %#v", dnsRecord.ResourceRecords[0].Value)
		}

		if nifcloud.StringValue(dnsRecord.XniftyComment) != "tfacc-memo" {
			return fmt.Errorf("bad comment state, expected \"tfacc-memo\", got: %#v", dnsRecord.XniftyComment)
		}

		if nifcloud.Int64Value(dnsRecord.Weight) != 90 {
			return fmt.Errorf("bad comment state, expected \"90\", got: %#v", dnsRecord.Weight)
		}

		return nil
	}
}

func testAccCheckDnsRecordFailoverValues(dnsRecord *dns.ResourceRecordSets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(dnsRecord.Name) != dnsRecordName {
			return fmt.Errorf("bad name state, expected %s, got: %#v", dnsRecordName, dnsRecord.Name)
		}

		if nifcloud.StringValue(dnsRecord.Type) != "A" {
			return fmt.Errorf("bad comment state, expected \"A\", got: %#v", dnsRecord.Type)
		}

		if nifcloud.Int64Value(dnsRecord.TTL) != 60 {
			return fmt.Errorf("bad comment state, expected 60, got: %#v", dnsRecord.TTL)
		}

		if nifcloud.StringValue(dnsRecord.ResourceRecords[0].Value) != "192.0.2.1" {
			return fmt.Errorf("bad comment state, expected \"192.0.2.1\", got: %#v", dnsRecord.ResourceRecords[0].Value)
		}

		if nifcloud.StringValue(dnsRecord.XniftyComment) != "tfacc-memo" {
			return fmt.Errorf("bad comment state, expected \"tfacc-memo\", got: %#v", dnsRecord.XniftyComment)
		}

		if nifcloud.StringValue(dnsRecord.Failover) != "PRIMARY" {
			return fmt.Errorf("bad comment state, expected \"PRIMARY\", got: %#v", dnsRecord.Failover)
		}

		if nifcloud.StringValue(dnsRecord.XniftyHealthCheckConfig.Protocol) != "HTTPS" {
			return fmt.Errorf("bad comment state, expected \"HTTPS\", got: %#v", dnsRecord.XniftyHealthCheckConfig.Protocol)
		}

		if nifcloud.StringValue(dnsRecord.XniftyHealthCheckConfig.IPAddress) != "192.0.2.2" {
			return fmt.Errorf("bad comment state, expected \"192.0.2.2\", got: %#v", dnsRecord.XniftyHealthCheckConfig.IPAddress)
		}

		if nifcloud.Int64Value(dnsRecord.XniftyHealthCheckConfig.Port) != 443 {
			return fmt.Errorf("bad comment state, expected \"443\", got: %#v", dnsRecord.XniftyHealthCheckConfig.Port)
		}

		if nifcloud.StringValue(dnsRecord.XniftyHealthCheckConfig.ResourcePath) != "test" {
			return fmt.Errorf("bad comment state, expected \"test\", got: %#v", dnsRecord.XniftyHealthCheckConfig.ResourcePath)
		}

		if nifcloud.StringValue(dnsRecord.XniftyHealthCheckConfig.FullyQualifiedDomainName) != "example.test" {
			return fmt.Errorf("bad comment state, expected \"example.test\", got: %#v", dnsRecord.XniftyHealthCheckConfig.FullyQualifiedDomainName)
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

		res, err := svc.ListResourceRecordSetsRequest(&dns.ListResourceRecordSetsInput{
			Identifier: nifcloud.String(rs.Primary.ID),
			ZoneID:     nifcloud.String(dnsZoneName),
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "NoSuchHostedZone" {
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

	res, err := svc.ListResourceRecordSetsRequest(&dns.ListResourceRecordSetsInput{
		ZoneID: nifcloud.String(dnsZoneName),
	}).Send(ctx)
	if err != nil {
		return err
	}

	for _, resourceRecordSet := range res.ResourceRecordSets {
		if strings.HasPrefix(nifcloud.StringValue(resourceRecordSet.XniftyComment), prefix) {
			input := &dns.ChangeResourceRecordSetsInput{
				ZoneID: nifcloud.String(dnsZoneName),
				RequestChangeBatch: &dns.RequestChangeBatch{
					ListOfRequestChanges: []dns.RequestChanges{{
						RequestChange: &dns.RequestChange{
							Action: nifcloud.String("DELETE"),
							RequestResourceRecordSet: &dns.RequestResourceRecordSet{
								Name:              resourceRecordSet.Name,
								SetIdentifier:     resourceRecordSet.SetIdentifier,
								TTL:               resourceRecordSet.TTL,
								Type:              resourceRecordSet.Type,
								XniftyComment:     resourceRecordSet.XniftyComment,
								XniftyDefaultHost: resourceRecordSet.XniftyDefaultHost,
								ListOfRequestResourceRecords: []dns.RequestResourceRecords{{
									RequestResourceRecord: &dns.RequestResourceRecord{
										Value: resourceRecordSet.ResourceRecords[0].Value,
									},
								}},
							},
						},
					}},
				},
			}

			_, err := svc.ChangeResourceRecordSetsRequest(input).Send(ctx)
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

		var parts []string
		parts = append(parts, setIdentifier)
		parts = append(parts, zoneId)

		id := strings.Join(parts, "_")
		return id, nil
	}
}
