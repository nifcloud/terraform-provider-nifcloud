package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
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
	var dbSecurityGroup rdb.DBSecurityGroup

	resourceName := "nifcloud_db_security_group.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

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
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
	)
}

func testAccCheckDbSecurityGroupExists(n string, dbSecurityGroup *rdb.DBSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no dbSecurityGroup resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no dbSecurityGroup id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).RDB
		res, err := svc.DescribeDBSecurityGroupsRequest(&rdb.DescribeDBSecurityGroupsInput{
			DBSecurityGroupName: nifcloud.String(saved.Primary.ID),
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.DBSecurityGroups) == 0 {
			return fmt.Errorf("dbSecurityGroup does not found in cloud: %s", saved.Primary.ID)
		}

		foundDbSecurityGroup := res.DBSecurityGroups[0]

		if nifcloud.StringValue(foundDbSecurityGroup.DBSecurityGroupName) != saved.Primary.ID {
			return fmt.Errorf("dbSecurityGroup does not found in cloud: %s", saved.Primary.ID)
		}

		*dbSecurityGroup = foundDbSecurityGroup
		return nil
	}
}

func testAccCheckDbSecurityGroupValues(dbSecurityGroup *rdb.DBSecurityGroup, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(dbSecurityGroup.IPRanges) != 1 {
			return fmt.Errorf("bad cidr_ip rules: %#v", dbSecurityGroup.IPRanges)
		}

		if len(dbSecurityGroup.EC2SecurityGroups) != 1 {
			return fmt.Errorf("bad security_group_name rules: %#v", dbSecurityGroup.EC2SecurityGroups)
		}

		if nifcloud.StringValue(dbSecurityGroup.DBSecurityGroupName) != groupName {
			return fmt.Errorf("bad group_name state, expected \"%s\", got: %#v", groupName, dbSecurityGroup.DBSecurityGroupName)
		}

		if nifcloud.StringValue(dbSecurityGroup.DBSecurityGroupDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", dbSecurityGroup.DBSecurityGroupDescription)
		}

		if nifcloud.StringValue(dbSecurityGroup.NiftyAvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", dbSecurityGroup.NiftyAvailabilityZone)
		}

		cidrIPRule := make(map[int]rdb.IPRange)
		securityGroupNameRule := make(map[int]rdb.EC2SecurityGroup)

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

		if nifcloud.StringValue(cidrIPRule[0].CIDRIP) != "0.0.0.0/0" {
			return fmt.Errorf("bad cide_ip rule, expected \"0.0.0.0/0\", got: %#v", nifcloud.StringValue(cidrIPRule[0].CIDRIP))
		}

		if nifcloud.StringValue(securityGroupNameRule[0].EC2SecurityGroupName) != groupName {
			return fmt.Errorf("bad security_group_name rule, expected \"%s\", got: %#v", groupName, nifcloud.StringValue(securityGroupNameRule[0].EC2SecurityGroupName))
		}

		return nil
	}
}

func testAccCheckDbSecurityGroupValuesUpdated(dbSecurityGroup *rdb.DBSecurityGroup, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(dbSecurityGroup.IPRanges) != 1 {
			return fmt.Errorf("bad cidr_ip rules: %#v", dbSecurityGroup.IPRanges)
		}

		if len(dbSecurityGroup.EC2SecurityGroups) != 1 {
			return fmt.Errorf("bad security_group_name rules: %#v", dbSecurityGroup.EC2SecurityGroups)
		}

		if nifcloud.StringValue(dbSecurityGroup.DBSecurityGroupName) != groupName {
			return fmt.Errorf("bad group_name state, expected \"%s\", got: %#v", groupName, dbSecurityGroup.DBSecurityGroupName)
		}

		if nifcloud.StringValue(dbSecurityGroup.DBSecurityGroupDescription) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", dbSecurityGroup.DBSecurityGroupDescription)
		}

		if nifcloud.StringValue(dbSecurityGroup.NiftyAvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", dbSecurityGroup.NiftyAvailabilityZone)
		}

		cidrIPRule := make(map[int]rdb.IPRange)
		securityGroupNameRule := make(map[int]rdb.EC2SecurityGroup)

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

		if nifcloud.StringValue(cidrIPRule[0].CIDRIP) != "192.168.0.1/32" {
			return fmt.Errorf("bad cide_ip rule, expected \"192.168.0.1/32\", got: %#v", nifcloud.StringValue(cidrIPRule[0].CIDRIP))
		}

		if nifcloud.StringValue(securityGroupNameRule[0].EC2SecurityGroupName) != groupName {
			return fmt.Errorf("bad security_group_name rule, expected \"%s\", got: %#v", groupName, nifcloud.StringValue(securityGroupNameRule[0].EC2SecurityGroupName))
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

		res, err := svc.DescribeDBSecurityGroupsRequest(&rdb.DescribeDBSecurityGroupsInput{
			DBSecurityGroupName: nifcloud.String(rs.Primary.ID),
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.DBSecurityGroup" {
				return fmt.Errorf("failed DescribeDBSecurityGroupsRequest: %s", err)
			}
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

	res, err := svc.DescribeDBSecurityGroupsRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	for _, dbSecurityGroup := range res.DBSecurityGroups {

		input := &rdb.DeleteDBSecurityGroupInput{
			DBSecurityGroupName: dbSecurityGroup.DBSecurityGroupName,
		}

		_, err := svc.DeleteDBSecurityGroupRequest(input).Send(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
