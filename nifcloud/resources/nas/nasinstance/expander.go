package nasinstance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas/types"
)

func expandCreateNASInstanceInput(d *schema.ResourceData) *nas.CreateNASInstanceInput {
	return &nas.CreateNASInstanceInput{
		AllocatedStorage:       nifcloud.Int32(int32(d.Get("allocated_storage").(int))),
		AvailabilityZone:       nifcloud.String(d.Get("availability_zone").(string)),
		MasterPrivateAddress:   nifcloud.String(d.Get("private_ip_address").(string) + d.Get("private_ip_address_subnet_mask").(string)),
		MasterUserPassword:     nifcloud.String(d.Get("master_user_password").(string)),
		MasterUsername:         nifcloud.String(d.Get("master_username").(string)),
		NASInstanceDescription: nifcloud.String(d.Get("description").(string)),
		NASInstanceIdentifier:  nifcloud.String(d.Get("identifier").(string)),
		NASInstanceType:        nifcloud.Int32(int32(d.Get("type").(int))),
		NASSecurityGroups:      []string{d.Get("nas_security_group_name").(string)},
		NetworkId:              nifcloud.String(d.Get("network_id").(string)),
		Protocol:               types.ProtocolOfCreateNASInstanceRequest(d.Get("protocol").(string)),
	}
}

func expandDescribeNASInstancesInput(d *schema.ResourceData) *nas.DescribeNASInstancesInput {
	return &nas.DescribeNASInstancesInput{
		NASInstanceIdentifier: nifcloud.String(d.Id()),
	}
}

func expandModifyNASInstanceInput(d *schema.ResourceData) *nas.ModifyNASInstanceInput {
	input := &nas.ModifyNASInstanceInput{
		AllocatedStorage:       nifcloud.Int32(int32(d.Get("allocated_storage").(int))),
		NASInstanceIdentifier:  nifcloud.String(d.Id()),
		NASInstanceDescription: nifcloud.String(d.Get("description").(string)),
		NASSecurityGroups:      []string{d.Get("nas_security_group_name").(string)},
		NetworkId:              nifcloud.String(d.Get("network_id").(string)),
		MasterPrivateAddress:   nifcloud.String(d.Get("private_ip_address").(string) + d.Get("private_ip_address_subnet_mask").(string)),
		NoRootSquash:           nifcloud.Bool(d.Get("no_root_squash").(bool)),
	}

	if d.Get("protocol").(string) == "cifs" {
		authenticationType := d.Get("authentication_type").(int)
		input.AuthenticationType = nifcloud.Int32(int32(authenticationType))
		input.MasterUserPassword = nifcloud.String(d.Get("master_user_password").(string))
	}

	if d.HasChange("identifier") && !d.IsNewResource() {
		input.NewNASInstanceIdentifier = nifcloud.String(d.Get("identifier").(string))
	}

	return input
}

func expandDeleteNASInstanceInput(d *schema.ResourceData) *nas.DeleteNASInstanceInput {
	return &nas.DeleteNASInstanceInput{
		NASInstanceIdentifier: nifcloud.String(d.Id()),
	}
}
