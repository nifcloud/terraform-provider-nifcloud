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
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_devops_firewall_group", &resource.Sweeper{
		Name: "nifcloud_devops_firewall_group",
		F:    testSweepDevOpsFirewallGroup,
		// Dependencies: []string{
		// 	"nifcloud_devops_instance",
		// },
	})
}

func TestAcc_DevOpsFirewallGroup(t *testing.T) {
	var group types.FirewallGroup

	resourceName := "nifcloud_devops_firewall_group.basic"
	randName := prefix + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDevOpsFirewallGroupResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevOpsFirewallGroup(t, "testdata/devops_firewall_group.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsFirewallGroupExists(resourceName, &group),
					testAccCheckDevOpsFirewallGroupValues(&group, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-14"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"rule.*",
						map[string]string{
							"protocol": "TCP",
							"port":     "443",
							"cidr_ip":  "172.16.0.0/24",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"rule.*",
						map[string]string{
							"protocol": "TCP",
							"port":     "22",
							"cidr_ip":  "172.16.0.0/24",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"rule.*",
						map[string]string{
							"protocol":    "ICMP",
							"cidr_ip":     "172.16.0.0/24",
							"description": "ping",
						},
					),
				),
			},
			{
				Config: testAccDevOpsFirewallGroup(t, "testdata/devops_firewall_group_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsFirewallGroupExists(resourceName, &group),
					testAccCheckDevOpsFirewallGroupValuesUpdated(&group, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-14"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol", "TCP"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"rule.*",
						map[string]string{
							"protocol": "TCP",
							"port":     "443",
							"cidr_ip":  "192.168.1.0/24",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"rule.*",
						map[string]string{
							"protocol": "TCP",
							"port":     "22",
							"cidr_ip":  "192.168.1.0/24",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"rule.*",
						map[string]string{
							"protocol":    "ICMP",
							"cidr_ip":     "192.168.1.0/24",
							"description": "pong",
						},
					),
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

func testAccDevOpsFirewallGroup(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), rName)
}

func testAccCheckDevOpsFirewallGroupExists(n string, group *types.FirewallGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no devops firewall group resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no devops firewall group id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DevOps
		res, err := svc.GetFirewallGroup(context.Background(), &devops.GetFirewallGroupInput{
			FirewallGroupName: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		if res.FirewallGroup == nil {
			return fmt.Errorf("devops firewall group is not found in cloud: %s", saved.Primary.ID)
		}

		if nifcloud.ToString(res.FirewallGroup.FirewallGroupName) != saved.Primary.ID {
			return fmt.Errorf("devops firewall group is not found in cloud: %s", saved.Primary.ID)
		}

		*group = *res.FirewallGroup

		return nil
	}
}

func testAccCheckDevOpsFirewallGroupValues(group *types.FirewallGroup, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(group.FirewallGroupName) != rName {
			return fmt.Errorf("bad firewall group name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(group.FirewallGroupName))
		}

		if nifcloud.ToString(group.AvailabilityZone) != "east-14" {
			return fmt.Errorf("bad availability_zone state, expected \"east-14\", got: %#v", nifcloud.ToString(group.AvailabilityZone))
		}

		if nifcloud.ToString(group.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", nifcloud.ToString(group.Description))
		}

		if len(group.Rules) != 3 {
			return fmt.Errorf("bad rules length: %#v", group.Rules)
		}

		var httpOk bool
		var sshOk bool
		var icmpOk bool

		for _, r := range group.Rules {
			if nifcloud.ToString(r.Protocol) == "TCP" && nifcloud.ToInt32(r.Port) == int32(443) && nifcloud.ToString(r.CidrIp) == "172.16.0.0/24" {
				httpOk = true
			} else if nifcloud.ToString(r.Protocol) == "TCP" && nifcloud.ToInt32(r.Port) == int32(22) && nifcloud.ToString(r.CidrIp) == "172.16.0.0/24" {
				sshOk = true
			} else if nifcloud.ToString(r.Protocol) == "ICMP" && nifcloud.ToString(r.CidrIp) == "172.16.0.0/24" && nifcloud.ToString(r.Description) == "ping" {
				icmpOk = true
			}
		}

		if !httpOk || !sshOk || !icmpOk {
			return fmt.Errorf("one or more rule(s) not found: %#v", group.Rules)
		}

		return nil
	}
}

func testAccCheckDevOpsFirewallGroupValuesUpdated(group *types.FirewallGroup, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(group.FirewallGroupName) != rName+"-upd" {
			return fmt.Errorf("bad firewall group name state, expected \"%s\", got: %#v", rName+"-upd", nifcloud.ToString(group.FirewallGroupName))
		}

		if nifcloud.ToString(group.AvailabilityZone) != "east-14" {
			return fmt.Errorf("bad availability_zone state, expected \"east-14\", got: %#v", nifcloud.ToString(group.AvailabilityZone))
		}

		if nifcloud.ToString(group.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", nifcloud.ToString(group.Description))
		}

		var httpOk bool
		var sshOk bool
		var icmpOk bool

		for _, r := range group.Rules {
			if nifcloud.ToString(r.Protocol) == "TCP" && nifcloud.ToInt32(r.Port) == int32(443) && nifcloud.ToString(r.CidrIp) == "192.168.1.0/24" {
				httpOk = true
			} else if nifcloud.ToString(r.Protocol) == "TCP" && nifcloud.ToInt32(r.Port) == int32(22) && nifcloud.ToString(r.CidrIp) == "192.168.1.0/24" {
				sshOk = true
			} else if nifcloud.ToString(r.Protocol) == "ICMP" && nifcloud.ToString(r.CidrIp) == "192.168.1.0/24" && nifcloud.ToString(r.Description) == "pong" {
				icmpOk = true
			}
		}

		if !httpOk || !sshOk || !icmpOk {
			return fmt.Errorf("one or more rule(s) not found: %#v", group.Rules)
		}

		return nil
	}
}

func testAccDevOpsFirewallGroupResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).DevOps

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_devops_firewall_group" {
			continue
		}

		_, err := svc.GetFirewallGroup(context.Background(), &devops.GetFirewallGroupInput{
			FirewallGroupName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.FirewallGroup" {
				return nil
			}
			return fmt.Errorf("failed GetFirewallGroup: %s", err)
		}

		return fmt.Errorf("devops firewall group (%s) still exists", rs.Primary.ID)
	}
	return nil
}

func testSweepDevOpsFirewallGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).DevOps

	res, err := svc.ListFirewallGroups(ctx, nil)
	if err != nil {
		return err
	}

	var sweepFirewallGroups []string
	for _, g := range res.FirewallGroups {
		if strings.HasPrefix(nifcloud.ToString(g.FirewallGroupName), prefix) {
			sweepFirewallGroups = append(sweepFirewallGroups, nifcloud.ToString(g.FirewallGroupName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepFirewallGroups {
		group := n
		eg.Go(func() error {
			_, err := svc.DeleteFirewallGroup(ctx, &devops.DeleteFirewallGroupInput{
				FirewallGroupName: nifcloud.String(group),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
