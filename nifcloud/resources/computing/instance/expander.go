package instance

import (
	"encoding/base64"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandRunInstancesInput(d *schema.ResourceData) *computing.RunInstancesInput {
	var networkInterface []types.RequestNetworkInterface
	for _, ni := range d.Get("network_interface").(*schema.Set).List() {
		if v, ok := ni.(map[string]interface{}); ok {
			n := types.RequestNetworkInterface{}
			if row, ok := v["network_id"]; ok {
				// net-MULTI_IP_ADDRESS cannot be specified during instance creation, so replace it with net-COMMON_GLOBAL
				// Multi IP address group association will be performed after instance creation
				if row.(string) == "net-MULTI_IP_ADDRESS" {
					n.NetworkId = nifcloud.String("net-COMMON_GLOBAL")
				} else {
					n.NetworkId = nifcloud.String(row.(string))
				}
			}
			if row, ok := v["network_name"]; ok {
				n.NetworkName = nifcloud.String(row.(string))
			}
			if row, ok := v["ip_address"]; ok {
				n.IpAddress = nifcloud.String(row.(string))
			}
			if nifcloud.ToString(n.NetworkId) != "" || nifcloud.ToString(n.NetworkName) != "" {
				networkInterface = append(networkInterface, n)
			}
		}
	}

	var securityGroup []string
	if row, ok := d.GetOk("security_group"); ok {
		securityGroup = append(securityGroup, row.(string))
	}

	input := &computing.RunInstancesInput{
		InstanceId:    nifcloud.String(d.Get("instance_id").(string)),
		ImageId:       nifcloud.String(d.Get("image_id").(string)),
		KeyName:       nifcloud.String(d.Get("key_name").(string)),
		SecurityGroup: securityGroup,
		InstanceType:  types.InstanceTypeOfRunInstancesRequest(d.Get("instance_type").(string)),
		Placement: &types.RequestPlacement{
			AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		},
		DisableApiTermination: nifcloud.Bool(d.Get("disable_api_termination").(bool)),
		AccountingType:        types.AccountingTypeOfRunInstancesRequest(d.Get("accounting_type").(string)),
		Description:           nifcloud.String(d.Get("description").(string)),
		Admin:                 nifcloud.String(d.Get("admin").(string)),
		Password:              nifcloud.String(d.Get("password").(string)),
		Agreement:             nifcloud.Bool(true),
		UserData: &types.RequestUserData{
			Content:  nifcloud.String(base64.StdEncoding.EncodeToString([]byte(d.Get("user_data").(string)))),
			Encoding: nifcloud.String("base64"),
		},
		NetworkInterface: networkInterface,
	}

	if raw, ok := d.GetOk("license_name"); ok {
		input.License = []types.RequestLicense{
			{
				LicenseName: types.LicenseNameOfLicenseForRunInstances(raw.(string)),
				LicenseNum:  nifcloud.String(strconv.Itoa(d.Get("license_num").(int))),
			},
		}
	}
	return input
}

func expandDescribeInstancesInput(d *schema.ResourceData) *computing.DescribeInstancesInput {
	return &computing.DescribeInstancesInput{
		InstanceId: []string{d.Id()},
	}
}

func expandDescribeInstanceAttributeInputWithDisableAPITermination(d *schema.ResourceData) *computing.DescribeInstanceAttributeInput {
	return &computing.DescribeInstanceAttributeInput{
		InstanceId: nifcloud.String(d.Id()),
		Attribute:  types.AttributeOfDescribeInstanceAttributeRequestDisableApiTermination,
	}
}

func expandStopInstancesInput(d *schema.ResourceData, force bool) *computing.StopInstancesInput {
	return &computing.StopInstancesInput{
		InstanceId: []string{d.Id()},
		Force:      nifcloud.Bool(force),
	}
}

func expandTerminateInstancesInput(d *schema.ResourceData) *computing.TerminateInstancesInput {
	return &computing.TerminateInstancesInput{
		InstanceId: []string{d.Id()},
	}
}

func expandModifyInstanceAttributeInputForAccountingType(d *schema.ResourceData) *computing.ModifyInstanceAttributeInput {
	return &computing.ModifyInstanceAttributeInput{
		InstanceId: nifcloud.String(d.Id()),
		Attribute:  types.AttributeOfModifyInstanceAttributeRequestAccountingType,
		Value:      nifcloud.String(d.Get("accounting_type").(string)),
	}
}

func expandModifyInstanceAttributeInputForDescription(d *schema.ResourceData) *computing.ModifyInstanceAttributeInput {
	return &computing.ModifyInstanceAttributeInput{
		InstanceId: nifcloud.String(d.Id()),
		Attribute:  types.AttributeOfModifyInstanceAttributeRequestDescription,
		Value:      nifcloud.String(d.Get("description").(string)),
	}
}

func expandModifyInstanceAttributeInputForDisableAPITermination(d *schema.ResourceData) *computing.ModifyInstanceAttributeInput {
	return &computing.ModifyInstanceAttributeInput{
		InstanceId: nifcloud.String(d.Id()),
		Attribute:  types.AttributeOfModifyInstanceAttributeRequestDisableApiTermination,
		Value:      nifcloud.String(strconv.FormatBool(d.Get("disable_api_termination").(bool))),
	}
}

func expandModifyInstanceAttributeInputForInstanceID(d *schema.ResourceData) *computing.ModifyInstanceAttributeInput {
	before, after := d.GetChange("instance_id")

	return &computing.ModifyInstanceAttributeInput{
		InstanceId: nifcloud.String(before.(string)),
		Attribute:  types.AttributeOfModifyInstanceAttributeRequestInstanceName,
		Value:      nifcloud.String(after.(string)),
	}
}
func expandModifyInstanceAttributeInputForInstanceType(d *schema.ResourceData) *computing.ModifyInstanceAttributeInput {
	return &computing.ModifyInstanceAttributeInput{
		InstanceId: nifcloud.String(d.Id()),
		Attribute:  types.AttributeOfModifyInstanceAttributeRequestInstanceType,
		Value:      nifcloud.String(d.Get("instance_type").(string)),
	}
}

func expandModifyInstanceAttributeInputForSecurityGroup(d *schema.ResourceData) *computing.ModifyInstanceAttributeInput {
	return &computing.ModifyInstanceAttributeInput{
		InstanceId: nifcloud.String(d.Id()),
		Attribute:  types.AttributeOfModifyInstanceAttributeRequestGroupId,
		Value:      nifcloud.String(d.Get("security_group").(string)),
	}
}

func expandNiftyUpdateInstanceNetworkInterfacesInput(d *schema.ResourceData) *computing.NiftyUpdateInstanceNetworkInterfacesInput {
	var networkInterface []types.RequestNetworkInterface
	for _, ni := range d.Get("network_interface").(*schema.Set).List() {
		if v, ok := ni.(map[string]interface{}); ok {
			n := types.RequestNetworkInterface{}
			if row, ok := v["network_id"]; ok {
				n.NetworkId = nifcloud.String(row.(string))
			}
			if row, ok := v["network_name"]; ok {
				n.NetworkName = nifcloud.String(row.(string))
			}
			if row, ok := v["ip_address"]; ok {
				n.IpAddress = nifcloud.String(row.(string))
			}

			if nifcloud.ToString(n.NetworkId) != "" || nifcloud.ToString(n.NetworkName) != "" {
				networkInterface = append(networkInterface, n)
			}
		}
	}

	return &computing.NiftyUpdateInstanceNetworkInterfacesInput{
		InstanceId:       nifcloud.String(d.Id()),
		NetworkInterface: networkInterface,
	}
}

func expandAttachNetworkInterfaceInput(d *schema.ResourceData, networkInterfaceID string) *computing.AttachNetworkInterfaceInput {
	return &computing.AttachNetworkInterfaceInput{
		InstanceId:         nifcloud.String(d.Id()),
		NetworkInterfaceId: nifcloud.String(networkInterfaceID),
		NiftyReboot:        types.NiftyRebootOfAttachNetworkInterfaceRequestForce,
	}
}

func expandDetachNetworkInterfaceInput(d *schema.ResourceData, attachmentID string) *computing.DetachNetworkInterfaceInput {
	return &computing.DetachNetworkInterfaceInput{
		AttachmentId: nifcloud.String(attachmentID),
		NiftyReboot:  types.NiftyRebootOfDetachNetworkInterfaceRequestForce,
	}
}

func expandDeregisterInstancesFromSecurityGroupInput(d *schema.ResourceData) *computing.DeregisterInstancesFromSecurityGroupInput {
	groupName, _ := d.GetChange("security_group")

	return &computing.DeregisterInstancesFromSecurityGroupInput{
		InstanceId: []string{d.Id()},
		GroupName:  nifcloud.String(groupName.(string)),
	}
}

func expandAssociateMultiIpAddressGroupInput(multiIPAddressGroupID, instanceUniqueID string) *computing.AssociateMultiIpAddressGroupInput {
	return &computing.AssociateMultiIpAddressGroupInput{
		MultiIpAddressGroupId: nifcloud.String(multiIPAddressGroupID),
		InstanceUniqueId:      nifcloud.String(instanceUniqueID),
		NiftyReboot:           types.NiftyRebootOfAssociateMultiIpAddressGroupRequestFalse,
	}
}

func expandDisassociateMultiIpAddressGroupInput(multiIPAddressGroupID, instanceUniqueID string) *computing.DisassociateMultiIpAddressGroupInput {
	return &computing.DisassociateMultiIpAddressGroupInput{
		MultiIpAddressGroupId: nifcloud.String(multiIPAddressGroupID),
		InstanceUniqueId:      nifcloud.String(instanceUniqueID),
		NiftyReboot:           types.NiftyRebootOfDisassociateMultiIpAddressGroupRequestFalse,
	}
}

func expandStartInstancesInputWithMultiIPAddressConfigurationUserData(d *schema.ResourceData) *computing.StartInstancesInput {
	return &computing.StartInstancesInput{
		InstanceId: []string{d.Get("instance_id").(string)},
		UserData: &types.RequestUserData{
			Content:  nifcloud.String(base64.StdEncoding.EncodeToString([]byte(d.Get("multi_ip_address_configuration_user_data").(string)))),
			Encoding: nifcloud.String("base64"),
		},
	}
}
