package dbsecuritygroup

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
)

func expandCreateDBSecurityGroupInput(d *schema.ResourceData) *rdb.CreateDBSecurityGroupInput {
	return &rdb.CreateDBSecurityGroupInput{
		NiftyAvailabilityZone:      nifcloud.String(d.Get("availability_zone").(string)),
		DBSecurityGroupName:        nifcloud.String(d.Get("group_name").(string)),
		DBSecurityGroupDescription: nifcloud.String(d.Get("description").(string)),
	}
}

func expandAuthorizeDBSecurityGroupIngressInput(d *schema.ResourceData, rule map[string]interface{}) *rdb.AuthorizeDBSecurityGroupIngressInput {
	return &rdb.AuthorizeDBSecurityGroupIngressInput{
		CIDRIP:               nifcloud.String(rule["cidr_ip"].(string)),
		EC2SecurityGroupName: nifcloud.String(rule["security_group_name"].(string)),
		DBSecurityGroupName:  nifcloud.String(d.Id()),
	}
}

func expandRevokeDBSecurityGroupIngressInput(d *schema.ResourceData, rule map[string]interface{}) *rdb.RevokeDBSecurityGroupIngressInput {
	return &rdb.RevokeDBSecurityGroupIngressInput{
		CIDRIP:               nifcloud.String(rule["cidr_ip"].(string)),
		EC2SecurityGroupName: nifcloud.String(rule["security_group_name"].(string)),
		DBSecurityGroupName:  nifcloud.String(d.Id()),
	}
}

func expandDescribeDBSecurityGroupsInput(d *schema.ResourceData) *rdb.DescribeDBSecurityGroupsInput {
	return &rdb.DescribeDBSecurityGroupsInput{
		DBSecurityGroupName: nifcloud.String(d.Id()),
	}
}

func expandDeleteDBSecurityGroupInput(d *schema.ResourceData) *rdb.DeleteDBSecurityGroupInput {
	return &rdb.DeleteDBSecurityGroupInput{
		DBSecurityGroupName: nifcloud.String(d.Id()),
	}
}
