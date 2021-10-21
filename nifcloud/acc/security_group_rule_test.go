package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func TestAcc_SecurityGroupRule_Cidr(t *testing.T) {
	var securityGroupRule computing.IpPermissions

	resourceName := "nifcloud_security_group_rule.basic_cidr"

	randName := prefix + acctest.RandString(6)

	fwName := randName + "1"
	fwNameUpd := randName + "2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccSecurityGroupRuleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRule(t, "testdata/security_group_rule.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(resourceName, &securityGroupRule, fwName),
					testAccCheckSecurityGroupRuleValuesWithCidr(&securityGroupRule),
					resource.TestCheckResourceAttr(resourceName, "protocol", "ANY"),
					resource.TestCheckResourceAttr(resourceName, "type", "OUT"),
					resource.TestCheckResourceAttr(resourceName, "cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "security_group_names.0", fwName),
				),
			},
			{
				Config: testAccSecurityGroupRule(t, "testdata/security_group_rule_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(resourceName, &securityGroupRule, fwNameUpd),
					testAccCheckSecurityGroupRuleValuesWithCidr(&securityGroupRule),
					resource.TestCheckResourceAttr(resourceName, "protocol", "ANY"),
					resource.TestCheckResourceAttr(resourceName, "type", "OUT"),
					resource.TestCheckResourceAttr(resourceName, "cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "security_group_names.0", fwNameUpd),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccSecurityGroupRuleImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_SecurityGroupRule_Source(t *testing.T) {
	var securityGroupRule computing.IpPermissions

	resourceName := "nifcloud_security_group_rule.basic_source"

	randName := prefix + acctest.RandString(6)

	fwName := randName + "3"
	fwNameUpd := randName + "4"
	fwNameSrc := randName + "5"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccSecurityGroupRuleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRule(t, "testdata/security_group_rule.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(resourceName, &securityGroupRule, fwName),
					testAccCheckSecurityGroupRuleValuesWithSource(&securityGroupRule, fwNameSrc),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "type", "IN"),
					resource.TestCheckResourceAttr(resourceName, "from_port", "1"),
					resource.TestCheckResourceAttr(resourceName, "to_port", "65535"),
					resource.TestCheckResourceAttr(resourceName, "source_security_group_name", fwNameSrc),
					resource.TestCheckResourceAttr(resourceName, "security_group_names.0", fwName),
				),
			},
			{
				Config: testAccSecurityGroupRule(t, "testdata/security_group_rule_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(resourceName, &securityGroupRule, fwNameUpd),
					testAccCheckSecurityGroupRuleValuesWithSource(&securityGroupRule, fwNameSrc),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "type", "IN"),
					resource.TestCheckResourceAttr(resourceName, "from_port", "1"),
					resource.TestCheckResourceAttr(resourceName, "to_port", "65535"),
					resource.TestCheckResourceAttr(resourceName, "source_security_group_name", fwNameSrc),
					resource.TestCheckResourceAttr(resourceName, "security_group_names.0", fwNameUpd),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccSecurityGroupRuleImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityGroupRule(t *testing.T, fileName, groupName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		groupName+"1",
		groupName+"2",
		groupName+"3",
		groupName+"4",
		groupName+"5",
	)
}

func testAccCheckSecurityGroupRuleExists(n string, securityGroupRule *computing.IpPermissions, groupName string) resource.TestCheckFunc {
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
			GroupName: []string{groupName},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.SecurityGroupInfo) == 0 {
			return fmt.Errorf("securityGroup does not found in cloud: %s", groupName)
		}

		foundSecurityGroup := res.SecurityGroupInfo[0]
		if nifcloud.StringValue(foundSecurityGroup.GroupName) != groupName {
			return fmt.Errorf("securityGroup does not found in cloud: %s", groupName)
		}

		foundSecurityGroupRule := foundSecurityGroup.IpPermissions[0]

		*securityGroupRule = foundSecurityGroupRule
		return nil
	}
}

