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

var devopsGitLabURL = os.Getenv("TF_VAR_devops_gitlab_url")
var devopsRunnerToken = os.Getenv("TF_VAR_devops_runner_token")

func init() {
	resource.AddTestSweepers("nifcloud_devops_runner_registration", &resource.Sweeper{
		Name: "nifcloud_devops_runner_registration",
		F:    testSweepDevOpsRunnerRegistration,
	})
}

func TestAcc_DevOpsRunnerRegistration(t *testing.T) {
	var registration types.Registrations

	resourceName := "nifcloud_devops_runner_registration.basic"
	randName := prefix + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDevOpsRunnerRegistrationResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevOpsRunnerRegistration(t, "testdata/devops_runner_registration.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsRunnerRegistrationExists(resourceName, &registration),
					testAccCheckDevOpsRunnerRegistrationValues(&registration, randName),
					resource.TestCheckResourceAttr(resourceName, "runner_name", randName),
					resource.TestCheckResourceAttr(resourceName, "gitlab_url", devopsGitLabURL),
					resource.TestCheckResourceAttr(resourceName, "parameter_group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "token", devopsRunnerToken[5:14]),
				),
			},
			{
				Config: testAccDevOpsRunnerRegistration(t, "testdata/devops_runner_registration_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsRunnerRegistrationExists(resourceName, &registration),
					testAccCheckDevOpsRunnerRegistrationValuesUpdated(&registration, randName),
					resource.TestCheckResourceAttr(resourceName, "runner_name", randName),
					resource.TestCheckResourceAttr(resourceName, "gitlab_url", devopsGitLabURL),
					resource.TestCheckResourceAttr(resourceName, "parameter_group_name", randName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "token", devopsRunnerToken[5:14]),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDevOpsRunnerRegistrationImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDevOpsRunnerRegistration(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), rName, rName, rName)
}

func testAccCheckDevOpsRunnerRegistrationExists(n string, registration *types.Registrations) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no devops runner resource: %s", n)
		}

		if saved.Primary.Attributes["runner_name"] == "" {
			return fmt.Errorf("no devops runner id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DevOpsRunner
		res, err := svc.ListRunnerRegistrations(context.Background(), &devopsrunner.ListRunnerRegistrationsInput{
			RunnerName: nifcloud.String(saved.Primary.Attributes["runner_name"]),
		})
		if err != nil {
			return err
		}

		if res.Registrations == nil {
			return fmt.Errorf("devops runner registrations are not found in cloud: %s", saved.Primary.ID)
		}

		if nifcloud.ToString(res.Registrations[0].RegistrationId) != saved.Primary.ID {
			return fmt.Errorf("devops runner registration is not found in cloud: %s", saved.Primary.ID)
		}

		*registration = res.Registrations[0]

		return nil
	}
}

func testAccCheckDevOpsRunnerRegistrationValues(registration *types.Registrations, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(registration.GitlabUrl) != devopsGitLabURL {
			return fmt.Errorf("bad gitlab_url state, expected \"%s\", got: %#v", devopsGitLabURL, nifcloud.ToString(registration.GitlabUrl))
		}

		if nifcloud.ToString(registration.ParameterGroupName) != rName {
			return fmt.Errorf("bad parameter_group_name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(registration.ParameterGroupName))
		}

		if nifcloud.ToString(registration.Token) != devopsRunnerToken[5:14] {
			return fmt.Errorf("bad token state")
		}

		return nil
	}
}

func testAccCheckDevOpsRunnerRegistrationValuesUpdated(registration *types.Registrations, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(registration.GitlabUrl) != devopsGitLabURL {
			return fmt.Errorf("bad gitlab_url state, expected \"%s\", got: %#v", devopsGitLabURL, nifcloud.ToString(registration.GitlabUrl))
		}

		if nifcloud.ToString(registration.ParameterGroupName) != rName+"-upd" {
			return fmt.Errorf("bad parameter_group_name state, expected \"%s\", got: %#v", rName+"-upd", nifcloud.ToString(registration.ParameterGroupName))
		}

		if nifcloud.ToString(registration.Token) != devopsRunnerToken[5:14] {
			return fmt.Errorf("bad token state")
		}

		return nil
	}
}

func testAccDevOpsRunnerRegistrationResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).DevOpsRunner

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_devops_runner_registration" {
			continue
		}

		_, err := svc.GetRunner(context.Background(), &devopsrunner.GetRunnerInput{
			RunnerName: nifcloud.String(rs.Primary.Attributes["runner_name"]),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Runner" {
				return nil
			}
			return fmt.Errorf("failed GetRunner: %s", err)
		}

		return fmt.Errorf("devops runner (%s) still exists", rs.Primary.Attributes["runner_name"])
	}
	return nil
}

func testSweepDevOpsRunnerRegistration(region string) error {
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

func testAccDevOpsRunnerRegistrationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		runnerName := rs.Primary.Attributes["runner_name"]
		registrationId := rs.Primary.ID

		var parts []string
		parts = append(parts, runnerName)
		parts = append(parts, registrationId)

		id := strings.Join(parts, "_")
		return id, nil
	}
}
