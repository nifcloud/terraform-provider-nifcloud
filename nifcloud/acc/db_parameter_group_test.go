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
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_db_parameter_group", &resource.Sweeper{
		Name: "nifcloud_db_parameter_group",
		F:    testSweepDBParameterGroup,
		Dependencies: []string{
			"nifcloud_db_instance",
		},
	})
}

func TestAcc_DBParameterGroup(t *testing.T) {
	var group types.DBParameterGroupsOfDescribeDBParameterGroups

	resourceName := "nifcloud_db_parameter_group.basic"
	randName := prefix + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDBParameterGroupResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBParameterGroup(t, "testdata/db_parameter_group.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBParameterGroupExists(resourceName, &group),
					testAccCheckDBParameterGroupValues(&group, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "family", "mysql5.6"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "3"),
				),
			},
			{
				Config: testAccDBParameterGroup(t, "testdata/db_parameter_group_update_only_parameter.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBParameterGroupExists(resourceName, &group),
					testAccCheckDBParameterGroupValuesUpdatedOnlyParameters(&group, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "family", "mysql5.6"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "3"),
				),
			},
			{
				Config: testAccDBParameterGroup(t, "testdata/db_parameter_group_update_all.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBParameterGroupExists(resourceName, &group),
					testAccCheckDBParameterGroupValuesUpdatedAll(&group, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "family", "mysql5.7"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "2"),
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

func testAccDBParameterGroup(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), rName)
}

func testAccCheckDBParameterGroupExists(n string, group *types.DBParameterGroupsOfDescribeDBParameterGroups) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no db parameter group resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no db parameter group id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).RDB
		res, err := svc.DescribeDBParameterGroups(context.Background(), &rdb.DescribeDBParameterGroupsInput{
			DBParameterGroupName: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		if len(res.DBParameterGroups) == 0 {
			return fmt.Errorf("db parameter group does not found in cloud: %s", saved.Primary.ID)
		}

		foundGroup := res.DBParameterGroups[0]

		if nifcloud.ToString(foundGroup.DBParameterGroupName) != saved.Primary.ID {
			return fmt.Errorf("db parameter group does not found in cloud: %s", saved.Primary.ID)
		}

		*group = foundGroup

		return nil
	}
}

func testAccCheckDBParameterGroupValues(group *types.DBParameterGroupsOfDescribeDBParameterGroups, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(group.DBParameterGroupName) != rName {
			return fmt.Errorf("bad db parameter group name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(group.DBParameterGroupName))
		}

		if nifcloud.ToString(group.DBParameterGroupFamily) != "mysql5.6" {
			return fmt.Errorf("bad db parameter group family state, expected \"mysql5.6\", got: %#v", nifcloud.ToString(group.DBParameterGroupFamily))
		}

		if nifcloud.ToString(group.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", nifcloud.ToString(group.Description))
		}

		svc := testAccProvider.Meta().(*client.Client).RDB
		res, err := svc.DescribeDBParameters(context.Background(), &rdb.DescribeDBParametersInput{
			DBParameterGroupName: group.DBParameterGroupName,
			Source:               types.SourceOfDescribeDBParametersRequestUser,
		})
		if err != nil {
			return fmt.Errorf("failed describe DBParameterGroup: %s", err)
		}

		if len(res.Parameters) != 3 {
			return fmt.Errorf("bad customized parameter size, expected len 3, got: %d", len(res.Parameters))
		}

		expected := map[string]string{
			"character_set_server":  "utf8",
			"character_set_client":  "utf8",
			"character_set_results": "utf8",
		}

		return checkParameter(res.Parameters, expected)
	}
}