func testAccCheckSecurityGroupRuleValuesWithCidr(securityGroupRule *computing.IpPermissions) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(securityGroupRule.Description) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", securityGroupRule.Description)
		}

		if nifcloud.StringValue(securityGroupRule.InOut) != "OUT" {
			return fmt.Errorf("bad type state,  expected \"OUT\", got: %#v", securityGroupRule.InOut)
		}

		if nifcloud.StringValue(securityGroupRule.IpProtocol) != "ANY" {
			return fmt.Errorf("bad protocol state,  expected \"ANY\", got: %#v", securityGroupRule.IpProtocol)
		}

		if nifcloud.StringValue(securityGroupRule.IpRanges[0].CidrIp) != "0.0.0.0/0" {
			return fmt.Errorf("bad cidr_ip state,  expected \"0.0.0.0/0\", got: %#v", securityGroupRule.IpRanges[0].CidrIp)
		}
		return nil
	}
}

func testAccCheckSecurityGroupRuleValuesWithSource(securityGroupRule *computing.IpPermissions, sourceGroupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(securityGroupRule.Description) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", securityGroupRule.Description)
		}

		if nifcloud.StringValue(securityGroupRule.InOut) != "IN" {
			return fmt.Errorf("bad type state,  expected \"IN\", got: %#v", securityGroupRule.InOut)
		}

		if nifcloud.StringValue(securityGroupRule.IpProtocol) != "TCP" {
			return fmt.Errorf("bad protocol state,  expected \"TCP\", got: %#v", securityGroupRule.IpProtocol)
		}

		if nifcloud.Int64Value(securityGroupRule.FromPort) != 1 {
			return fmt.Errorf("bad from_port state,  expected \"1\", got: %#v", securityGroupRule.FromPort)
		}

		if nifcloud.Int64Value(securityGroupRule.ToPort) != 65535 {
			return fmt.Errorf("bad to_port state,  expected \"65535\", got: %#v", securityGroupRule.ToPort)
		}

		if nifcloud.StringValue(securityGroupRule.Groups[0].GroupName) != sourceGroupName {
			return fmt.Errorf("bad source_security_group_name state,  expected \"%s\", got: %#v", sourceGroupName, securityGroupRule.Groups[0].GroupName)
		}
		return nil
	}
}
func testAccSecurityGroupRuleResourceDestroy(s *terraform.State) error {
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
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroup" {
				return nil
			}
			return fmt.Errorf("failed DescribeSecurityGroupsRequest: %s", err)
		}

		if len(res.SecurityGroupInfo) > 0 {
			return fmt.Errorf("securityGroup (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccSecurityGroupRuleImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ruleType := rs.Primary.Attributes["type"]
		protocol := rs.Primary.Attributes["protocol"]

		var parts []string
		parts = append(parts, ruleType)
		parts = append(parts, protocol)

		if fromPort, ok := rs.Primary.Attributes["from_port"]; ok {
			parts = append(parts, fromPort)
		} else {
			parts = append(parts, "-")
		}

		if toPort, ok := rs.Primary.Attributes["to_port"]; ok {
			parts = append(parts, toPort)
		} else {
			parts = append(parts, "-")
		}

		if sgSource, ok := rs.Primary.Attributes["source_security_group_name"]; ok {
			parts = append(parts, sgSource)
		}

		if cidrIP, ok := rs.Primary.Attributes["cidr_ip"]; ok {
			parts = append(parts, cidrIP)
		}

		if countStr, ok := rs.Primary.Attributes[fmt.Sprintf("%s.#", "security_group_names")]; ok && countStr != "0" {
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return "", err
			}
			for i := 0; i < count; i++ {
				parts = append(parts, rs.Primary.Attributes[fmt.Sprintf("%s.%d", "security_group_names", i)])
			}
		}

		id := strings.Join(parts, "_")
		return id, nil
	}
}
