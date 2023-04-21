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
	resource.AddTestSweepers("nifcloud_db_security_group", &resource.Sweeper{
		Name: "nifcloud_db_security_group",
		F:    testSweepDbSecurityGroup,
		Dependencies: []string{
			"nifcloud_db_instance",
		},
	})
}

func TestAcc_DbSecurityGroup(t *testing.T) {
	var dbSecurityGroup types.DBSecurityGroupsOfDescribeDBSecurityGroups

	resourceName := "nifcloud_db_security_group.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDbSecurityGroupResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbSecurityGroup(t, "testdata/db_security_group.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbSecurityGroupExists(resourceName, &dbSecurityGroup),
					testAccCheckDbSecurityGroupValues(&dbSecurityGroup, randName),
					resource.TestCheckResourceAttr(resourceName, "group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.security_group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "rule.1.cidr_ip", "0.0.0.0/0"),
				),
			},
			{
				Config: testAccDbSecurityGroup(t, "testdata/db_security_group_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbSecurityGroupExists(resourceName, &dbSecurityGroup),
					testAccCheckDbSecurityGroupValuesUpdated(&dbSecurityGroup, randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "group_name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.security_group_name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.cidr_ip", "192.168.0.1/32"),
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

func testAccDbSecurityGroup(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
	)
}

func testAccCheckDbSecurityGroupExists(n string, dbSecurityGroup *types.DBSecurityGroupsOfDescribeDBSecurityGroups) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no dbSecurityGroup resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no dbSecurityGroup id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).RDB
		res, err := svc.DescribeDBSecurityGroups(context.Background(), &rdb.DescribeDBSecurityGroupsInput{
			DBSecurityGroupName: nifcloud.String(saved.Primary.ID),
		})

		if err != nil {
			return err
		}

		if len(res.DBSecurityGroups) == 0 {
			return fmt.Errorf("dbSecurityGroup does not found in cloud: %s", saved.Primary.ID)
		}

		foundDbSecurityGroup := res.DBSecurityGroups[0]

		if nifcloud.ToString(foundDbSecurityGroup.DBSecurityGroupName) != saved.Primary.ID {
			return fmt.Errorf("dbSecurityGroup does not found in cloud: %s", saved.Primary.ID)
		}

		*dbSecurityGroup = foundDbSecurityGroup
		return nil
	}
}

func testAccCheckDbSecurityGroupValues(dbSecurityGroup *types.DBSecurityGroupsOfDescribeDBSecurityGroups, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(dbSecurityGroup.IPRanges) != 1 {
			return fmt.Errorf("bad cidr_ip rules: %#v", dbSecurityGroup.IPRanges)
		}

		if len(dbSecurityGroup.EC2SecurityGroups) != 1 {
			return fmt.Errorf("bad security_group_name rules: %#v", dbSecurityGroup.EC2SecurityGroups)
		}

		if nifcloud.ToString(dbSecurityGroup.DBSecurityGroupName) != groupName {
			return fmt.Errorf("bad group_name state, expected \"%s\", got: %#v", groupName, dbSecurityGroup.DBSecurityGroupName)
		}

		if nifcloud.ToString(dbSecurityGroup.DBSecurityGroupDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", dbSecurityGroup.DBSecurityGroupDescription)
		}

		if nifcloud.ToString(dbSecurityGroup.NiftyAvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", dbSecurityGroup.NiftyAvailabilityZone)
		}

		cidrIPRule := make(map[int]types.IPRanges)
		securityGroupNameRule := make(map[int]types.EC2SecurityGroups)

		for _, c := range dbSecurityGroup.IPRanges {
			cidrIPRule[0] = c
		}
		for _, s := range dbSecurityGroup.EC2SecurityGroups {
			securityGroupNameRule[0] = s
		}

		if _, ok := cidrIPRule[0]; !ok {
			return fmt.Errorf("bad cidr_ip rule: %#v", dbSecurityGroup.IPRanges)
		}

		if _, ok := securityGroupNameRule[0]; !ok {
			return fmt.Errorf("bad security_group_name rule: %#v", dbSecurityGroup.EC2SecurityGroups)
		}

		if nifcloud.ToString(cidrIPRule[0].CIDRIP) != "0.0.0.0/0" {
			return fmt.Errorf("bad cide_ip rule, expected \"0.0.0.0/0\", got: %#v", nifcloud.ToString(cidrIPRule[0].CIDRIP))
		}

		if nifcloud.ToString(securityGroupNameRule[0].EC2SecurityGroupName) != groupName {
			return fmt.Errorf("bad security_group_name rule, expected \"%s\", got: %#v", groupName, nifcloud.ToString(securityGroupNameRule[0].EC2SecurityGroupName))
		}

		return nil
	}
}