func testAccCheckDBParameterGroupValuesUpdatedOnlyParameters(group *types.DBParameterGroupsOfDescribeDBParameterGroups, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(group.DBParameterGroupName) != rName {
			return fmt.Errorf("bad db parameter group name state, expected \"%s\", got: %#v", rName, nifcloud.ToString(group.DBParameterGroupName))
		}

		if nifcloud.ToString(group.DBParameterGroupFamily) != "mysql5.6" {
			return fmt.Errorf("bad db parameter group family state, expected \"mysql5.6\", got: %#v", nifcloud.ToString(group.DBParameterGroupFamily))
		}

		if nifcloud.ToString(group.Description) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", nifcloud.ToString(group.Description))
		}

		svc := testAccProvider.Meta().(*client.Client).RDB
		res, err := svc.DescribeDBParameters(context.Background(), &rdb.DescribeDBParametersInput{
			DBParameterGroupName: group.DBParameterGroupName,
			Source:               types.SourceOfDescribeDBParametersRequestUser,
		})
		if err != nil {
			return fmt.Errorf("failed describe DBParameterGroup: %s", err)
		}

		if len(res.Parameters) != 3 {
			return fmt.Errorf("bad customized parameter size, expected len 3, got: %d", len(res.Parameters))
		}

		expected := map[string]string{
			"character_set_server":  "ascii",
			"character_set_client":  "utf8",
			"character_set_results": "utf8",
		}

		return checkParameter(res.Parameters, expected)
	}
}

func testAccCheckDBParameterGroupValuesUpdatedAll(group *types.DBParameterGroupsOfDescribeDBParameterGroups, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(group.DBParameterGroupName) != rName+"upd" {
			return fmt.Errorf("bad db parameter group name state, expected \"%s\", got: %#v", rName+"upd", group.DBParameterGroupName)
		}

		if nifcloud.ToString(group.DBParameterGroupFamily) != "mysql5.7" {
			return fmt.Errorf("bad db parameter group family state,  expected \"mysql5.7\", got: %#v", group.DBParameterGroupFamily)
		}

		if nifcloud.ToString(group.Description) != "tfacc-memo-upd" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo-upd\", got: %#v", group.Description)
		}

		svc := testAccProvider.Meta().(*client.Client).RDB
		res, err := svc.DescribeDBParameters(context.Background(), &rdb.DescribeDBParametersInput{
			DBParameterGroupName: group.DBParameterGroupName,
			Source:               types.SourceOfDescribeDBParametersRequestUser,
		})
		if err != nil {
			return fmt.Errorf("failed describe DBParameterGroup: %s", err)
		}

		if len(res.Parameters) != 2 {
			return fmt.Errorf("bad customized parameter size, expected len 2, got: %d", len(res.Parameters))
		}

		expected := map[string]string{
			"character_set_server":  "ascii",
			"character_set_results": "ascii",
		}

		return checkParameter(res.Parameters, expected)
	}
}

func checkParameter(params []types.Parameters, expected map[string]string) error {
	for _, p := range params {
		if val, ok := expected[*p.ParameterName]; ok {
			if nifcloud.ToString(p.ParameterValue) != val {
				return fmt.Errorf("bad parameter state, expected \"%s\", got: %#v", val, nifcloud.ToString(p.ParameterValue))
			}
		} else {
			return fmt.Errorf("bad parameter state, %s is unexpected parameter", nifcloud.ToString(p.ParameterName))
		}
	}

	return nil
}

func testAccDBParameterGroupResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).RDB

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_db_parameter_group" {
			continue
		}

		res, err := svc.DescribeDBParameterGroups(context.Background(), &rdb.DescribeDBParameterGroupsInput{
			DBParameterGroupName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.DBParameterGroup" {
				return nil
			}
			return fmt.Errorf("failed DescribeDBParameterGroups: %s", err)
		}

		if len(res.DBParameterGroups) > 0 {
			return fmt.Errorf("db parameter group (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepDBParameterGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).RDB

	res, err := svc.DescribeDBParameterGroups(ctx, nil)
	if err != nil {
		return err
	}

	var sweepDBParameterGroups []string
	for _, g := range res.DBParameterGroups {
		if strings.HasPrefix(nifcloud.ToString(g.DBParameterGroupName), prefix) {
			sweepDBParameterGroups = append(sweepDBParameterGroups, nifcloud.ToString(g.DBParameterGroupName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepDBParameterGroups {
		group := n
		eg.Go(func() error {
			_, err := svc.DeleteDBParameterGroup(ctx, &rdb.DeleteDBParameterGroupInput{
				DBParameterGroupName: nifcloud.String(group),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
