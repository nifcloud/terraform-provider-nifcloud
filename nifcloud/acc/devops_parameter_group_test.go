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
	resource.AddTestSweepers("nifcloud_devops_parameter_group", &resource.Sweeper{
		Name: "nifcloud_devops_parameter_group",
		F:    testSweepDevOpsParameterGroup,
		// Dependencies: []string{
		// 	"nifcloud_devops_instance",
		// },
	})
}

func TestAcc_DevOpsParameterGroup(t *testing.T) {
	var group types.ParameterGroup

	resourceName := "nifcloud_devops_parameter_group.basic"
	randName := prefix + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDevOpsParameterGroupResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevOpsParameterGroup(t, "testdata/devops_parameter_group.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsParameterGroupExists(resourceName, &group),
					testAccCheckDevOpsParameterGroupValues(&group, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"parameter.*",
						map[string]string{
							"name":  "smtp_user_name",
							"value": "user1",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"parameter.*",
						map[string]string{
							"name":  "gitlab_email_from",
							"value": "from@mail.com",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"parameter.*",
						map[string]string{
							"name":  "gitlab_email_reply_to",
							"value": "reply-to@mail.com",
						},
					),
					resource.TestCheckResourceAttr(resourceName, "sensitive_parameter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"sensitive_parameter.*",
						map[string]string{
							"name":  "smtp_password",
							"value": "mystrongpassword",
						},
					),
				),
			},
			{
				Config: testAccDevOpsParameterGroup(t, "testdata/devops_parameter_group_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDevOpsParameterGroupExists(resourceName, &group),
					testAccCheckDevOpsParameterGroupValuesUpdated(&group, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"parameter.*",
						map[string]string{
							"name":  "smtp_user_name",
							"value": "user101",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"parameter.*",
						map[string]string{
							"name":  "gitlab_email_from",
							"value": "from@mail.com",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"parameter.*",
						map[string]string{
							"name":  "gitlab_email_reply_to",
							"value": "reply-to@mail.com",
						},
					),
					resource.TestCheckResourceAttr(resourceName, "sensitive_parameter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"sensitive_parameter.*",
						map[string]string{
							"name":  "smtp_password",
							"value": "mynewstrongpassword",
						},
					),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"parameter",
					"sensitive_parameter",
				},
			},
		},
	})
}

func testAccDevOpsParameterGroup(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), rName)
}

func testAccCheckDevOpsParameterGroupExists(n string, group *types.ParameterGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no devops parameter group resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no devops parameter group id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).DevOps
		res, err := svc.GetParameterGroup(context.Background(), &devops.GetParameterGroupInput{
			ParameterGroupName: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		if res.ParameterGroup == nil {
			return fmt.Errorf("devops parameter group is not found in cloud: %s", saved.Primary.ID)
		}

		if nifcloud.ToString(res.ParameterGroup.ParameterGroupName) != saved.Primary.ID {
			return fmt.Errorf("devops parameter group is not found in cloud: %s", saved.Primary.ID)
		}

		*group = *res.ParameterGroup

		return nil
	}
}

func testAccCheckDevOpsParameterGroupValues(group *types.ParameterGroup, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(group.ParameterGroupName) != rName {
			return fmt.Errorf("bad parameter group name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(group.ParameterGroupName))
		}

		if nifcloud.ToString(group.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", nifcloud.ToString(group.Description))
		}

		expected := map[string]string{
			"smtp_password":         "mystrongpassword",
			"smtp_user_name":        "user1",
			"gitlab_email_from":     "from@mail.com",
			"gitlab_email_reply_to": "reply-to@mail.com",
		}

		return checkDevOpsParameter(group.Parameters, expected)
	}
}

func testAccCheckDevOpsParameterGroupValuesUpdated(group *types.ParameterGroup, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(group.ParameterGroupName) != rName+"-upd" {
			return fmt.Errorf("bad parameter group name state, expected \"%s\", got: %#v", rName+"-upd", nifcloud.ToString(group.ParameterGroupName))
		}

		if nifcloud.ToString(group.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", nifcloud.ToString(group.Description))
		}

		expected := map[string]string{
			"smtp_password":         "mynewstrongpassword",
			"smtp_user_name":        "user101",
			"gitlab_email_from":     "from@mail.com",
			"gitlab_email_reply_to": "reply-to@mail.com",
		}

		return checkDevOpsParameter(group.Parameters, expected)
	}
}

func checkDevOpsParameter(params []types.Parameters, expected map[string]string) error {
	for key, val := range expected {
		found := false
		for _, p := range params {
			if nifcloud.ToString(p.Name) == key {
				if nifcloud.ToString(p.Value) == val {
					found = true
					break
				} else {
					// MEMO: This can be removed if the API returns a non masked value.
					if key == "smtp_password" && nifcloud.ToString(p.Value) == "********" {
						found = true
						break
					}
					return fmt.Errorf("bad parameter state, expected \"%s\", got: %#v", val, nifcloud.ToString(p.Value))
				}
			}
		}
		if !found {
			return fmt.Errorf("bad parameter state, %s is not found", key)
		}
	}

	return nil
}

func testAccDevOpsParameterGroupResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).DevOps

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_devops_parameter_group" {
			continue
		}

		_, err := svc.GetParameterGroup(context.Background(), &devops.GetParameterGroupInput{
			ParameterGroupName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.ParameterGroup" {
				return nil
			}
			return fmt.Errorf("failed GetParameterGroup: %s", err)
		}

		return fmt.Errorf("devops parameter group (%s) still exists", rs.Primary.ID)
	}
	return nil
}

func testSweepDevOpsParameterGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).DevOps

	res, err := svc.ListParameterGroups(ctx, nil)
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
			_, err := svc.DeleteParameterGroup(ctx, &devops.DeleteParameterGroupInput{
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