func testAccCheckDbSecurityGroupValuesUpdated(dbSecurityGroup *types.DBSecurityGroupsOfDescribeDBSecurityGroups, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(dbSecurityGroup.IPRanges) != 1 {
			return fmt.Errorf("bad cidr_ip rules: %#v", dbSecurityGroup.IPRanges)
		}

		if len(dbSecurityGroup.EC2SecurityGroups) != 1 {
			return fmt.Errorf("bad security_group_name rules: %#v", dbSecurityGroup.EC2SecurityGroups)
		}

		if nifcloud.ToString(dbSecurityGroup.DBSecurityGroupName) != groupName {
			return fmt.Errorf("bad group_name state, expected \"%s\", got: %#v", groupName, dbSecurityGroup.DBSecurityGroupName)
		}

		if nifcloud.ToString(dbSecurityGroup.DBSecurityGroupDescription) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", dbSecurityGroup.DBSecurityGroupDescription)
		}

		if nifcloud.ToString(dbSecurityGroup.NiftyAvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", dbSecurityGroup.NiftyAvailabilityZone)
		}

		cidrIPRule := make(map[int]types.IPRanges)
		securityGroupNameRule := make(map[int]types.EC2SecurityGroups)

		for _, c := range dbSecurityGroup.IPRanges {
			cidrIPRule[0] = c
		}
		for _, s := range dbSecurityGroup.EC2SecurityGroups {
			securityGroupNameRule[0] = s
		}

		if _, ok := cidrIPRule[0]; !ok {
			return fmt.Errorf("bad cidr_ip rule: %#v", dbSecurityGroup.IPRanges)
		}

		if _, ok := securityGroupNameRule[0]; !ok {
			return fmt.Errorf("bad security_group_name rule: %#v", dbSecurityGroup.EC2SecurityGroups)
		}

		if nifcloud.ToString(cidrIPRule[0].CIDRIP) != "192.168.0.1/32" {
			return fmt.Errorf("bad cide_ip rule, expected \"192.168.0.1/32\", got: %#v", nifcloud.ToString(cidrIPRule[0].CIDRIP))
		}

		if nifcloud.ToString(securityGroupNameRule[0].EC2SecurityGroupName) != groupName {
			return fmt.Errorf("bad security_group_name rule, expected \"%s\", got: %#v", groupName, nifcloud.ToString(securityGroupNameRule[0].EC2SecurityGroupName))
		}

		return nil
	}
}

func testAccDbSecurityGroupResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).RDB

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_db_security_group" {
			continue
		}

		res, err := svc.DescribeDBSecurityGroups(context.Background(), &rdb.DescribeDBSecurityGroupsInput{
			DBSecurityGroupName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.DBSecurityGroup" {
				return nil
			}
			return fmt.Errorf("failed DescribeDBSecurityGroupsRequest: %s", err)
		}

		if len(res.DBSecurityGroups) > 0 {
			return fmt.Errorf("dbSecurityGroup (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepDbSecurityGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).RDB

	res, err := svc.DescribeDBSecurityGroups(ctx, nil)
	if err != nil {
		return err
	}

	var sweepDbSecurityGroups []string
	for _, k := range res.DBSecurityGroups {
		if strings.HasPrefix(nifcloud.ToString(k.DBSecurityGroupName), prefix) {
			sweepDbSecurityGroups = append(sweepDbSecurityGroups, nifcloud.ToString(k.DBSecurityGroupName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepDbSecurityGroups {
		groupName := n
		eg.Go(func() error {
			_, err := svc.DeleteDBSecurityGroup(ctx, &rdb.DeleteDBSecurityGroupInput{
				DBSecurityGroupName: nifcloud.String(groupName),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
