package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func init() {
	resource.AddTestSweepers("nifcloud_hatoba_firewall_group", &resource.Sweeper{
		Name:         "nifcloud_hatoba_firewall_group",
		F:            testSweepHatobaFirewallGroup,
		Dependencies: []string{},
	})
}

func TestAcc_HatobaFirewallGroup(t *testing.T) {
	var firewallGroup hatoba.FirewallGroupResponse

	resourceName := "nifcloud_hatoba_firewall_group.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

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
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccCheckHatobaFirewallGroupExists(n string, firewallGroup *hatoba.FirewallGroupResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Hatoba firewall group resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no Hatoba firewall group id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Hatoba
		res, err := svc.GetFirewallGroupRequest(&hatoba.GetFirewallGroupInput{
			FirewallGroupName: nifcloud.String(saved.Primary.ID),
		}).Send(context.Background())
		if err != nil {
			return err
		}

		foundFirewallGroup := res.FirewallGroup

		if nifcloud.StringValue(foundFirewallGroup.Name) != saved.Primary.ID {
			return fmt.Errorf("Hatoba firewall group does not found in cloud: %s", saved.Primary.ID)
		}

		*firewallGroup = *foundFirewallGroup

		return nil
	}
}

func testAccCheckHatobaFirewallGroupValues(firewallGroup *hatoba.FirewallGroupResponse, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(firewallGroup.Name) != name {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", name, nifcloud.StringValue(firewallGroup.Name))
		}

		if nifcloud.StringValue(firewallGroup.Description) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", nifcloud.StringValue(firewallGroup.Description))
		}

		wantRules := []hatoba.FirewallRule{
			{
				Protocol:    nifcloud.String("TCP"),
				Direction:   nifcloud.String("IN"),
				FromPort:    nifcloud.Int64(80),
				ToPort:      nifcloud.Int64(80),
				CidrIp:      nifcloud.String("0.0.0.0/0"),
				Description: nifcloud.String("rule memo"),
			},
			{
				Protocol:  nifcloud.String("TCP"),
				Direction: nifcloud.String("IN"),
				FromPort:  nifcloud.Int64(443),
				ToPort:    nifcloud.Int64(443),
				CidrIp:    nifcloud.String("0.0.0.0/0"),
			},
			{
				Protocol:  nifcloud.String("TCP"),
				Direction: nifcloud.String("OUT"),
				FromPort:  nifcloud.Int64(53),
				ToPort:    nifcloud.Int64(53),
				CidrIp:    nifcloud.String("8.8.8.8"),
			},
			{
				Protocol:  nifcloud.String("UDP"),
				Direction: nifcloud.String("OUT"),
				FromPort:  nifcloud.Int64(53),
				ToPort:    nifcloud.Int64(53),
				CidrIp:    nifcloud.String("8.8.8.8"),
			},
		}

		if len(firewallGroup.Rules) != len(wantRules) {
			return fmt.Errorf("bad rule[*] state, expected length is %d, got length: %d", len(wantRules), len(firewallGroup.Rules))
		}

		found := map[string]hatoba.FirewallRule{}
		for _, gr := range firewallGroup.Rules {
			for _, wr := range wantRules {
				if nifcloud.StringValue(gr.Protocol) == nifcloud.StringValue(wr.Protocol) &&
					nifcloud.StringValue(gr.Direction) == nifcloud.StringValue(wr.Direction) &&
					nifcloud.Int64Value(gr.FromPort) == nifcloud.Int64Value(wr.FromPort) &&
					nifcloud.Int64Value(gr.ToPort) == nifcloud.Int64Value(wr.ToPort) &&
					nifcloud.StringValue(gr.CidrIp) == nifcloud.StringValue(wr.CidrIp) &&
					nifcloud.StringValue(gr.Description) == nifcloud.StringValue(wr.Description) {
					found[nifcloud.StringValue(gr.Id)] = gr
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

func testAccCheckHatobaFirewallGroupValuesUpdated(firewallGroup *hatoba.FirewallGroupResponse, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(firewallGroup.Name) != name+"upd" {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", name+"upd", nifcloud.StringValue(firewallGroup.Name))
		}

		if nifcloud.StringValue(firewallGroup.Description) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", nifcloud.StringValue(firewallGroup.Description))
		}

		wantRules := []hatoba.FirewallRule{
			{
				Protocol:    nifcloud.String("TCP"),
				Direction:   nifcloud.String("IN"),
				FromPort:    nifcloud.Int64(443),
				ToPort:      nifcloud.Int64(443),
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

		found := map[string]hatoba.FirewallRule{}
		for _, gr := range firewallGroup.Rules {
			for _, wr := range wantRules {
				if nifcloud.StringValue(gr.Protocol) == nifcloud.StringValue(wr.Protocol) &&
					nifcloud.StringValue(gr.Direction) == nifcloud.StringValue(wr.Direction) &&
					nifcloud.Int64Value(gr.FromPort) == nifcloud.Int64Value(wr.FromPort) &&
					nifcloud.Int64Value(gr.ToPort) == nifcloud.Int64Value(wr.ToPort) &&
					nifcloud.StringValue(gr.CidrIp) == nifcloud.StringValue(wr.CidrIp) &&
					nifcloud.StringValue(gr.Description) == nifcloud.StringValue(wr.Description) {
					found[nifcloud.StringValue(gr.Id)] = gr
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

		_, err := svc.GetFirewallGroupRequest(&hatoba.GetFirewallGroupInput{
			FirewallGroupName: nifcloud.String(rs.Primary.ID),
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.FirewallGroup" {
				return fmt.Errorf("failed GetFirewallGroupRequest: %s", err)
			}
		}
	}
	return nil
}

func testSweepHatobaFirewallGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Hatoba

	res, err := svc.ListFirewallGroupsRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepHatobaFirewallGroups []string
	for _, k := range res.FirewallGroups {
		if strings.HasPrefix(nifcloud.StringValue(k.Name), prefix) {
			sweepHatobaFirewallGroups = append(sweepHatobaFirewallGroups, nifcloud.StringValue(k.Name))
		}
	}

	if _, err := svc.DeleteFirewallGroupsRequest(&hatoba.DeleteFirewallGroupsInput{
		Names: nifcloud.String(strings.Join(sweepHatobaFirewallGroups, ",")),
	}).Send(ctx); err != nil {
		return err
	}

	return nil
}
