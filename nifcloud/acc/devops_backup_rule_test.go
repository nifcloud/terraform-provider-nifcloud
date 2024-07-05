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
	resource.AddTestSweepers("nifcloud_devops_backup_rule", &resource.Sweeper{
		Name: "nifcloud_devops_backup_rule",
		F:    testSweepDevOpsBackupRule,
	})
}

func TestAcc_DevOpsBackupRule(t *testing.T) {
	var rule types.BackupRule

	resourceName := "nifcloud_devops_backup_rule.basic"
	randName := prefix + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDevOpsBackupRuleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevOpsBackupRule(t, "testdata/devops_backup_rule.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsBackupRuleExists(resourceName, &rule),
					testAccCheckDevOpsBackupRuleValues(&rule, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "instance_id", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
				),
			},
			{
				Config: testAccDevOpsBackupRule(t, "testdata/devops_backup_rule_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsBackupRuleExists(resourceName, &rule),
					testAccCheckDevOpsBackupRuleValuesUpdated(&rule, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo-upd"),
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

func testAccDevOpsBackupRule(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), rName, rName, rName, rName)
}

func testAccCheckDevOpsBackupRuleExists(n string, rule *types.BackupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no devops backup rule resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no devops backup rule id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DevOps
		res, err := svc.GetBackupRule(context.Background(), &devops.GetBackupRuleInput{
			BackupRuleName: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		if res.BackupRule == nil {
			return fmt.Errorf("devops backup rule is not found in cloud: %s", saved.Primary.ID)
		}

		if nifcloud.ToString(res.BackupRule.BackupRuleName) != saved.Primary.ID {
			return fmt.Errorf("devops backup rule is not found in cloud: %s", saved.Primary.ID)
		}

		*rule = *res.BackupRule

		return nil
	}
}

func testAccCheckDevOpsBackupRuleValues(rule *types.BackupRule, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(rule.BackupRuleName) != rName {
			return fmt.Errorf("bad backup rule name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(rule.BackupRuleName))
		}

		if nifcloud.ToString(rule.InstanceId) != rName {
			return fmt.Errorf("bad instance_id state, expected \"%s\", got: %#v", rName, nifcloud.ToString(rule.InstanceId))
		}

		if nifcloud.ToString(rule.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", nifcloud.ToString(rule.Description))
		}

		return nil
	}
}

func testAccCheckDevOpsBackupRuleValuesUpdated(rule *types.BackupRule, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(rule.BackupRuleName) != rName+"-upd" {
			return fmt.Errorf("bad backup rule name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(rule.BackupRuleName))
		}

		if nifcloud.ToString(rule.InstanceId) != rName {
			return fmt.Errorf("bad instance_id state, expected \"%s\", got: %#v", rName, nifcloud.ToString(rule.InstanceId))
		}

		if nifcloud.ToString(rule.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", nifcloud.ToString(rule.Description))
		}

		return nil
	}
}

func testAccDevOpsBackupRuleResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).DevOps

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_devops_backup_rule" {
			continue
		}

		_, err := svc.GetBackupRule(context.Background(), &devops.GetBackupRuleInput{
			BackupRuleName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.BackupRule" {
				return nil
			}
			return fmt.Errorf("failed GetBackupRule: %s", err)
		}

		return fmt.Errorf("devops backup rule (%s) still exists", rs.Primary.ID)
	}
	return nil
}

func testSweepDevOpsBackupRule(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).DevOps

	res, err := svc.ListBackupRules(ctx, nil)
	if err != nil {
		return err
	}

	var sweepBackupRules []string
	for _, g := range res.BackupRules {
		if strings.HasPrefix(nifcloud.ToString(g.BackupRuleName), prefix) {
			sweepBackupRules = append(sweepBackupRules, nifcloud.ToString(g.BackupRuleName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepBackupRules {
		group := n
		eg.Go(func() error {
			_, err := svc.DeleteBackupRule(ctx, &devops.DeleteBackupRuleInput{
				BackupRuleName: nifcloud.String(group),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
