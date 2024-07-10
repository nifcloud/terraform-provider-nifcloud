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
	resource.AddTestSweepers("nifcloud_devops_runner_parameter_group", &resource.Sweeper{
		Name: "nifcloud_devops_runner_parameter_group",
		F:    testSweepDevOpsRunnerParameterGroup,
		// Dependencies: []string{
		// 	"nifcloud_devops_runner",
		// },
	})
}

func TestAcc_DevOpsRunnerParameterGroup(t *testing.T) {
	var group types.ParameterGroup

	resourceName := "nifcloud_devops_runner_parameter_group.basic"
	randName := prefix + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDevOpsRunnerParameterGroupResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevOpsRunnerParameterGroup(t, "testdata/devops_runner_parameter_group.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsRunnerParameterGroupExists(resourceName, &group),
					testAccCheckDevOpsRunnerParameterGroupValues(&group, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "docker_image", "ruby"),
					resource.TestCheckResourceAttr(resourceName, "docker_privileged", "true"),
					resource.TestCheckResourceAttr(resourceName, "docker_shm_size", "300000"),
					resource.TestCheckResourceAttr(resourceName, "docker_extra_host.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"docker_extra_host.*",
						map[string]string{
							"host_name":  "example.test",
							"ip_address": "192.168.1.2",
						},
					),
					resource.TestCheckResourceAttr(resourceName, "docker_volume.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "docker_volume.*", "/user_data:/cache"),
				),
			},
			{
				Config: testAccDevOpsRunnerParameterGroup(t, "testdata/devops_runner_parameter_group_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsRunnerParameterGroupExists(resourceName, &group),
					testAccCheckDevOpsRunnerParameterGroupValuesUpdated(&group, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "docker_image", "ruby:3"),
					resource.TestCheckResourceAttr(resourceName, "docker_privileged", "false"),
					resource.TestCheckResourceAttr(resourceName, "docker_shm_size", "600000"),
					resource.TestCheckResourceAttr(resourceName, "docker_extra_host.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"docker_extra_host.*",
						map[string]string{
							"host_name":  "example.test",
							"ip_address": "192.168.1.3",
						},
					),
					resource.TestCheckResourceAttr(resourceName, "docker_volume.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "docker_volume.*", "/user_data"),
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

func testAccDevOpsRunnerParameterGroup(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), rName)
}

func testAccCheckDevOpsRunnerParameterGroupExists(n string, group *types.ParameterGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no devops runner parameter group resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no devops runner parameter group id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DevOpsRunner
		res, err := svc.GetRunnerParameterGroup(context.Background(), &devopsrunner.GetRunnerParameterGroupInput{
			ParameterGroupName: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		if res.ParameterGroup == nil {
			return fmt.Errorf("devops runner parameter group is not found in cloud: %s", saved.Primary.ID)
		}

		if nifcloud.ToString(res.ParameterGroup.ParameterGroupName) != saved.Primary.ID {
			return fmt.Errorf("devops runner parameter group is not found in cloud: %s", saved.Primary.ID)
		}

		*group = *res.ParameterGroup

		return nil
	}
}

func testAccCheckDevOpsRunnerParameterGroupValues(group *types.ParameterGroup, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(group.ParameterGroupName) != rName {
			return fmt.Errorf("bad runner parameter group name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(group.ParameterGroupName))
		}

		if nifcloud.ToString(group.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", nifcloud.ToString(group.Description))
		}

		if nifcloud.ToBool(group.DockerParameters.DisableCache) != false {
			return fmt.Errorf("bad docker_disable_cache state, expected false, got: %#v", nifcloud.ToBool(group.DockerParameters.DisableCache))
		}

		if nifcloud.ToBool(group.DockerParameters.DisableEntrypointOverwrite) != false {
			return fmt.Errorf("bad docker_disable_entrypoint_overwrite state, expected false, got: %#v", nifcloud.ToBool(group.DockerParameters.DisableEntrypointOverwrite))
		}

		if len(group.DockerParameters.ExtraHosts) != 1 {
			return fmt.Errorf("bad docker_extra_host length: %#v", group.DockerParameters.ExtraHosts)
		}

		if nifcloud.ToString(group.DockerParameters.ExtraHosts[0].HostName) != "example.test" {
			return fmt.Errorf("bad host_name state, expected \"example.test\", got: %#v", group.DockerParameters.ExtraHosts[0].HostName)
		}

		if nifcloud.ToString(group.DockerParameters.ExtraHosts[0].IpAddress) != "192.168.1.2" {
			return fmt.Errorf("bad ip_address state, expected \"192.168.1.2\", got: %#v", group.DockerParameters.ExtraHosts[0].IpAddress)
		}

		if nifcloud.ToString(group.DockerParameters.Image) != "ruby" {
			return fmt.Errorf("bad docker_image state, expected \"ruby\", got: %#v", nifcloud.ToString(group.DockerParameters.Image))
		}

		if nifcloud.ToBool(group.DockerParameters.OomKillDisable) != false {
			return fmt.Errorf("bad docker_oom_kill_disable state, expected false, got: %#v", nifcloud.ToBool(group.DockerParameters.OomKillDisable))
		}

		if nifcloud.ToBool(group.DockerParameters.Privileged) != true {
			return fmt.Errorf("bad docker_privileged state, expected true, got: %#v", nifcloud.ToBool(group.DockerParameters.Privileged))
		}

		if nifcloud.ToInt32(group.DockerParameters.ShmSize) != int32(300000) {
			return fmt.Errorf("bad docker_shm_size state, expected 300000, got: %#v", nifcloud.ToInt32(group.DockerParameters.ShmSize))
		}

		if nifcloud.ToBool(group.DockerParameters.TlsVerify) != false {
			return fmt.Errorf("bad docker_tls_verify state, expected false, got: %#v", nifcloud.ToBool(group.DockerParameters.TlsVerify))
		}

		if len(group.DockerParameters.Volumes) != 1 {
			return fmt.Errorf("bad docker_volume length: %#v", group.DockerParameters.Volumes)
		}

		if group.DockerParameters.Volumes[0] != "/user_data:/cache" {
			return fmt.Errorf("bad docker_volume state, expected \"/user_data:/cache\", got: %#v", group.DockerParameters.Volumes[0])
		}

		return nil
	}
}

func testAccCheckDevOpsRunnerParameterGroupValuesUpdated(group *types.ParameterGroup, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(group.ParameterGroupName) != rName+"-upd" {
			return fmt.Errorf("bad parameter group name state, expected \"%s\", got: %#v", rName+"-upd", nifcloud.ToString(group.ParameterGroupName))
		}

		if nifcloud.ToString(group.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", nifcloud.ToString(group.Description))
		}

		if nifcloud.ToBool(group.DockerParameters.DisableCache) != false {
			return fmt.Errorf("bad docker_disable_cache state, expected false, got: %#v", nifcloud.ToBool(group.DockerParameters.DisableCache))
		}

		if nifcloud.ToBool(group.DockerParameters.DisableEntrypointOverwrite) != false {
			return fmt.Errorf("bad docker_disable_entrypoint_overwrite state, expected false, got: %#v", nifcloud.ToBool(group.DockerParameters.DisableEntrypointOverwrite))
		}

		if len(group.DockerParameters.ExtraHosts) != 1 {
			return fmt.Errorf("bad docker_extra_host length: %#v", group.DockerParameters.ExtraHosts)
		}

		if nifcloud.ToString(group.DockerParameters.ExtraHosts[0].HostName) != "example.test" {
			return fmt.Errorf("bad host_name state, expected \"example.test\", got: %#v", group.DockerParameters.ExtraHosts[0].HostName)
		}

		if nifcloud.ToString(group.DockerParameters.ExtraHosts[0].IpAddress) != "192.168.1.3" {
			return fmt.Errorf("bad ip_address state, expected \"192.168.1.3\", got: %#v", group.DockerParameters.ExtraHosts[0].IpAddress)
		}

		if nifcloud.ToString(group.DockerParameters.Image) != "ruby:3" {
			return fmt.Errorf("bad docker_image state, expected \"ruby:3\", got: %#v", nifcloud.ToString(group.DockerParameters.Image))
		}

		if nifcloud.ToBool(group.DockerParameters.OomKillDisable) != false {
			return fmt.Errorf("bad docker_oom_kill_disable state, expected false, got: %#v", nifcloud.ToBool(group.DockerParameters.OomKillDisable))
		}

		if nifcloud.ToBool(group.DockerParameters.Privileged) != false {
			return fmt.Errorf("bad docker_privileged state, expected false, got: %#v", nifcloud.ToBool(group.DockerParameters.Privileged))
		}

		if nifcloud.ToInt32(group.DockerParameters.ShmSize) != int32(600000) {
			return fmt.Errorf("bad docker_shm_size state, expected 600000, got: %#v", nifcloud.ToInt32(group.DockerParameters.ShmSize))
		}

		if nifcloud.ToBool(group.DockerParameters.TlsVerify) != false {
			return fmt.Errorf("bad docker_tls_verify state, expected false, got: %#v", nifcloud.ToBool(group.DockerParameters.TlsVerify))
		}

		if len(group.DockerParameters.Volumes) != 1 {
			return fmt.Errorf("bad docker_volume length: %#v", group.DockerParameters.Volumes)
		}

		if group.DockerParameters.Volumes[0] != "/user_data" {
			return fmt.Errorf("bad docker_volume state, expected \"/user_data\", got: %#v", group.DockerParameters.Volumes[0])
		}

		return nil
	}
}

func testAccDevOpsRunnerParameterGroupResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).DevOpsRunner

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_devops_runner_parameter_group" {
			continue
		}

		_, err := svc.GetRunnerParameterGroup(context.Background(), &devopsrunner.GetRunnerParameterGroupInput{
			ParameterGroupName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.ParameterGroup" {
				return nil
			}
			return fmt.Errorf("failed GetRunnerParameterGroup: %s", err)
		}

		return fmt.Errorf("devops runner parameter group (%s) still exists", rs.Primary.ID)
	}
	return nil
}

func testSweepDevOpsRunnerParameterGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).DevOpsRunner

	res, err := svc.ListRunnerParameterGroups(ctx, nil)
	if err != nil {
		return err
	}

	var sweepParameterGroups []string
	for _, g := range res.ParameterGroups {
		if strings.HasPrefix(nifcloud.ToString(g.ParameterGroupName), prefix) {
			sweepParameterGroups = append(sweepParameterGroups, nifcloud.ToString(g.ParameterGroupName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepParameterGroups {
		group := n
		eg.Go(func() error {
			_, err := svc.DeleteRunnerParameterGroup(ctx, &devopsrunner.DeleteRunnerParameterGroupInput{
				ParameterGroupName: nifcloud.String(group),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
