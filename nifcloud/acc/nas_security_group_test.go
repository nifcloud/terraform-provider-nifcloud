package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_nas_security_group", &resource.Sweeper{
		Name: "nifcloud_nas_security_group",
		F:    testSweepNASSecurityGroup,
		Dependencies: []string{
			"nifcloud_nas_instance",
		},
	})
}

func TestAcc_NASSecurityGroup(t *testing.T) {
	var nasSecurityGroup types.NASSecurityGroupsOfDescribeNASSecurityGroups

	resourceName := "nifcloud_nas_security_group.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccNASSecurityGroupResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNASSecurityGroup(t, "testdata/nas_security_group.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNASSecurityGroupExists(resourceName, &nasSecurityGroup),
					testAccCheckNASSecurityGroupValues(&nasSecurityGroup, randName),
					resource.TestCheckResourceAttr(resourceName, "group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "3"),
				),
			},
			{
				Config: testAccNASSecurityGroup(t, "testdata/nas_security_group_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNASSecurityGroupExists(resourceName, &nasSecurityGroup),
					testAccCheckNASSecurityGroupValuesUpdated(&nasSecurityGroup, randName),
					resource.TestCheckResourceAttr(resourceName, "group_name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "3"),
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

func testAccNASSecurityGroup(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
	)
}

func testAccCheckNASSecurityGroupExists(n string, nasSecurityGroup *types.NASSecurityGroupsOfDescribeNASSecurityGroups) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no nasSecurityGroup resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no nasSecurityGroup id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).NAS
		res, err := svc.DescribeNASSecurityGroups(context.Background(), &nas.DescribeNASSecurityGroupsInput{
			NASSecurityGroupName: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		if len(res.NASSecurityGroups) == 0 {
			return fmt.Errorf("nasSecurityGroup does not found in cloud: %s", saved.Primary.ID)
		}

		foundNASSecurityGroup := res.NASSecurityGroups[0]

		if nifcloud.ToString(foundNASSecurityGroup.NASSecurityGroupName) != saved.Primary.ID {
			return fmt.Errorf("nasSecurityGroup does not found in cloud: %s", saved.Primary.ID)
		}

		*nasSecurityGroup = foundNASSecurityGroup

		return nil
	}
}

func testAccCheckNASSecurityGroupValues(nasSecurityGroup *types.NASSecurityGroupsOfDescribeNASSecurityGroups, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(nasSecurityGroup.NASSecurityGroupName) != groupName {
			return fmt.Errorf("bad group_name state, expected \"%s\", got: %#v", groupName, nifcloud.ToString(nasSecurityGroup.NASSecurityGroupName))
		}

		if nifcloud.ToString(nasSecurityGroup.NASSecurityGroupDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", nifcloud.ToString(nasSecurityGroup.NASSecurityGroupDescription))
		}

		if nifcloud.ToString(nasSecurityGroup.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", nifcloud.ToString(nasSecurityGroup.AvailabilityZone))
		}

		wantCidrIPs := []string{"192.168.0.1/32"}
		if len(nasSecurityGroup.IPRanges) != len(wantCidrIPs) {
			return fmt.Errorf("bad rule[*] state, expected length that having cidr_ip is %d, got length: %d", len(wantCidrIPs), len(nasSecurityGroup.IPRanges))
		}

		gotCidrIps := []string{}
		for _, ipRange := range nasSecurityGroup.IPRanges {
			gotCidrIps = append(gotCidrIps, nifcloud.ToString(ipRange.CIDRIP))
		}

		sort.Strings(wantCidrIPs)
		sort.Strings(gotCidrIps)

		for i, want := range wantCidrIPs {
			if want != gotCidrIps[i] {
				return fmt.Errorf("bad rule[*].cidr_ip state, expected %q, got: %#v", want, gotCidrIps[i])
			}
		}

		wantGroupNames := []string{groupName + "01", groupName + "02"}
		if len(nasSecurityGroup.SecurityGroups) != len(wantGroupNames) {
			return fmt.Errorf("bad rule[*] state, expected length that having security_group_name is %d, got length: %d", len(wantGroupNames), len(nasSecurityGroup.SecurityGroups))
		}

		gotGroupNames := []string{}
		for _, securityGroup := range nasSecurityGroup.SecurityGroups {
			gotGroupNames = append(gotGroupNames, nifcloud.ToString(securityGroup.SecurityGroupName))
		}

		sort.Strings(wantGroupNames)
		sort.Strings(gotGroupNames)

		for i, want := range wantGroupNames {
			if want != gotGroupNames[i] {
				return fmt.Errorf("bad rule[*].security_group_name state, expected %q, got: %#v", want, gotGroupNames[i])
			}
		}

		return nil
	}
}

