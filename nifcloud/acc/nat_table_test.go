package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
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
	var natTable types.NatTableSet

	resourceName := "nifcloud_nat_table.basic"
	randName := prefix + acctest.RandString(7)

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
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccCheckNatTableExists(n string, natTable *types.NatTableSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no natTable resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no natTable id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.NiftyDescribeNatTables(context.Background(), &computing.NiftyDescribeNatTablesInput{
			NatTableId: []string{saved.Primary.ID},
		})

		if err != nil {
			return err
		}

		if len(res.NatTableSet) == 0 {
			return fmt.Errorf("natTable does not found in cloud: %s", saved.Primary.ID)
		}

		foundNatTable := res.NatTableSet[0]

		if nifcloud.ToString(foundNatTable.NatTableId) != saved.Primary.ID {
			return fmt.Errorf("natTable does not found in cloud: %s", saved.Primary.ID)
		}

		*natTable = foundNatTable
		return nil
	}
}

func testAccCheckNatTableValues(natTable *types.NatTableSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(natTable.NatRuleSet) != 2 {
			return fmt.Errorf("bad nat rules: %#v", natTable.NatRuleSet)
		}

		rules := make(map[string]types.NatRuleSet)
		for _, r := range natTable.NatRuleSet {
			rules[*r.Description] = r
		}

		if _, ok := rules["snat-memo"]; !ok {
			return fmt.Errorf("bad snat rule: %#v", natTable.NatRuleSet)
		}

		if _, ok := rules["dnat-memo"]; !ok {
			return fmt.Errorf("bad dnat rule: %#v", natTable.NatRuleSet)
		}

		if nifcloud.ToString(rules["snat-memo"].RuleNumber) != "1" {
			return fmt.Errorf("bad snat rule number, expected \"1\", got: %#v", nifcloud.ToString(rules["snat-memo"].RuleNumber))
		}

		if nifcloud.ToString(rules["snat-memo"].Protocol) != "TCP" {
			return fmt.Errorf("bad snat protocol, expected \"TCP\", got: %#v", nifcloud.ToString(rules["snat-memo"].Protocol))
		}

		if nifcloud.ToString(rules["snat-memo"].Source.Address) != "192.0.2.1" {
			return fmt.Errorf("bad snat source address, expected \"192.0.2.1\", got: %#v", nifcloud.ToString(rules["snat-memo"].Source.Address))
		}

		if nifcloud.ToInt32(rules["snat-memo"].Source.Port) != 80 {
			return fmt.Errorf("bad snat source port, expected \"80\", got: %#v", nifcloud.ToInt32(rules["snat-memo"].Source.Port))
		}

		if nifcloud.ToInt32(rules["snat-memo"].Translation.Port) != 81 {
			return fmt.Errorf("bad snat translation port, expected \"81\", got: %#v", nifcloud.ToInt32(rules["snat-memo"].Translation.Port))
		}

		if nifcloud.ToString(rules["snat-memo"].OutboundInterface.NetworkId) != "net-COMMON_PRIVATE" {
			return fmt.Errorf("bad snat outbound interface network id, expected \"net-COMMON_PRIVATE\", got: %#v", nifcloud.ToString(rules["snat-memo"].OutboundInterface.NetworkId))
		}

		if nifcloud.ToString(rules["dnat-memo"].RuleNumber) != "1" {
			return fmt.Errorf("bad dnat rule number, expected \"1\", got: %#v", nifcloud.ToString(rules["dnat-memo"].RuleNumber))
		}

		if nifcloud.ToString(rules["dnat-memo"].Protocol) != "ALL" {
			return fmt.Errorf("bad dnat protocol, expected \"ALL\", got: %#v", nifcloud.ToString(rules["dnat-memo"].Protocol))
		}

		if nifcloud.ToString(rules["dnat-memo"].Translation.Address) != "192.168.1.1" {
			return fmt.Errorf("bad dnat translation address, expected \"192.168.1.1\", got: %#v", nifcloud.ToString(rules["dnat-memo"].Translation.Address))
		}

		if nifcloud.ToString(rules["dnat-memo"].InboundInterface.NetworkId) == "" {
			return fmt.Errorf("bad dnat inbound interface network id, expected \"not null\", got: %#v", nifcloud.ToString(rules["dnat-memo"].InboundInterface.NetworkId))
		}

		return nil
	}
}

func testAccCheckNatTableValuesUpdated(natTable *types.NatTableSet, privateLanName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(natTable.NatRuleSet) != 1 {
			return fmt.Errorf("bad nat rules: %#v", natTable.NatRuleSet)
		}

		rules := make(map[string]types.NatRuleSet)
		for _, r := range natTable.NatRuleSet {
			rules[*r.Description] = r
		}

		if _, ok := rules["dnat-memo-upd"]; !ok {
			return fmt.Errorf("bad dnat rule: %#v", natTable.NatRuleSet)
		}

		if nifcloud.ToString(rules["dnat-memo-upd"].RuleNumber) != "2" {
			return fmt.Errorf("bad dnat rule number, expected \"2\", got: %#v", nifcloud.ToString(rules["dnat-memo-upd"].RuleNumber))
		}

		if nifcloud.ToString(rules["dnat-memo-upd"].Protocol) != "TCP" {
			return fmt.Errorf("bad dnat protocol, expected \"TCP\", got: %#v", nifcloud.ToString(rules["dnat-memo-upd"].Protocol))
		}

		if nifcloud.ToInt32(rules["dnat-memo-upd"].Destination.Port) != 80 {
			return fmt.Errorf("bad dnat destination port, expected \"80\", got: %#v", nifcloud.ToInt32(rules["dnat-memo-upd"].Destination.Port))
		}

		if nifcloud.ToString(rules["dnat-memo-upd"].Translation.Address) != "192.168.1.2" {
			return fmt.Errorf("bad dnat translation address, expected \"192.168.1.2\", got: %#v", nifcloud.ToString(rules["dnat-memo-upd"].Translation.Address))
		}

		if nifcloud.ToInt32(rules["dnat-memo-upd"].Translation.Port) != 81 {
			return fmt.Errorf("bad dnat translation port, expected \"81\", got: %#v", nifcloud.ToInt32(rules["dnat-memo-upd"].Translation.Port))
		}

		if nifcloud.ToString(rules["dnat-memo-upd"].InboundInterface.NetworkName) != privateLanName {
			return fmt.Errorf("bad dnat inbound interface network name, expected \"%s\", got: %#v", privateLanName, nifcloud.ToString(rules["dnat-memo-upd"].InboundInterface.NetworkId))
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

		res, err := svc.NiftyDescribeNatTables(context.Background(), &computing.NiftyDescribeNatTablesInput{
			NatTableId: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.NatTableID" {
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

	res, err := svc.NiftyDescribeNatTables(ctx, nil)
	if err != nil {
		return err
	}

	for _, natTable := range res.NatTableSet {

		for _, natTableAssociation := range natTable.AssociationSet {

			input := &computing.NiftyDisassociateNatTableInput{
				AssociationId: natTableAssociation.AssociationId,
			}

			_, err := svc.NiftyDisassociateNatTable(ctx, input)
			if err != nil {
				return err
			}
		}

		input := &computing.NiftyDeleteNatTableInput{
			NatTableId: natTable.NatTableId,
		}

		_, err := svc.NiftyDeleteNatTable(ctx, input)
		if err != nil {
			return err
		}
	}
	return nil
}
