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
	resource.AddTestSweepers("nifcloud_security_group", &resource.Sweeper{
		Name: "nifcloud_security_group",
		F:    testSweepSecurityGroup,
		Dependencies: []string{
			"nifcloud_instance",
			"nifcloud_nas_security_group",
		},
	})
}

func TestAcc_SecurityGroup(t *testing.T) {
	var securityGroup computing.SecurityGroupInfo

	resourceName := "nifcloud_security_group.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccSecurityGroupResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroup(t, "testdata/security_group.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &securityGroup),
					testAccCheckSecurityGroupValues(&securityGroup, randName),
					resource.TestCheckResourceAttr(resourceName, "group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "log_limit", "1000"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "revoke_rules_on_delete", "false"),
				),
			},
			{
				Config: testAccSecurityGroup(t, "testdata/security_group_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &securityGroup),
					testAccCheckSecurityGroupValuesUpdated(&securityGroup, randName+"upd"),

					resource.TestCheckResourceAttr(resourceName, "group_name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "log_limit", "100000"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "revoke_rules_on_delete", "true"),
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

func testAccSecurityGroup(t *testing.T, fileName, groupName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		groupName,
	)
}

func testAccCheckSecurityGroupExists(n string, securityGroup *computing.SecurityGroupInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no securityGroup resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no securityGroup id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeSecurityGroupsRequest(&computing.DescribeSecurityGroupsInput{
			GroupName: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.SecurityGroupInfo) == 0 {
			return fmt.Errorf("securityGroup does not found in cloud: %s", saved.Primary.ID)
		}

		foundSecurityGroup := res.SecurityGroupInfo[0]

		if nifcloud.StringValue(foundSecurityGroup.GroupName) != saved.Primary.ID {
			return fmt.Errorf("securityGroup does not found in cloud: %s", saved.Primary.ID)
		}

		*securityGroup = foundSecurityGroup
		return nil
	}
}

func testAccCheckSecurityGroupValues(securityGroup *computing.SecurityGroupInfo, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(securityGroup.GroupName) != groupName {
			return fmt.Errorf("bad group_name state, expected \"%s\", got: %#v", groupName, securityGroup.GroupName)
		}

		if nifcloud.StringValue(securityGroup.GroupDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", securityGroup.GroupDescription)
		}

		if nifcloud.StringValue(securityGroup.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", securityGroup.AvailabilityZone)
		}
		if nifcloud.Int64Value(securityGroup.GroupLogLimit) != 1000 {
			return fmt.Errorf("bad log_limit state,  expected \"1000\", got: %#v", securityGroup.GroupLogLimit)
		}
		return nil
	}
}

func testAccCheckSecurityGroupValuesUpdated(securityGroup *computing.SecurityGroupInfo, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(securityGroup.GroupName) != groupName {
			return fmt.Errorf("bad group_name state, expected \"%s\", got: %#v", groupName, securityGroup.GroupName)
		}

		if nifcloud.StringValue(securityGroup.GroupDescription) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", securityGroup.GroupDescription)
		}

		if nifcloud.StringValue(securityGroup.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", securityGroup.AvailabilityZone)
		}
		if nifcloud.Int64Value(securityGroup.GroupLogLimit) != 100000 {
			return fmt.Errorf("bad log_limit state,  expected \"100000\", got: %#v", securityGroup.GroupLogLimit)
		}
		return nil
	}

}

func testAccSecurityGroupResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_security_group" {
			continue
		}

		res, err := svc.DescribeSecurityGroupsRequest(&computing.DescribeSecurityGroupsInput{
			GroupName: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.SecurityGroup" {
				return fmt.Errorf("failed DescribeSecurityGroupsRequest: %s", err)
			}
		}

		if len(res.SecurityGroupInfo) > 0 {
			return fmt.Errorf("securityGroup (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepSecurityGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeSecurityGroupsRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepSecurityGroups []string
	for _, k := range res.SecurityGroupInfo {
		if strings.HasPrefix(nifcloud.StringValue(k.GroupName), prefix) {
			sweepSecurityGroups = append(sweepSecurityGroups, nifcloud.StringValue(k.GroupName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepSecurityGroups {
		groupName := n
		eg.Go(func() error {
			_, err := svc.DeleteSecurityGroupRequest(&computing.DeleteSecurityGroupInput{
				GroupName: nifcloud.String(groupName),
			}).Send(ctx)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