func testAccCheckNASSecurityGroupValuesUpdated(nasSecurityGroup *types.NASSecurityGroupsOfDescribeNASSecurityGroups, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(nasSecurityGroup.NASSecurityGroupName) != groupName+"upd" {
			return fmt.Errorf("bad group_name state, expected \"%supd\", got: %#v", groupName, nifcloud.ToString(nasSecurityGroup.NASSecurityGroupName))
		}

		if nifcloud.ToString(nasSecurityGroup.NASSecurityGroupDescription) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", nifcloud.ToString(nasSecurityGroup.NASSecurityGroupDescription))
		}

		if nifcloud.ToString(nasSecurityGroup.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", nifcloud.ToString(nasSecurityGroup.AvailabilityZone))
		}

		wantCidrIPs := []string{"192.168.0.2/32", "192.168.0.3/32"}
		if len(nasSecurityGroup.IPRanges) != len(wantCidrIPs) {
			return fmt.Errorf("bad rule[*] state, expected length that having cidr_ip is %d, got length: %d", len(wantCidrIPs), len(nasSecurityGroup.IPRanges))
		}

		gotCidrIps := []string{}
		for _, ipRange := range nasSecurityGroup.IPRanges {
			gotCidrIps = append(gotCidrIps, nifcloud.ToString(ipRange.CIDRIP))
		}

		sort.Strings(wantCidrIPs)
		sort.Strings(gotCidrIps)

		for i, want := range wantCidrIPs {
			if want != gotCidrIps[i] {
				return fmt.Errorf("bad rule[*].cidr_ip state, expected %q, got: %#v", want, gotCidrIps[i])
			}
		}

		wantGroupNames := []string{groupName + "01"}
		if len(nasSecurityGroup.SecurityGroups) != len(wantGroupNames) {
			return fmt.Errorf("bad rule[*] state, expected length that having security_group_name is %d, got length: %d", len(wantGroupNames), len(nasSecurityGroup.SecurityGroups))
		}

		gotGroupNames := []string{}
		for _, securityGroup := range nasSecurityGroup.SecurityGroups {
			gotGroupNames = append(gotGroupNames, nifcloud.ToString(securityGroup.SecurityGroupName))
		}

		sort.Strings(wantGroupNames)
		sort.Strings(gotGroupNames)

		for i, want := range wantGroupNames {
			if want != gotGroupNames[i] {
				return fmt.Errorf("bad rule[*].security_group_name state, expected %q, got: %#v", want, gotGroupNames[i])
			}
		}

		return nil
	}
}

func testAccNASSecurityGroupResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).NAS

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_nas_security_group" {
			continue
		}

		res, err := svc.DescribeNASSecurityGroups(context.Background(), &nas.DescribeNASSecurityGroupsInput{
			NASSecurityGroupName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameter.NotFound.NASSecurityGroupName" {
				return nil
			}
			return fmt.Errorf("failed DescribeNASSecurityGroupsRequest: %s", err)
		}

		if len(res.NASSecurityGroups) > 0 {
			return fmt.Errorf("nasSecurityGroup (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepNASSecurityGroup(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).NAS

	res, err := svc.DescribeNASSecurityGroups(ctx, nil)
	if err != nil {
		return err
	}

	var sweepNASSecurityGroups []string
	for _, g := range res.NASSecurityGroups {
		if strings.HasPrefix(nifcloud.ToString(g.NASSecurityGroupName), prefix) {
			sweepNASSecurityGroups = append(sweepNASSecurityGroups, nifcloud.ToString(g.NASSecurityGroupName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepNASSecurityGroups {
		groupName := n
		eg.Go(func() error {
			_, err := svc.DeleteNASSecurityGroup(ctx, &nas.DeleteNASSecurityGroupInput{
				NASSecurityGroupName: nifcloud.String(groupName),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
