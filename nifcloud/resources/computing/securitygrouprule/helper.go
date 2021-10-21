package securitygrouprule

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

var mutexKV = mutexkv.NewMutexKV()

type securityGroupNotFound struct {
	name           string
	securityGroups []computing.SecurityGroupInfo
}

func (err securityGroupNotFound) Error() string {
	if err.securityGroups == nil {
		return fmt.Sprintf("No security group with name %q", err.name)
	}
	return fmt.Sprintf("Expected to find one security group with name %q, got: %#v",
		err.name, err.securityGroups)
}

func checkSecurityGroupExist(securityGroupInfo []computing.SecurityGroupInfo, groupName string) error {
	for _, s := range securityGroupInfo {
		if nifcloud.StringValue(s.GroupName) == groupName {
			return nil
		}
	}
	return securityGroupNotFound{name: groupName, securityGroups: securityGroupInfo}
}

func idHash(inputList []*computing.AuthorizeSecurityGroupIngressInput) string {
	var buf bytes.Buffer
	ip := inputList[0].IpPermissions[0]
	if ip.FromPort != nil && nifcloud.Int64Value(ip.FromPort) > 0 {
		buf.WriteString(fmt.Sprintf("%d-", nifcloud.Int64Value(ip.FromPort)))
	}
	if ip.ToPort != nil && nifcloud.Int64Value(ip.ToPort) > 0 {
		buf.WriteString(fmt.Sprintf("%d-", nifcloud.Int64Value(ip.ToPort)))
	}
	buf.WriteString(fmt.Sprintf("%s-", ip.IpProtocol))
	buf.WriteString(fmt.Sprintf("%s-", ip.InOut))
	buf.WriteString(fmt.Sprintf("%s-", nifcloud.StringValue(ip.Description)))

	s := make([]string, len(inputList))
	for i, input := range inputList {
		s[i] = nifcloud.StringValue(input.GroupName)
		sort.Strings(s)
	}

	for _, v := range s {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	if ip.ListOfRequestIpRanges != nil && len(ip.ListOfRequestIpRanges) > 0 {
		buf.WriteString(fmt.Sprintf("%s-", ip.ListOfRequestIpRanges[0]))

	}
	if ip.ListOfRequestGroups != nil && len(ip.ListOfRequestGroups) > 0 {
		buf.WriteString(fmt.Sprintf("%s-", ip.ListOfRequestGroups[0]))
	}

	hashcode := 0
	v := int(crc32.ChecksumIEEE(buf.Bytes()))
	if v >= 0 {
		hashcode = v
	}
	if -v >= 0 {
		hashcode = -v
	}
	return fmt.Sprintf("sgrule-%d", hashcode)
}

func validateSecurityGroupRuleImportString(importStr string) ([]string, error) {
	// example: IN_TCP_8000_8000_10.0.3.0/24_example

	importParts := strings.Split(importStr, "_")
	errStr := "unexpected format of import string (%q), expected SECURITYGROUPNAME_TYPE_PROTOCOL_FROMPORT_TOPORT_SOURCE: %s"
	if len(importParts) < 6 {
		return nil, fmt.Errorf(errStr, importStr, "invalid parts")
	}

	ruleType := importParts[0]
	protocol := importParts[1]
	fromPort := importParts[2]
	toPort := importParts[3]
	source := importParts[4]
	sgName := importParts[5:]

	if len(sgName) == 0 {
		return nil, fmt.Errorf(errStr, importStr, "security group name must be required")
	}

	if ruleType != "IN" && ruleType != "OUT" {
		return nil, fmt.Errorf(errStr, importStr, "expecting 'IN' or 'OUT'")
	}

	if protocol != "ANY" &&
		protocol != "TCP" &&
		protocol != "UDP" &&
		protocol != "ICMP" &&
		protocol != "GRE" &&
		protocol != "AH" &&
		protocol != "VRRP" &&
		protocol != "ICMPv6-all" {
		return nil, fmt.Errorf(errStr, importStr, "protocol must be ANY/TCP/UDP/ICMP/GRE/AH/VRRP/ICMPv6-all")
	}

	if fromPort != "-" && toPort != "-" {
		if p1, err := strconv.Atoi(fromPort); err != nil {
			return nil, fmt.Errorf(errStr, importStr, "invalid port")
		} else if p2, err := strconv.Atoi(toPort); err != nil || p2 < p1 {
			return nil, fmt.Errorf(errStr, importStr, "invalid port")
		}
	}

	if source == "" {
		return nil, fmt.Errorf(errStr, importStr, "source must be required")
	}

	return importParts, nil
}

func populateSecurityGroupRuleFromImport(d *schema.ResourceData, importParts []string) error {
	ruleType := importParts[0]
	protocol := importParts[1]
	fromPort := importParts[2]
	toPort := importParts[3]
	source := importParts[4]
	sgName := importParts[5:]

	if err := d.Set("type", ruleType); err != nil {
		return err
	}

	if err := d.Set("protocol", protocol); err != nil {
		return err
	}

	if fromPort != "-" {
		p, err := strconv.Atoi(fromPort)
		if err != nil {
			return err
		}

		if err := d.Set("from_port", p); err != nil {
			return err
		}
	}

	if toPort != "-" {
		p, err := strconv.Atoi(toPort)
		if err != nil {
			return err
		}

		if err := d.Set("to_port", p); err != nil {
			return err
		}
	}

	if strings.Contains(source, ".") || strings.Contains(source, ":") {
		if err := d.Set("cidr_ip", source); err != nil {
			return err
		}
	} else {
		if err := d.Set("source_security_group_name", source); err != nil {
			return err
		}
	}

	if err := d.Set("security_group_names", sgName); err != nil {
		return err
	}
	return nil
}
