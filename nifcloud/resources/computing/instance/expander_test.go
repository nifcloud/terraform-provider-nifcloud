package instance

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandRunInstancesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":         "test_accounting_type",
		"admin":                   "test_admin",
		"availability_zone":       "test_availability_zone",
		"description":             "test_description",
		"disable_api_termination": false,
		"image_id":                "test_image_id",
		"instance_id":             "test_instance_id",
		"instance_type":           "test_instance_type",
		"key_name":                "test_key_name",
		"license_name":            "test_license_name",
		"license_num":             200,
		"network_interface": []interface{}{map[string]interface{}{
			"network_id":   "test_network_id",
			"network_name": "test_network_name",
			"ip_address":   "test_ip_address",
		}},
		"password":       "test_password",
		"security_group": "test_security_group",
		"user_data":      "test_user_data",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.RunInstancesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.RunInstancesInput{
				InstanceId:    nifcloud.String("test_instance_id"),
				ImageId:       nifcloud.String("test_image_id"),
				KeyName:       nifcloud.String("test_key_name"),
				SecurityGroup: []string{"test_security_group"},
				InstanceType:  types.InstanceTypeOfRunInstancesRequest("test_instance_type"),
				Placement: &types.RequestPlacement{
					AvailabilityZone: nifcloud.String(("test_availability_zone")),
				},
				DisableApiTermination: nifcloud.Bool(false),
				AccountingType:        types.AccountingTypeOfRunInstancesRequest("test_accounting_type"),
				Description:           nifcloud.String("test_description"),
				Admin:                 nifcloud.String("test_admin"),
				Password:              nifcloud.String("test_password"),
				Agreement:             nifcloud.Bool(true),
				UserData: &types.RequestUserData{
					Content:  nifcloud.String(base64.StdEncoding.EncodeToString([]byte("test_user_data"))),
					Encoding: nifcloud.String("base64"),
				},
				License: []types.RequestLicense{
					{
						LicenseName: types.LicenseNameOfLicenseForRunInstances("test_license_name"),
						LicenseNum:  nifcloud.String("200"),
					},
				},
				NetworkInterface: []types.RequestNetworkInterface{{
					NetworkId:   nifcloud.String("test_network_id"),
					NetworkName: nifcloud.String("test_network_name"),
					IpAddress:   nifcloud.String("test_ip_address"),
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRunInstancesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeInstancesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id": "test_instance_id",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeInstancesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeInstancesInput{
				InstanceId: []string{"test_instance_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeInstancesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeInstanceAttributeInputWithDisableAPITermination(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id": "test_instance_id",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeInstanceAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeInstanceAttributeInput{
				InstanceId: nifcloud.String("test_instance_id"),
				Attribute:  types.AttributeOfDescribeInstanceAttributeRequestDisableApiTermination,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeInstanceAttributeInputWithDisableAPITermination(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandStopInstancesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id": "test_instance_id",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.StopInstancesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.StopInstancesInput{
				InstanceId: []string{"test_instance_id"},
				Force:      nifcloud.Bool(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandStopInstancesInput(tt.args, true)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandTerminateInstancesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id": "test_instance_id",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.TerminateInstancesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.TerminateInstancesInput{
				InstanceId: []string{"test_instance_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandTerminateInstancesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyInstanceAttributeInputForAccountingType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":     "test_instance_id",
		"accounting_type": "test_accounting_type",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyInstanceAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyInstanceAttributeInput{
				InstanceId: nifcloud.String("test_instance_id"),
				Attribute:  types.AttributeOfModifyInstanceAttributeRequestAccountingType,
				Value:      nifcloud.String("test_accounting_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyInstanceAttributeInputForAccountingType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyInstanceAttributeInputForDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id": "test_instance_id",
		"description": "test_description",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyInstanceAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyInstanceAttributeInput{
				InstanceId: nifcloud.String("test_instance_id"),
				Attribute:  types.AttributeOfModifyInstanceAttributeRequestDescription,
				Value:      nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyInstanceAttributeInputForDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyInstanceAttributeInputForDisableAPITermination(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":             "test_instance_id",
		"disable_api_termination": false,
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyInstanceAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyInstanceAttributeInput{
				InstanceId: nifcloud.String("test_instance_id"),
				Attribute:  types.AttributeOfModifyInstanceAttributeRequestDisableApiTermination,
				Value:      nifcloud.String("false"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyInstanceAttributeInputForDisableAPITermination(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyInstanceAttributeInputForInstanceID(t *testing.T) {
	r := New()
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id": "test_instance_id",
	})
	rd.SetId("test_instance_id")
	dn := r.Data(rd.State())

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyInstanceAttributeInput
	}{
		{
			name: "expands the resource data",
			args: dn,
			want: &computing.ModifyInstanceAttributeInput{
				InstanceId: nifcloud.String("test_instance_id"),
				Attribute:  types.AttributeOfModifyInstanceAttributeRequestInstanceName,
				Value:      nifcloud.String("test_instance_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyInstanceAttributeInputForInstanceID(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyInstanceAttributeInputForInstanceType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":   "test_instance_id",
		"instance_type": "test_instance_type",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyInstanceAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyInstanceAttributeInput{
				InstanceId: nifcloud.String("test_instance_id"),
				Attribute:  types.AttributeOfModifyInstanceAttributeRequestInstanceType,
				Value:      nifcloud.String("test_instance_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyInstanceAttributeInputForInstanceType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyInstanceAttributeInputForSecurityGroup(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":    "test_instance_id",
		"security_group": "test_security_group",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyInstanceAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyInstanceAttributeInput{
				InstanceId: nifcloud.String("test_instance_id"),
				Attribute:  types.AttributeOfModifyInstanceAttributeRequestGroupId,
				Value:      nifcloud.String("test_security_group"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyInstanceAttributeInputForSecurityGroup(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyUpdateInstanceNetworkInterfacesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id": "test_instance_id",
		"network_interface": []interface{}{map[string]interface{}{
			"network_id":   "test_network_id",
			"network_name": "test_network_name",
			"ip_address":   "test_ip_address",
		}},
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyUpdateInstanceNetworkInterfacesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyUpdateInstanceNetworkInterfacesInput{
				InstanceId: nifcloud.String("test_instance_id"),
				NetworkInterface: []types.RequestNetworkInterface{{
					NetworkId:   nifcloud.String("test_network_id"),
					NetworkName: nifcloud.String("test_network_name"),
					IpAddress:   nifcloud.String("test_ip_address"),
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyUpdateInstanceNetworkInterfacesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandAttachNetworkInterfaceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id": "test_instance_id",
	})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.AttachNetworkInterfaceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.AttachNetworkInterfaceInput{
				InstanceId:         nifcloud.String("test_instance_id"),
				NetworkInterfaceId: nifcloud.String("test_network_interface_id"),
				NiftyReboot:        types.NiftyRebootOfAttachNetworkInterfaceRequestForce,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandAttachNetworkInterfaceInput(tt.args, "test_network_interface_id")
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDetachNetworkInterfaceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_instance_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DetachNetworkInterfaceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DetachNetworkInterfaceInput{
				AttachmentId: nifcloud.String("test_network_interface_attachment_id"),
				NiftyReboot:  types.NiftyRebootOfDetachNetworkInterfaceRequestForce,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDetachNetworkInterfaceInput(tt.args, "test_network_interface_attachment_id")
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeregisterInstancesFromSecurityGroupInput(t *testing.T) {
	r := New()
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":    "test_instance_id",
		"security_group": "test_security_group",
	})
	rd.SetId("test_instance_id")

	dn := r.Data(rd.State())

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeregisterInstancesFromSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: dn,
			want: &computing.DeregisterInstancesFromSecurityGroupInput{
				InstanceId: []string{"test_instance_id"},
				GroupName:  nifcloud.String("test_security_group"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeregisterInstancesFromSecurityGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandAssociateMultiIpAddressGroupInput(t *testing.T) {
	type args struct {
		multiIPAddressGroupID string
		instanceUniqueID      string
	}

	tests := []struct {
		name string
		args args
		want *computing.AssociateMultiIpAddressGroupInput
	}{
		{
			name: "expands the args",
			args: args{
				multiIPAddressGroupID: "test_multi_ip_address_group_id",
				instanceUniqueID:      "test_instance_unique_id",
			},
			want: &computing.AssociateMultiIpAddressGroupInput{
				MultiIpAddressGroupId: nifcloud.String("test_multi_ip_address_group_id"),
				InstanceUniqueId:      nifcloud.String("test_instance_unique_id"),
				NiftyReboot:           types.NiftyRebootOfAssociateMultiIpAddressGroupRequestFalse,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandAssociateMultiIpAddressGroupInput(tt.args.multiIPAddressGroupID, tt.args.instanceUniqueID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDisassociateMultiIpAddressGroupInput(t *testing.T) {
	type args struct {
		multiIPAddressGroupID string
		instanceUniqueID      string
	}

	tests := []struct {
		name string
		args args
		want *computing.DisassociateMultiIpAddressGroupInput
	}{
		{
			name: "expands the args",
			args: args{
				multiIPAddressGroupID: "test_multi_ip_address_group_id",
				instanceUniqueID:      "test_instance_unique_id",
			},
			want: &computing.DisassociateMultiIpAddressGroupInput{
				MultiIpAddressGroupId: nifcloud.String("test_multi_ip_address_group_id"),
				InstanceUniqueId:      nifcloud.String("test_instance_unique_id"),
				NiftyReboot:           types.NiftyRebootOfDisassociateMultiIpAddressGroupRequestFalse,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDisassociateMultiIpAddressGroupInput(tt.args.multiIPAddressGroupID, tt.args.instanceUniqueID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandStartInstancesInputWithMultiIPAddressConfigurationUserData(t *testing.T) {
	r := New()
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id": "test_instance_id",
		"multi_ip_address_configuration_user_data": "test_user_data",
	})
	rd.SetId("test_instance_id")

	dn := r.Data(rd.State())

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.StartInstancesInput
	}{
		{
			name: "expands the resource data",
			args: dn,
			want: &computing.StartInstancesInput{
				InstanceId: []string{"test_instance_id"},
				UserData: &types.RequestUserData{
					Content:  nifcloud.String(base64.StdEncoding.EncodeToString([]byte("test_user_data"))),
					Encoding: nifcloud.String("base64"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandStartInstancesInputWithMultiIPAddressConfigurationUserData(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
