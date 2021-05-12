package nassecuritygroup

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
)

func expandCreateNASSecurityGroupInput(d *schema.ResourceData) *nas.CreateNASSecurityGroupInput {
	return &nas.CreateNASSecurityGroupInput{
		AvailabilityZone:            nifcloud.String(d.Get("availability_zone").(string)),
		NASSecurityGroupName:        nifcloud.String(d.Get("group_name").(string)),
		NASSecurityGroupDescription: nifcloud.String(d.Get("description").(string)),
	}
}

func expandAuthorizeNASSecurityGroupIngressInput(d *schema.ResourceData, rule map[string]interface{}) *nas.AuthorizeNASSecurityGroupIngressInput {
	return &nas.AuthorizeNASSecurityGroupIngressInput{
		CIDRIP:               nifcloud.String(rule["cidr_ip"].(string)),
		SecurityGroupName:    nifcloud.String(rule["security_group_name"].(string)),
		NASSecurityGroupName: nifcloud.String(d.Id()),
	}
}

func expandRevokeNASSecurityGroupIngressInput(d *schema.ResourceData, rule map[string]interface{}) *nas.RevokeNASSecurityGroupIngressInput {
	return &nas.RevokeNASSecurityGroupIngressInput{
		CIDRIP:               nifcloud.String(rule["cidr_ip"].(string)),
		SecurityGroupName:    nifcloud.String(rule["security_group_name"].(string)),
		NASSecurityGroupName: nifcloud.String(d.Id()),
	}
}

func expandDescribeNASSecurityGroupsInput(d *schema.ResourceData) *nas.DescribeNASSecurityGroupsInput {
	return &nas.DescribeNASSecurityGroupsInput{
		NASSecurityGroupName: nifcloud.String(d.Id()),
	}
}

func expandModifyNASSecurityGroupInput(d *schema.ResourceData) *nas.ModifyNASSecurityGroupInput {
	input := &nas.ModifyNASSecurityGroupInput{
		NASSecurityGroupName:        nifcloud.String(d.Id()),
		NASSecurityGroupDescription: nifcloud.String(d.Get("description").(string)),
	}

	if d.HasChange("group_name") && !d.IsNewResource() {
		input.NewNASSecurityGroupName = nifcloud.String(d.Get("group_name").(string))
	}

	return input
}

func expandDeleteNASSecurityGroupInput(d *schema.ResourceData) *nas.DeleteNASSecurityGroupInput {
	return &nas.DeleteNASSecurityGroupInput{
		NASSecurityGroupName: nifcloud.String(d.Id()),
	}
}
