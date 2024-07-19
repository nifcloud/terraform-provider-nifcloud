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
	resource.AddTestSweepers("nifcloud_devops_instance", &resource.Sweeper{
		Name: "nifcloud_devops_instance",
		F:    testSweepDevOpsInstance,
	})
}

func TestAcc_DevOpsInstance(t *testing.T) {
	var instance types.Instance

	resourceName := "nifcloud_devops_instance.basic"
	randName := prefix + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDevOpsInstanceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevOpsInstance(t, "testdata/devops_instance.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsInstanceExists(resourceName, &instance),
					testAccCheckDevOpsInstanceValues(&instance, randName),
					resource.TestCheckResourceAttr(resourceName, "instance_id", randName),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "c-large"),
					resource.TestCheckResourceAttr(resourceName, "firewall_group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "parameter_group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "disk_size", "100"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-14"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "to", "email@example.com"),
					resource.TestCheckResourceAttr(resourceName, "gitlab_url", "https://"+randName+".jp-east-1.gitlab.devops.nifcloud.com"),
					resource.TestCheckResourceAttr(resourceName, "registry_url", "https://registry-"+randName+".jp-east-1.gitlab.devops.nifcloud.com"),
				),
			},
			{
				Config: testAccDevOpsInstance(t, "testdata/devops_instance_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsInstanceExists(resourceName, &instance),
					testAccCheckDevOpsInstanceValuesUpdated(&instance, randName),
					resource.TestCheckResourceAttr(resourceName, "instance_id", randName),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "e-large"),
					resource.TestCheckResourceAttr(resourceName, "firewall_group_name", randName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "parameter_group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "disk_size", "300"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-14"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "to", "email-upd@example.com"),
					resource.TestCheckResourceAttr(resourceName, "gitlab_url", "https://"+randName+".jp-east-1.gitlab.devops.nifcloud.com"),
					resource.TestCheckResourceAttr(resourceName, "registry_url", "https://registry-"+randName+".jp-east-1.gitlab.devops.nifcloud.com"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"initial_root_password",
				},
			},
		},
	})
}

func testAccDevOpsInstance(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), rName, rName, rName, rName)
}

