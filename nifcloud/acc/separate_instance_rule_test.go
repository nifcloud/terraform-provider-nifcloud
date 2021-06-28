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
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
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
	var separateInstanceRules computing.SeparateInstanceRulesInfo

	resourceName := "nifcloud_separate_instance_rule.basic"
	randName := prefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

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
					//Terraform上のステートの値がTerrafromConfigで指定した値の通りになっているか確認
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-11"),
					resource.TestCheckResourceAttr(resourceName, "instance_id.0", randName),
					resource.TestCheckResourceAttr(resourceName, "instance_id.1", randName),
				),
			},
			{
				Config: testAccSeparateInstanceRule(t, "testdata/separate_instance_rule_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSeparateInstanceRuleExists(resourceName, &separateInstanceRules),
					testAccCheckSeparateInstanceRuleValuesUpdated(&separateInstanceRules, randName),
					//TerrafromConfigを更新した場合、Terraform上のステートの値がTerrafromConfigで指定した値の通りになっているか確認
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-11"),
					resource.TestCheckResourceAttr(resourceName, "instance_id.0", randName),
					resource.TestCheckResourceAttr(resourceName, "instance_id.1", randName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"reboot",
				},
			},
		},
	})
}

func testAccSeparateInstanceRule(t *testing.T, fileName, SeparateInstanceRuleName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		SeparateInstanceRuleName,
	)
}

//作成されたリソースが実際のクラウド上に存在するかの確認
func testAccCheckSeparateInstanceRuleExists(n string, separateInstanceRules *computing.SeparateInstanceRulesInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no SeparateInstanceRule resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no SeparateInstanceRule id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.NiftyDescribeSeparateInstanceRulesRequest(&computing.NiftyDescribeSeparateInstanceRulesInput{
			SeparateInstanceRuleName: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.SeparateInstanceRulesInfo) == 0 {
			return fmt.Errorf("SeparateInstanceRule does not found in cloud: %s", saved.Primary.ID)
		}

		foundSeparateInstanceRule := res.SeparateInstanceRulesInfo[0]

		if nifcloud.StringValue(foundSeparateInstanceRule.SeparateInstanceRuleName) != saved.Primary.ID {
			return fmt.Errorf("SeparateInstanceRule does not found in cloud: %s", saved.Primary.ID)
		}

		*separateInstanceRules = foundSeparateInstanceRule
		return nil
	}
}

//クラウドに作成されたリソースの値がTerrafromConfigで指定した値の通りになっているか確認
func testAccCheckSeparateInstanceRuleValues(separateInstanceRules *computing.SeparateInstanceRulesInfo, SeparateInstanceRuleName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(separateInstanceRules.SeparateInstanceRuleName) != SeparateInstanceRuleName {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", SeparateInstanceRuleName, separateInstanceRules.SeparateInstanceRuleName)
		}

		if nifcloud.StringValue(separateInstanceRules.SeparateInstanceRuleDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", separateInstanceRules.SeparateInstanceRuleDescription)
		}

		if nifcloud.StringValue(separateInstanceRules.AvailabilityZone) != "east-11" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-11\", got: %#v", separateInstanceRules.AvailabilityZone)
		}

		/*if nifcloud.StringValue(volume.AttachmentSet[0].InstanceId) != rName {
			return fmt.Errorf("bad instance_id state, expected \"%s\", got: %#v", rName, volume.AttachmentSet[0].InstanceId)
		}*/
		return nil
	}
}

//TerrafromConfigを更新した場合、クラウド上のリソースの値も更新されていることの確認
func testAccCheckSeparateInstanceRuleValuesUpdated(separateInstanceRules *computing.SeparateInstanceRulesInfo, SeparateInstanceRuleName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(separateInstanceRules.SeparateInstanceRuleName) != SeparateInstanceRuleName {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", SeparateInstanceRuleName, separateInstanceRules.SeparateInstanceRuleName)
		}

		if nifcloud.StringValue(separateInstanceRules.SeparateInstanceRuleDescription) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", separateInstanceRules.SeparateInstanceRuleDescription)
		}

		if nifcloud.StringValue(separateInstanceRules.AvailabilityZone) != "east-11" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-11\", got: %#v", separateInstanceRules.AvailabilityZone)
		}
		return nil
	}
}

//テスト終了後にクラウド上のリソースが削除されているか確認
func testAccSeparateInstanceRuleResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_separate_instance_rule" {
			continue
		}

		res, err := svc.NiftyDescribeSeparateInstanceRulesRequest(&computing.NiftyDescribeSeparateInstanceRulesInput{
			SeparateInstanceRuleName: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.SeparateInstanceRule" {
				return fmt.Errorf("failed NiftyDescribeSeparateInstanceRulesRequest: %s", err)
			}
		}

		if len(res.SeparateInstanceRulesInfo) > 0 {
			return fmt.Errorf("SeparateInstanceRule (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

//テスト実施前にクラウド上リソースを掃除する用の関数
func testSweepSeparateInstanceRule(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.NiftyDescribeSeparateInstanceRulesRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepSeparateInstanceRules []string
	for _, k := range res.SeparateInstanceRulesInfo {
		if strings.HasPrefix(nifcloud.StringValue(k.SeparateInstanceRuleName), prefix) {
			sweepSeparateInstanceRules = append(sweepSeparateInstanceRules, nifcloud.StringValue(k.SeparateInstanceRuleName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepSeparateInstanceRules {
		SeparateInstanceRuleName := n
		eg.Go(func() error {
			_, err := svc.NiftyDeleteSeparateInstanceRuleRequest(&computing.NiftyDeleteSeparateInstanceRuleInput{
				SeparateInstanceRuleName: nifcloud.String(SeparateInstanceRuleName),
			}).Send(ctx)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
