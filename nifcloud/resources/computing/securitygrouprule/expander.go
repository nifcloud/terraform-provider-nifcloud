package securitygrouprule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandAuthorizeSecurityGroupIngressInputList(d *schema.ResourceData) []*computing.AuthorizeSecurityGroupIngressInput {
	protocol := d.Get("protocol").(string)
	ipPermission := computing.RequestIpPermissions{
		IpProtocol:  computing.IpProtocolOfIpPermissionsForAuthorizeSecurityGroupIngress(protocol),
		InOut:       computing.InOutOfIpPermissionsForAuthorizeSecurityGroupIngress(d.Get("type").(string)),
		Description: nifcloud.String(d.Get("description").(string)),
	}

	if protocol == "TCP" || protocol == "UDP" {
		ipPermission.FromPort = nifcloud.Int64(int64(d.Get("from_port").(int)))
		if raw, ok := d.GetOk("to_port"); ok {
			ipPermission.ToPort = nifcloud.Int64(int64(raw.(int)))
		}
	}

	if raw, ok := d.GetOk("cidr_ip"); ok {
		ipPermission.ListOfRequestIpRanges = append(
			ipPermission.ListOfRequestIpRanges,
			computing.RequestIpRanges{CidrIp: nifcloud.String(raw.(string))},
		)
	}

	if raw, ok := d.GetOk("source_security_group_name"); ok {
		ipPermission.ListOfRequestGroups = append(
			ipPermission.ListOfRequestGroups,
			computing.RequestGroups{GroupName: nifcloud.String(raw.(string))},
		)
	}

	groupNames := d.Get("security_group_names").([]interface{})
	inputList := make([]*computing.AuthorizeSecurityGroupIngressInput, len(groupNames))

	for i, name := range groupNames {
		inputList[i] = &computing.AuthorizeSecurityGroupIngressInput{
			GroupName:     nifcloud.String(name.(string)),
			IpPermissions: []computing.RequestIpPermissions{ipPermission},
		}
	}
	return inputList
}

func expandDescribeSecurityGroupsInput(d *schema.ResourceData) *computing.DescribeSecurityGroupsInput {
	groupNames := d.Get("security_group_names").([]interface{})
	groupName := make([]string, len(groupNames))

	for i, n := range groupNames {
		groupName[i] = n.(string)
	}

	return &computing.DescribeSecurityGroupsInput{
		GroupName: groupName,
	}
}

func expandRevokeSecurityGroupIngressInputList(d *schema.ResourceData) []*computing.RevokeSecurityGroupIngressInput {
	protocol := d.Get("protocol").(string)
	ipPermission := computing.RequestIpPermissionsOfRevokeSecurityGroupIngress{
		IpProtocol: computing.IpProtocolOfIpPermissionsForRevokeSecurityGroupIngress(protocol),
		InOut:      computing.InOutOfIpPermissionsForRevokeSecurityGroupIngress(d.Get("type").(string)),
	}

	if protocol == "TCP" || protocol == "UDP" {
		ipPermission.FromPort = nifcloud.Int64(int64(d.Get("from_port").(int)))
		if raw, ok := d.GetOk("to_port"); ok {
			ipPermission.ToPort = nifcloud.Int64(int64(raw.(int)))
		}
	}

	if raw, ok := d.GetOk("cidr_ip"); ok {
		ipPermission.ListOfRequestIpRanges = append(
			ipPermission.ListOfRequestIpRanges,
			computing.RequestIpRanges{CidrIp: nifcloud.String(raw.(string))},
		)
	}

	if raw, ok := d.GetOk("source_security_group_name"); ok {
		ipPermission.ListOfRequestGroups = append(
			ipPermission.ListOfRequestGroups,
			computing.RequestGroups{GroupName: nifcloud.String(raw.(string))},
		)
	}
	groupNames := d.Get("security_group_names").([]interface{})
	inputList := make([]*computing.RevokeSecurityGroupIngressInput, len(groupNames))

	for i, name := range groupNames {
		inputList[i] = &computing.RevokeSecurityGroupIngressInput{
			GroupName:     nifcloud.String(name.(string)),
			IpPermissions: []computing.RequestIpPermissionsOfRevokeSecurityGroupIngress{ipPermission},
		}
	}
	return inputList
}
