package securitygroup

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandCreateSecurityGroupInput(d *schema.ResourceData) *computing.CreateSecurityGroupInput {
	return &computing.CreateSecurityGroupInput{
		GroupName:        nifcloud.String(d.Get("group_name").(string)),
		GroupDescription: nifcloud.String(d.Get("description").(string)),
		Placement: &computing.RequestPlacementOfCreateSecurityGroup{
			AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		},
	}
}

func expandUpdateSecurityGroupInputForLogLimit(d *schema.ResourceData) *computing.UpdateSecurityGroupInput {
	return &computing.UpdateSecurityGroupInput{
		GroupName:           nifcloud.String(d.Id()),
		GroupLogLimitUpdate: nifcloud.Int64(int64(d.Get("log_limit").(int))),
	}
}

func expandUpdateSecurityGroupInputForName(d *schema.ResourceData) *computing.UpdateSecurityGroupInput {
	before, after := d.GetChange("group_name")

	return &computing.UpdateSecurityGroupInput{
		GroupName:       nifcloud.String(before.(string)),
		GroupNameUpdate: nifcloud.String(after.(string)),
	}
}
func expandUpdateSecurityGroupInputForDescription(d *schema.ResourceData) *computing.UpdateSecurityGroupInput {
	return &computing.UpdateSecurityGroupInput{
		GroupName:              nifcloud.String(d.Id()),
		GroupDescriptionUpdate: nifcloud.String(d.Get("description").(string)),
	}
}

func expandDescribeSecurityGroupsInput(d *schema.ResourceData) *computing.DescribeSecurityGroupsInput {
	return &computing.DescribeSecurityGroupsInput{
		GroupName: []string{d.Id()},
	}
}

func expandDeleteSecurityGroupInput(d *schema.ResourceData) *computing.DeleteSecurityGroupInput {
	return &computing.DeleteSecurityGroupInput{
		GroupName: nifcloud.String(d.Id()),
	}
}

func expandRevokeSecurityGroupIngressInput(
	d *schema.ResourceData,
	ipPermissions []computing.RequestIpPermissionsOfRevokeSecurityGroupIngress,
) *computing.RevokeSecurityGroupIngressInput {
	return &computing.RevokeSecurityGroupIngressInput{
		GroupName:     nifcloud.String(d.Id()),
		IpPermissions: ipPermissions,
	}
}
