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
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func init() {
	resource.AddTestSweepers("nifcloud_hatoba_firewall_group", &resource.Sweeper{
		Name: "nifcloud_hatoba_firewall_group",
		F:    testSweepHatobaFirewallGroup,
		Dependencies: []string{
			"nifcloud_hatoba_cluster",
		},
	})
}

func TestAcc_HatobaFirewallGroup(t *testing.T) {
	var firewallGroup types.FirewallGroup

	resourceName := "nifcloud_hatoba_firewall_group.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccHatobaFirewallGroupResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHatobaFirewallGroup(t, "testdata/hatoba_firewall_group.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHatobaFirewallGroupExists(resourceName, &firewallGroup),
					testAccCheckHatobaFirewallGroupValues(&firewallGroup, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "4"),
				),
			},
			{
				Config: testAccHatobaFirewallGroup(t, "testdata/hatoba_firewall_group_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHatobaFirewallGroupExists(resourceName, &firewallGroup),
					testAccCheckHatobaFirewallGroupValuesUpdated(&firewallGroup, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "3"),
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

func testAccHatobaFirewallGroup(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccCheckHatobaFirewallGroupExists(n string, firewallGroup *types.FirewallGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no types firewall group resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no types firewall group id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Hatoba
		res, err := svc.GetFirewallGroup(context.Background(), &hatoba.GetFirewallGroupInput{
			FirewallGroupName: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		foundFirewallGroup := res.FirewallGroup

		if nifcloud.ToString(foundFirewallGroup.Name) != saved.Primary.ID {
			return fmt.Errorf("Hatoba firewall group does not found in cloud: %s", saved.Primary.ID)
		}

		*firewallGroup = *foundFirewallGroup

		return nil
	}
}

func testAccCheckHatobaFirewallGroupValues(firewallGroup *types.FirewallGroup, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(firewallGroup.Name) != name {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", name, nifcloud.ToString(firewallGroup.Name))
		}

		if nifcloud.ToString(firewallGroup.Description) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", nifcloud.ToString(firewallGroup.Description))
		}

		wantRules := []types.Rules{
			{
				Protocol:    nifcloud.String("TCP"),
				Direction:   nifcloud.String("IN"),
				FromPort:    nifcloud.Int32(80),
				ToPort:      nifcloud.Int32(80),
				CidrIp:      nifcloud.String("0.0.0.0/0"),
				Description: nifcloud.String("rule memo"),
			},
			{
				Protocol:  nifcloud.String("TCP"),
				Direction: nifcloud.String("IN"),
				FromPort:  nifcloud.Int32(443),
				ToPort:    nifcloud.Int32(443),
				CidrIp:    nifcloud.String("0.0.0.0/0"),
			},
			{
				Protocol:  nifcloud.String("TCP"),
				Direction: nifcloud.String("OUT"),
				FromPort:  nifcloud.Int32(53),
				ToPort:    nifcloud.Int32(53),
				CidrIp:    nifcloud.String("8.8.8.8"),
			},
			{
				Protocol:  nifcloud.String("UDP"),
				Direction: nifcloud.String("OUT"),
				FromPort:  nifcloud.Int32(53),
				ToPort:    nifcloud.Int32(53),
				CidrIp:    nifcloud.String("8.8.8.8"),
			},
		}

		if len(firewallGroup.Rules) != len(wantRules) {
			return fmt.Errorf("bad rule[*] state, expected length is %d, got length: %d", len(wantRules), len(firewallGroup.Rules))
		}

		found := map[string]types.Rules{}
		for _, gr := range firewallGroup.Rules {
			for _, wr := range wantRules {
				if nifcloud.ToString(gr.Protocol) == nifcloud.ToString(wr.Protocol) &&
					nifcloud.ToString(gr.Direction) == nifcloud.ToString(wr.Direction) &&
					nifcloud.ToInt32(gr.FromPort) == nifcloud.ToInt32(wr.FromPort) &&
					nifcloud.ToInt32(gr.ToPort) == nifcloud.ToInt32(wr.ToPort) &&
					nifcloud.ToString(gr.CidrIp) == nifcloud.ToString(wr.CidrIp) &&
					nifcloud.ToString(gr.Description) == nifcloud.ToString(wr.Description) {
					found[nifcloud.ToString(gr.Id)] = gr
					break
				}
			}
		}

		if len(found) != len(wantRules) {
			return fmt.Errorf("bad rule[*] state, expected rules not found in cloud")
		}

		return nil
	}
}

func testAccCheckHatobaFirewallGroupValuesUpdated(firewallGroup *types.FirewallGroup, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(firewallGroup.Name) != name+"upd" {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", name+"upd", nifcloud.ToString(firewallGroup.Name))
		}

		if nifcloud.ToString(firewallGroup.Description) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", nifcloud.ToString(firewallGroup.Description))
		}

		wantRules := []types.Rules{
			{
				Protocol:    nifcloud.String("TCP"),
				Direction:   nifcloud.String("IN"),
				FromPort:    nifcloud.Int32(443),
				ToPort:      nifcloud.Int32(443),
				CidrIp:      nifcloud.String("0.0.0.0/0"),
				Description: nifcloud.String("HTTPS incomming"),
			},
			{
				Protocol:  nifcloud.String("ANY"),
				Direction: nifcloud.String("IN"),
				CidrIp:    nifcloud.String("192.168.0.0/24"),
			},
			{
				Protocol:  nifcloud.String("ANY"),
				Direction: nifcloud.String("OUT"),
				CidrIp:    nifcloud.String("0.0.0.0/0"),
			},
		}

		if len(firewallGroup.Rules) != len(wantRules) {
			return fmt.Errorf("bad rule[*] state, expected length is %d, got length: %d", len(wantRules), len(firewallGroup.Rules))
		}

		found := map[string]types.Rules{}
		for _, gr := range firewallGroup.Rules {
			for _, wr := range wantRules {
				if nifcloud.ToString(gr.Protocol) == nifcloud.ToString(wr.Protocol) &&
					nifcloud.ToString(gr.Direction) == nifcloud.ToString(wr.Direction) &&
					nifcloud.ToInt32(gr.FromPort) == nifcloud.ToInt32(wr.FromPort) &&
					nifcloud.ToInt32(gr.ToPort) == nifcloud.ToInt32(wr.ToPort) &&
					nifcloud.ToString(gr.CidrIp) == nifcloud.ToString(wr.CidrIp) &&
					nifcloud.ToString(gr.Description) == nifcloud.ToString(wr.Description) {
					found[nifcloud.ToString(gr.Id)] = gr
					break
				}
			}
		}

		if len(found) != len(wantRules) {
			return fmt.Errorf("bad rule[*] state, expected rules not found in cloud")
		}

		return nil
	}
}

func testAccHatobaFirewallGroupResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Hatoba

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_hatoba_firewall_group" {
			continue
		}

		_, err := svc.GetFirewallGroup(context.Background(), &hatoba.GetFirewallGroupInput{
			FirewallGroupName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.FirewallGroup" {
				return nil
			}
			return fmt.Errorf("failed GetFirewallGroupRequest: %s", err)
		}
	}
	return nil
}

func testSweepHatobaFirewallGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Hatoba

	res, err := svc.ListFirewallGroups(ctx, nil)
	if err != nil {
		return err
	}

	var sweepHatobaFirewallGroups []string
	for _, k := range res.FirewallGroups {
		if strings.HasPrefix(nifcloud.ToString(k.Name), prefix) {
			sweepHatobaFirewallGroups = append(sweepHatobaFirewallGroups, nifcloud.ToString(k.Name))
		}
	}

	if _, err := svc.DeleteFirewallGroups(ctx, &hatoba.DeleteFirewallGroupsInput{
		Names: nifcloud.String(strings.Join(sweepHatobaFirewallGroups, ",")),
	}); err != nil {
		return err
	}

	return nil
}
