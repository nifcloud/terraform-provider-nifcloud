package nasinstance

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
)

func expandCreateNASInstanceInput(d *schema.ResourceData) *nas.CreateNASInstanceInput {
	return &nas.CreateNASInstanceInput{
		AllocatedStorage:       nifcloud.Int64(int64(d.Get("allocated_storage").(int))),
		AvailabilityZone:       nifcloud.String(d.Get("availability_zone").(string)),
		MasterPrivateAddress:   nifcloud.String(d.Get("private_ip_address").(string) + d.Get("private_ip_address_subnet_mask").(string)),
		MasterUserPassword:     nifcloud.String(d.Get("master_user_password").(string)),
		MasterUsername:         nifcloud.String(d.Get("master_username").(string)),
		NASInstanceDescription: nifcloud.String(d.Get("description").(string)),
		NASInstanceIdentifier:  nifcloud.String(d.Get("identifier").(string)),
		NASInstanceType:        nifcloud.Int64(int64(d.Get("type").(int))),
		NASSecurityGroups:      []string{d.Get("nas_security_group_name").(string)},
		NetworkId:              nifcloud.String(d.Get("network_id").(string)),
		Protocol:               nifcloud.String(d.Get("protocol").(string)),
	}
}

func expandDescribeNASInstancesInput(d *schema.ResourceData) *nas.DescribeNASInstancesInput {
	return &nas.DescribeNASInstancesInput{
		NASInstanceIdentifier: nifcloud.String(d.Id()),
	}
}

func expandModifyNASInstanceInput(d *schema.ResourceData) *nas.ModifyNASInstanceInput {
	input := &nas.ModifyNASInstanceInput{
		AllocatedStorage:       nifcloud.Int64(int64(d.Get("allocated_storage").(int))),
		NASInstanceIdentifier:  nifcloud.String(d.Id()),
		NASInstanceDescription: nifcloud.String(d.Get("description").(string)),
		NASSecurityGroups:      []string{d.Get("nas_security_group_name").(string)},
		NetworkId:              nifcloud.String(d.Get("network_id").(string)),
		MasterPrivateAddress:   nifcloud.String(d.Get("private_ip_address").(string) + d.Get("private_ip_address_subnet_mask").(string)),
		NoRootSquash:           nifcloud.String(strconv.FormatBool(d.Get("no_root_squash").(bool))),
	}

	if d.Get("protocol").(string) == "cifs" {
		authenticationType := d.Get("authentication_type").(int)
		if authenticationType == 1 {
			domainControllers := []nas.RequestDomainControllers{}
			for _, controller := range d.Get("domain_controllers").(*schema.Set).List() {
				if v, ok := controller.(map[string]interface{}); ok {
					domainController := nas.RequestDomainControllers{}
					if row, ok := v["hostname"]; ok {
						domainController.Hostname = nifcloud.String(row.(string))
					}
					if row, ok := v["ip_address"]; ok {
						domainController.IPAddress = nifcloud.String(row.(string))
					}
					domainControllers = append(domainControllers, domainController)
				}
			}
			input.DirectoryServiceAdministratorName = nifcloud.String(d.Get("directory_service_administrator_name").(string))
			input.DirectoryServiceAdministratorPassword = nifcloud.String(d.Get("directory_service_administrator_password").(string))
			input.DirectoryServiceDomainName = nifcloud.String(d.Get("directory_service_domain_name").(string))
			input.DomainControllers = domainControllers
		}

		input.AuthenticationType = nifcloud.Int64(int64(authenticationType))
		input.MasterUserPassword = nifcloud.String(d.Get("master_user_password").(string))
	}

	if d.HasChange("identifier") && !d.IsNewResource() {
		input.NewNASInstanceIdentifier = nifcloud.String(d.Get("identifier").(string))
	}

	return input
}

func expandDeleteNASInstanceInput(d *schema.ResourceData) *nas.DeleteNASInstanceInput {
	return &nas.DeleteNASInstanceInput{
		NASInstanceIdentifier:                 nifcloud.String(d.Id()),
		DirectoryServiceAdministratorName:     nifcloud.String(d.Get("directory_service_administrator_name").(string)),
		DirectoryServiceAdministratorPassword: nifcloud.String(d.Get("directory_service_administrator_password").(string)),
	}
}
