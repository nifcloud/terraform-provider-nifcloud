package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_separate_instance_rule", &resource.Sweeper{
		Name: "nifcloud_separate_instance_rule",
		F:    testSweepSeparateInstanceRule,
	})
}

func TestAcc_SeparateInstanceRule(t *testing.T) {
	var separateInstanceRules types.SeparateInstanceRulesInfo

	resourceName := "nifcloud_separate_instance_rule.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccSeparateInstanceRuleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSeparateInstanceRule(t, "testdata/separate_instance_rule.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSeparateInstanceRuleExists(resourceName, &separateInstanceRules),
					testAccCheckSeparateInstanceRuleValues(&separateInstanceRules, randName),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id.0"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id.1"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
				),
			},
			{
				Config: testAccSeparateInstanceRule(t, "testdata/separate_instance_rule_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSeparateInstanceRuleExists(resourceName, &separateInstanceRules),
					testAccCheckSeparateInstanceRuleValuesUpdated(&separateInstanceRules, randName),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id.0"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id.1"),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
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

func TestAcc_SeparateInstanceRule_Unique_Id(t *testing.T) {
	var separateInstanceRules types.SeparateInstanceRulesInfo

	resourceName := "nifcloud_separate_instance_rule.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccSeparateInstanceRuleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSeparateInstanceRule(t, "testdata/separate_instance_rule_unique_id.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSeparateInstanceRuleExists(resourceName, &separateInstanceRules),
					testAccCheckSeparateInstanceRuleUniqueIDValues(&separateInstanceRules, randName),
					resource.TestCheckResourceAttrSet(resourceName, "instance_unique_id.0"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_unique_id.1"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id.",
					"instance_unique_id.",
				},
			},
		},
	})
}

func testAccSeparateInstanceRule(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
		rName,
		rName,
	)
}

func testAccCheckSeparateInstanceRuleExists(n string, separateInstanceRules *types.SeparateInstanceRulesInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no SeparateInstanceRule resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no SeparateInstanceRule id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.NiftyDescribeSeparateInstanceRules(context.Background(), &computing.NiftyDescribeSeparateInstanceRulesInput{
			SeparateInstanceRuleName: []string{saved.Primary.ID},
		})

		if err != nil {
			return err
		}

		if len(res.SeparateInstanceRulesInfo) == 0 {
			return fmt.Errorf("SeparateInstanceRule does not found in cloud: %s", saved.Primary.ID)
		}

		foundSeparateInstanceRule := res.SeparateInstanceRulesInfo[0]

		if nifcloud.ToString(foundSeparateInstanceRule.SeparateInstanceRuleName) != saved.Primary.ID {
			return fmt.Errorf("SeparateInstanceRule does not found in cloud: %s", saved.Primary.ID)
		}

		*separateInstanceRules = foundSeparateInstanceRule
		return nil
	}
}

func testAccCheckSeparateInstanceRuleValues(separateInstanceRules *types.SeparateInstanceRulesInfo, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(separateInstanceRules.SeparateInstanceRuleName) != rName {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", rName, separateInstanceRules.SeparateInstanceRuleName)
		}

		if nifcloud.ToString(separateInstanceRules.SeparateInstanceRuleDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", separateInstanceRules.SeparateInstanceRuleDescription)
		}

		if nifcloud.ToString(separateInstanceRules.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", separateInstanceRules.AvailabilityZone)
		}

		if nifcloud.ToString(separateInstanceRules.InstancesSet[0].InstanceId) == "" {
			return fmt.Errorf("bad instance_id state, expected not nil, got: nil")
		}

		if nifcloud.ToString(separateInstanceRules.InstancesSet[1].InstanceId) == "" {
			return fmt.Errorf("bad instance_id state, expected not nil, got: nil")
		}
		return nil
	}
}

func testAccCheckSeparateInstanceRuleValuesUpdated(separateInstanceRules *types.SeparateInstanceRulesInfo, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(separateInstanceRules.SeparateInstanceRuleName) != rName+"upd" {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", rName+"upd", separateInstanceRules.SeparateInstanceRuleName)
		}

		if nifcloud.ToString(separateInstanceRules.SeparateInstanceRuleDescription) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", separateInstanceRules.SeparateInstanceRuleDescription)
		}

		if nifcloud.ToString(separateInstanceRules.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", separateInstanceRules.AvailabilityZone)
		}

		if nifcloud.ToString(separateInstanceRules.InstancesSet[0].InstanceId) == "" {
			return fmt.Errorf("bad instance_id state, expected not nil, got: nil")
		}

		if nifcloud.ToString(separateInstanceRules.InstancesSet[1].InstanceId) == "" {
			return fmt.Errorf("bad instance_id state, expected not nil, got: nil")
		}
		return nil
	}
}

func testAccCheckSeparateInstanceRuleUniqueIDValues(separateInstanceRules *types.SeparateInstanceRulesInfo, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(separateInstanceRules.SeparateInstanceRuleName) != rName {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", rName, separateInstanceRules.SeparateInstanceRuleName)
		}

		if nifcloud.ToString(separateInstanceRules.SeparateInstanceRuleDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", separateInstanceRules.SeparateInstanceRuleDescription)
		}

		if nifcloud.ToString(separateInstanceRules.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", separateInstanceRules.AvailabilityZone)
		}

		if nifcloud.ToString(separateInstanceRules.InstancesSet[0].InstanceUniqueId) == "" {
			return fmt.Errorf("bad instance_unique_id state, expected not nil, got: nil")
		}

		if nifcloud.ToString(separateInstanceRules.InstancesSet[1].InstanceUniqueId) == "" {
			return fmt.Errorf("bad instance_unique_id state, expected not nil, got: nil")
		}
		return nil
	}
}

func testAccSeparateInstanceRuleResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_separate_instance_rule" {
			continue
		}

		res, err := svc.NiftyDescribeSeparateInstanceRules(context.Background(), &computing.NiftyDescribeSeparateInstanceRulesInput{
			SeparateInstanceRuleName: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() != "Client.InvalidParameterNotFound.SeparateInstanceRule" {
				return fmt.Errorf("failed NiftyDescribeSeparateInstanceRulesRequest: %s", err)
			}
		}

		if len(res.SeparateInstanceRulesInfo) > 0 {
			return fmt.Errorf("SeparateInstanceRule (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepSeparateInstanceRule(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.NiftyDescribeSeparateInstanceRules(ctx, nil)
	if err != nil {
		return err
	}

	var sweepSeparateInstanceRules []string
	for _, k := range res.SeparateInstanceRulesInfo {
		if strings.HasPrefix(nifcloud.ToString(k.SeparateInstanceRuleName), prefix) {
			sweepSeparateInstanceRules = append(sweepSeparateInstanceRules, nifcloud.ToString(k.SeparateInstanceRuleName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepSeparateInstanceRules {
		separateInstanceRuleName := n
		eg.Go(func() error {
			_, err := svc.NiftyDeleteSeparateInstanceRule(ctx, &computing.NiftyDeleteSeparateInstanceRuleInput{
				SeparateInstanceRuleName: nifcloud.String(separateInstanceRuleName),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