func testAccCheckDevOpsInstanceExists(n string, instance *types.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no devops instance resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no devops instance id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DevOps
		res, err := svc.GetInstance(context.Background(), &devops.GetInstanceInput{
			InstanceId: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		if res.Instance == nil {
			return fmt.Errorf("devops instance is not found in cloud: %s", saved.Primary.ID)
		}

		if nifcloud.ToString(res.Instance.InstanceId) != saved.Primary.ID {
			return fmt.Errorf("devops instance is not found in cloud: %s", saved.Primary.ID)
		}

		*instance = *res.Instance

		return nil
	}
}

func testAccCheckDevOpsInstanceValues(instance *types.Instance, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(instance.InstanceId) != rName {
			return fmt.Errorf("bad instance id state, expected \"%s\", got: %#v", rName, nifcloud.ToString(instance.InstanceId))
		}

		if nifcloud.ToString(instance.InstanceType) != "c-large" {
			return fmt.Errorf("bad instance_type state, expected \"c-large\", got: %#v", nifcloud.ToString(instance.InstanceType))
		}

		if nifcloud.ToString(instance.FirewallGroupName) != rName {
			return fmt.Errorf("bad firewall_group_name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(instance.FirewallGroupName))
		}

		if nifcloud.ToString(instance.ParameterGroupName) != rName {
			return fmt.Errorf("bad parameter_group_name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(instance.ParameterGroupName))
		}

		if nifcloud.ToInt32(instance.DiskSize) != int32(100) {
			return fmt.Errorf("bad disk_size state, expected 100, got: %#v", nifcloud.ToInt32(instance.DiskSize))
		}

		if nifcloud.ToString(instance.AvailabilityZone) != "east-14" {
			return fmt.Errorf("bad availability_zone state, expected \"east-14\", got: %#v", nifcloud.ToString(instance.AvailabilityZone))
		}

		if nifcloud.ToString(instance.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", nifcloud.ToString(instance.Description))
		}

		if nifcloud.ToString(instance.To) != "email@example.com" {
			return fmt.Errorf("bad to state, expected \"email@example.com\", got: %#v", nifcloud.ToString(instance.To))
		}

		if nifcloud.ToString(instance.GitlabUrl) != "https://"+rName+".jp-east-1.gitlab.devops.nifcloud.com" {
			return fmt.Errorf("bad gitlab_url state, expected \"https://%s.jp-east-1.gitlab.devops.nifcloud.com\", got: %#v", rName, nifcloud.ToString(instance.GitlabUrl))
		}

		if nifcloud.ToString(instance.RegistryUrl) != "https://registry-"+rName+".jp-east-1.gitlab.devops.nifcloud.com" {
			return fmt.Errorf("bad registry_url state, expected \"https://registry-%s.jp-east-1.gitlab.devops.nifcloud.com\", got: %#v", rName, nifcloud.ToString(instance.RegistryUrl))
		}

		return nil
	}
}

func testAccCheckDevOpsInstanceValuesUpdated(instance *types.Instance, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(instance.InstanceId) != rName {
			return fmt.Errorf("bad instance id state, expected \"%s\", got: %#v", rName, nifcloud.ToString(instance.InstanceId))
		}

		if nifcloud.ToString(instance.InstanceType) != "e-large" {
			return fmt.Errorf("bad instance_type state, expected \"e-large\", got: %#v", nifcloud.ToString(instance.InstanceType))
		}

		if nifcloud.ToString(instance.FirewallGroupName) != rName+"-upd" {
			return fmt.Errorf("bad firewall_group_name state, expected \"%s\", got: %#v", rName+"-upd", nifcloud.ToString(instance.FirewallGroupName))
		}

		if nifcloud.ToString(instance.ParameterGroupName) != rName {
			return fmt.Errorf("bad parameter_group_name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(instance.ParameterGroupName))
		}

		if nifcloud.ToInt32(instance.DiskSize) != int32(300) {
			return fmt.Errorf("bad disk_size state, expected 300, got: %#v", nifcloud.ToInt32(instance.DiskSize))
		}

		if nifcloud.ToString(instance.AvailabilityZone) != "east-14" {
			return fmt.Errorf("bad availability_zone state, expected \"east-14\", got: %#v", nifcloud.ToString(instance.AvailabilityZone))
		}

		if nifcloud.ToString(instance.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", nifcloud.ToString(instance.Description))
		}

		if nifcloud.ToString(instance.To) != "email-upd@example.com" {
			return fmt.Errorf("bad to state, expected \"email-upd@example.com\", got: %#v", nifcloud.ToString(instance.To))
		}

		if nifcloud.ToString(instance.GitlabUrl) != "https://"+rName+".jp-east-1.gitlab.devops.nifcloud.com" {
			return fmt.Errorf("bad gitlab_url state, expected \"https://%s.jp-east-1.gitlab.devops.nifcloud.com\", got: %#v", rName, nifcloud.ToString(instance.GitlabUrl))
		}

		if nifcloud.ToString(instance.RegistryUrl) != "https://registry-"+rName+".jp-east-1.gitlab.devops.nifcloud.com" {
			return fmt.Errorf("bad registry_url state, expected \"https://registry-%s.jp-east-1.gitlab.devops.nifcloud.com\", got: %#v", rName, nifcloud.ToString(instance.RegistryUrl))
		}

		return nil
	}
}

func testAccDevOpsInstanceResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).DevOps

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_devops_instance" {
			continue
		}

		_, err := svc.GetInstance(context.Background(), &devops.GetInstanceInput{
			InstanceId: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Instance" {
				return nil
			}
			return fmt.Errorf("failed GetInstance: %s", err)
		}

		return fmt.Errorf("devops instance (%s) still exists", rs.Primary.ID)
	}
	return nil
}

func testSweepDevOpsInstance(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).DevOps

	res, err := svc.ListInstances(ctx, nil)
	if err != nil {
		return err
	}

	var sweepInstances []string
	for _, i := range res.Instances {
		if strings.HasPrefix(nifcloud.ToString(i.InstanceId), prefix) {
			sweepInstances = append(sweepInstances, nifcloud.ToString(i.InstanceId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepInstances {
		instance := n
		eg.Go(func() error {
			_, err := svc.DeleteInstance(ctx, &devops.DeleteInstanceInput{
				InstanceId: nifcloud.String(instance),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
