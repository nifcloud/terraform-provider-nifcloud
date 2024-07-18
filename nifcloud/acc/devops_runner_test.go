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
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_devops_runner", &resource.Sweeper{
		Name: "nifcloud_devops_runner",
		F:    testSweepDevOpsRunner,
	})
}

func TestAcc_DevOpsRunner(t *testing.T) {
	var runner types.Runner

	resourceName := "nifcloud_devops_runner.basic"
	randName := prefix + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDevOpsRunnerResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevOpsRunner(t, "testdata/devops_runner.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsRunnerExists(resourceName, &runner),
					testAccCheckDevOpsRunnerValues(&runner, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "c-small"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-14"),
					resource.TestCheckResourceAttr(resourceName, "concurrent", "10"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
				),
			},
			{
				Config: testAccDevOpsRunner(t, "testdata/devops_runner_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsRunnerExists(resourceName, &runner),
					testAccCheckDevOpsRunnerValuesUpdated(&runner, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "e-small"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-14"),
					resource.TestCheckResourceAttr(resourceName, "concurrent", "50"),
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

func testAccDevOpsRunner(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), rName)
}

func testAccCheckDevOpsRunnerExists(n string, runner *types.Runner) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no devops runner resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no devops runner id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DevOpsRunner
		res, err := svc.GetRunner(context.Background(), &devopsrunner.GetRunnerInput{
			RunnerName: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		if res.Runner == nil {
			return fmt.Errorf("devops runner is not found in cloud: %s", saved.Primary.ID)
		}

		if nifcloud.ToString(res.Runner.RunnerName) != saved.Primary.ID {
			return fmt.Errorf("devops runner is not found in cloud: %s", saved.Primary.ID)
		}

		*runner = *res.Runner

		return nil
	}
}

func testAccCheckDevOpsRunnerValues(runner *types.Runner, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(runner.RunnerName) != rName {
			return fmt.Errorf("bad runner name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(runner.RunnerName))
		}

		if nifcloud.ToString(runner.InstanceType) != "c-small" {
			return fmt.Errorf("bad instance_type state, expected \"c-small\", got: %#v", nifcloud.ToString(runner.InstanceType))
		}

		if nifcloud.ToString(runner.AvailabilityZone) != "east-14" {
			return fmt.Errorf("bad availability_zone state, expected \"east-14\", got: %#v", nifcloud.ToString(runner.AvailabilityZone))
		}

		if nifcloud.ToInt32(runner.Concurrent) != int32(10) {
			return fmt.Errorf("bad concurrent state, expected 10, got: %#v", nifcloud.ToInt32(runner.Concurrent))
		}

		if nifcloud.ToString(runner.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", nifcloud.ToString(runner.Description))
		}

		return nil
	}
}

func testAccCheckDevOpsRunnerValuesUpdated(runner *types.Runner, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(runner.RunnerName) != rName+"-upd" {
			return fmt.Errorf("bad runner name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(runner.RunnerName))
		}

		if nifcloud.ToString(runner.InstanceType) != "e-small" {
			return fmt.Errorf("bad instance_type state, expected \"e-small\", got: %#v", nifcloud.ToString(runner.InstanceType))
		}

		if nifcloud.ToString(runner.AvailabilityZone) != "east-14" {
			return fmt.Errorf("bad availability_zone state, expected \"east-14\", got: %#v", nifcloud.ToString(runner.AvailabilityZone))
		}

		if nifcloud.ToInt32(runner.Concurrent) != int32(50) {
			return fmt.Errorf("bad concurrent state, expected 50, got: %#v", nifcloud.ToInt32(runner.Concurrent))
		}

		if nifcloud.ToString(runner.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", nifcloud.ToString(runner.Description))
		}

		return nil
	}
}

func testAccDevOpsRunnerResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).DevOpsRunner

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_devops_runner" {
			continue
		}

		_, err := svc.GetRunner(context.Background(), &devopsrunner.GetRunnerInput{
			RunnerName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Runner" {
				return nil
			}
			return fmt.Errorf("failed GetRunner: %s", err)
		}

		return fmt.Errorf("devops runner (%s) still exists", rs.Primary.ID)
	}
	return nil
}

func testSweepDevOpsRunner(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).DevOpsRunner

	res, err := svc.ListRunners(ctx, nil)
	if err != nil {
		return err
	}

	var sweepRunners []string
	for _, r := range res.Runners {
		if strings.HasPrefix(nifcloud.ToString(r.RunnerName), prefix) {
			sweepRunners = append(sweepRunners, nifcloud.ToString(r.RunnerName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepRunners {
		runner := n
		eg.Go(func() error {
			_, err := svc.DeleteRunner(ctx, &devopsrunner.DeleteRunnerInput{
				RunnerName: nifcloud.String(runner),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
