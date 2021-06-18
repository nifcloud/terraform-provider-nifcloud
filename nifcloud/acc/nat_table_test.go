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
	resource.AddTestSweepers("nifcloud_nat_table", &resource.Sweeper{
		Name: "nifcloud_nat_table",
		F:    testSweepNatTable,
		Dependencies: []string{
			"nifcloud_router",
		},
	})
}

func TestAcc_NatTable(t *testing.T) {
	var natTable computing.NatTableSet

	resourceName := "nifcloud_nat_table.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccNatTableResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatTable(t, "testdata/nat_table.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatTableExists(resourceName, &natTable),
					testAccCheckNatTableValues(&natTable),
					resource.TestCheckResourceAttr(resourceName, "snat.0.rule_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "snat.0.description", "snat-memo"),
					resource.TestCheckResourceAttr(resourceName, "snat.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "snat.0.source_address", "192.0.2.1"),
					resource.TestCheckResourceAttr(resourceName, "snat.0.source_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "snat.0.translation_port", "81"),
					resource.TestCheckResourceAttr(resourceName, "snat.0.outbound_interface_network_id", "net-COMMON_PRIVATE"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.rule_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.description", "dnat-memo"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.protocol", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.translation_address", "192.168.1.1"),
					resource.TestCheckResourceAttrSet(resourceName, "dnat.0.inbound_interface_network_id"),
				),
			},
			{
				Config: testAccNatTable(t, "testdata/nat_table_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatTableExists(resourceName, &natTable),
					testAccCheckNatTableValuesUpdated(&natTable, randName),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.rule_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.description", "dnat-memo"+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.destination_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.translation_address", "192.168.1.2"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.translation_port", "81"),
					resource.TestCheckResourceAttr(resourceName, "dnat.0.inbound_interface_network_name", randName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"dnat.0.inbound_interface_network",
				},
			},
		},
	})
}

func testAccNatTable(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccCheckNatTableExists(n string, natTable *computing.NatTableSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no natTable resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no natTable id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.NiftyDescribeNatTablesRequest(&computing.NiftyDescribeNatTablesInput{
			NatTableId: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.NatTableSet) == 0 {
			return fmt.Errorf("natTable does not found in cloud: %s", saved.Primary.ID)
		}

		foundNatTable := res.NatTableSet[0]

		if nifcloud.StringValue(foundNatTable.NatTableId) != saved.Primary.ID {
			return fmt.Errorf("natTable does not found in cloud: %s", saved.Primary.ID)
		}

		*natTable = foundNatTable
		return nil
	}
}

func testAccCheckNatTableValues(natTable *computing.NatTableSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(natTable.NatRuleSet) != 2 {
			return fmt.Errorf("bad nat rules: %#v", natTable.NatRuleSet)
		}

		rules := make(map[string]computing.NatRuleSet)
		for _, r := range natTable.NatRuleSet {
			rules[*r.Description] = r
		}

		if _, ok := rules["snat-memo"]; !ok {
			return fmt.Errorf("bad snat rule: %#v", natTable.NatRuleSet)
		}

		if _, ok := rules["dnat-memo"]; !ok {
			return fmt.Errorf("bad dnat rule: %#v", natTable.NatRuleSet)
		}

		if nifcloud.StringValue(rules["snat-memo"].RuleNumber) != "1" {
			return fmt.Errorf("bad snat rule number, expected \"1\", got: %#v", nifcloud.StringValue(rules["snat-memo"].RuleNumber))
		}

		if nifcloud.StringValue(rules["snat-memo"].Protocol) != "TCP" {
			return fmt.Errorf("bad snat protocol, expected \"TCP\", got: %#v", nifcloud.StringValue(rules["snat-memo"].Protocol))
		}

		if nifcloud.StringValue(rules["snat-memo"].Source.Address) != "192.0.2.1" {
			return fmt.Errorf("bad snat source address, expected \"192.0.2.1\", got: %#v", nifcloud.StringValue(rules["snat-memo"].Source.Address))
		}

		if nifcloud.Int64Value(rules["snat-memo"].Source.Port) != 80 {
			return fmt.Errorf("bad snat source port, expected \"80\", got: %#v", nifcloud.Int64Value(rules["snat-memo"].Source.Port))
		}

		if nifcloud.Int64Value(rules["snat-memo"].Translation.Port) != 81 {
			return fmt.Errorf("bad snat translation port, expected \"81\", got: %#v", nifcloud.Int64Value(rules["snat-memo"].Translation.Port))
		}

		if nifcloud.StringValue(rules["snat-memo"].OutboundInterface.NetworkId) != "net-COMMON_PRIVATE" {
			return fmt.Errorf("bad snat outbound interface network id, expected \"net-COMMON_PRIVATE\", got: %#v", nifcloud.StringValue(rules["snat-memo"].OutboundInterface.NetworkId))
		}

		if nifcloud.StringValue(rules["dnat-memo"].RuleNumber) != "1" {
			return fmt.Errorf("bad dnat rule number, expected \"1\", got: %#v", nifcloud.StringValue(rules["dnat-memo"].RuleNumber))
		}

		if nifcloud.StringValue(rules["dnat-memo"].Protocol) != "ALL" {
			return fmt.Errorf("bad dnat protocol, expected \"ALL\", got: %#v", nifcloud.StringValue(rules["dnat-memo"].Protocol))
		}

		if nifcloud.StringValue(rules["dnat-memo"].Translation.Address) != "192.168.1.1" {
			return fmt.Errorf("bad dnat translation address, expected \"192.168.1.1\", got: %#v", nifcloud.StringValue(rules["dnat-memo"].Translation.Address))
		}

		if nifcloud.StringValue(rules["dnat-memo"].InboundInterface.NetworkId) == "" {
			return fmt.Errorf("bad dnat inbound interface network id, expected \"not null\", got: %#v", nifcloud.StringValue(rules["dnat-memo"].InboundInterface.NetworkId))
		}

		return nil
	}
}

