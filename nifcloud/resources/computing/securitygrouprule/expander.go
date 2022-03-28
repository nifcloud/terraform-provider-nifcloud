package securitygrouprule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandAuthorizeSecurityGroupIngressInputList(d *schema.ResourceData) []*computing.AuthorizeSecurityGroupIngressInput {
	protocol := d.Get("protocol").(string)
	ipPermission := types.RequestIpPermissions{
		IpProtocol:  types.IpProtocolOfIpPermissionsForAuthorizeSecurityGroupIngress(protocol),
		InOut:       types.InOutOfIpPermissionsForAuthorizeSecurityGroupIngress(d.Get("type").(string)),
		Description: nifcloud.String(d.Get("description").(string)),
	}

	if protocol == "TCP" || protocol == "UDP" {
		ipPermission.FromPort = nifcloud.Int32(int32(d.Get("from_port").(int)))
		if raw, ok := d.GetOk("to_port"); ok {
			ipPermission.ToPort = nifcloud.Int32(int32(raw.(int)))
		}
	}

	if raw, ok := d.GetOk("cidr_ip"); ok {
		ipPermission.ListOfRequestIpRanges = append(
			ipPermission.ListOfRequestIpRanges,
			types.RequestIpRanges{CidrIp: nifcloud.String(raw.(string))},
		)
	}

	if raw, ok := d.GetOk("source_security_group_name"); ok {
		ipPermission.ListOfRequestGroups = append(
			ipPermission.ListOfRequestGroups,
			types.RequestGroups{GroupName: nifcloud.String(raw.(string))},
		)
	}

	groupNames := d.Get("security_group_names").([]interface{})
	inputList := make([]*computing.AuthorizeSecurityGroupIngressInput, len(groupNames))

	for i, name := range groupNames {
		inputList[i] = &computing.AuthorizeSecurityGroupIngressInput{
			GroupName:     nifcloud.String(name.(string)),
			IpPermissions: []types.RequestIpPermissions{ipPermission},
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
	ipPermission := types.RequestIpPermissionsOfRevokeSecurityGroupIngress{
		IpProtocol: types.IpProtocolOfIpPermissionsForRevokeSecurityGroupIngress(protocol),
		InOut:      types.InOutOfIpPermissionsForRevokeSecurityGroupIngress(d.Get("type").(string)),
	}

	if protocol == "TCP" || protocol == "UDP" {
		ipPermission.FromPort = nifcloud.Int32(int32(d.Get("from_port").(int)))
		if raw, ok := d.GetOk("to_port"); ok {
			ipPermission.ToPort = nifcloud.Int32(int32(raw.(int)))
		}
	}

	if raw, ok := d.GetOk("cidr_ip"); ok {
		ipPermission.ListOfRequestIpRanges = append(
			ipPermission.ListOfRequestIpRanges,
			types.RequestIpRanges{CidrIp: nifcloud.String(raw.(string))},
		)
	}

	if raw, ok := d.GetOk("source_security_group_name"); ok {
		ipPermission.ListOfRequestGroups = append(
			ipPermission.ListOfRequestGroups,
			types.RequestGroups{GroupName: nifcloud.String(raw.(string))},
		)
	}
	groupNames := d.Get("security_group_names").([]interface{})
	inputList := make([]*computing.RevokeSecurityGroupIngressInput, len(groupNames))

	for i, name := range groupNames {
		inputList[i] = &computing.RevokeSecurityGroupIngressInput{
			GroupName:     nifcloud.String(name.(string)),
			IpPermissions: []types.RequestIpPermissionsOfRevokeSecurityGroupIngress{ipPermission},
		}
	}
	return inputList
}