func testAccCheckNatTableValuesUpdated(natTable *computing.NatTableSet, privateLanName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(natTable.NatRuleSet) != 1 {
			return fmt.Errorf("bad nat rules: %#v", natTable.NatRuleSet)
		}

		rules := make(map[string]computing.NatRuleSet)
		for _, r := range natTable.NatRuleSet {
			rules[*r.Description] = r
		}

		if _, ok := rules["dnat-memo-upd"]; !ok {
			return fmt.Errorf("bad dnat rule: %#v", natTable.NatRuleSet)
		}

		if nifcloud.StringValue(rules["dnat-memo-upd"].RuleNumber) != "2" {
			return fmt.Errorf("bad dnat rule number, expected \"2\", got: %#v", nifcloud.StringValue(rules["dnat-memo-upd"].RuleNumber))
		}

		if nifcloud.StringValue(rules["dnat-memo-upd"].Protocol) != "TCP" {
			return fmt.Errorf("bad dnat protocol, expected \"TCP\", got: %#v", nifcloud.StringValue(rules["dnat-memo-upd"].Protocol))
		}

		if nifcloud.Int64Value(rules["dnat-memo-upd"].Destination.Port) != 80 {
			return fmt.Errorf("bad dnat destination port, expected \"80\", got: %#v", nifcloud.Int64Value(rules["dnat-memo-upd"].Destination.Port))
		}

		if nifcloud.StringValue(rules["dnat-memo-upd"].Translation.Address) != "192.168.1.2" {
			return fmt.Errorf("bad dnat translation address, expected \"192.168.1.2\", got: %#v", nifcloud.StringValue(rules["dnat-memo-upd"].Translation.Address))
		}

		if nifcloud.Int64Value(rules["dnat-memo-upd"].Translation.Port) != 81 {
			return fmt.Errorf("bad dnat translation port, expected \"81\", got: %#v", nifcloud.Int64Value(rules["dnat-memo-upd"].Translation.Port))
		}

		if nifcloud.StringValue(rules["dnat-memo-upd"].InboundInterface.NetworkName) != privateLanName {
			return fmt.Errorf("bad dnat inbound interface network name, expected \"%s\", got: %#v", privateLanName, nifcloud.StringValue(rules["dnat-memo-upd"].InboundInterface.NetworkId))
		}
		return nil
	}
}

func testAccNatTableResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_nat_table" {
			continue
		}

		res, err := svc.NiftyDescribeNatTablesRequest(&computing.NiftyDescribeNatTablesInput{
			NatTableId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.NatTableID" {
				return nil
			}
			return fmt.Errorf("failed NiftyDescribeNatTablesRequest: %s", err)
		}

		if len(res.NatTableSet) > 0 {
			return fmt.Errorf("natTable (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepNatTable(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.NiftyDescribeNatTablesRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	for _, natTable := range res.NatTableSet {

		for _, natTableAssociation := range natTable.AssociationSet {

			input := &computing.NiftyDisassociateNatTableInput{
				AssociationId: natTableAssociation.AssociationId,
			}

			_, err := svc.NiftyDisassociateNatTableRequest(input).Send(ctx)
			if err != nil {
				return err
			}
		}

		input := &computing.NiftyDeleteNatTableInput{
			NatTableId: natTable.NatTableId,
		}

		_, err := svc.NiftyDeleteNatTableRequest(input).Send(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
