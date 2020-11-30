package instance

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
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
				InstanceType:  computing.InstanceTypeOfRunInstancesRequest("test_instance_type"),
				Placement: &computing.RequestPlacement{
					AvailabilityZone: nifcloud.String(("test_availability_zone")),
				},
				DisableApiTermination: nifcloud.Bool(false),
				AccountingType:        computing.AccountingTypeOfRunInstancesRequest("test_accounting_type"),
				Description:           nifcloud.String("test_description"),
				Admin:                 nifcloud.String("test_admin"),
				Password:              nifcloud.String("test_password"),
				Agreement:             nifcloud.Bool(true),
				UserData: &computing.RequestUserDataOfRunInstances{
					Content:  nifcloud.String(base64.StdEncoding.EncodeToString([]byte("test_user_data"))),
					Encoding: nifcloud.String("base64"),
				},
				License: []computing.RequestLicense{
					{
						LicenseName: computing.LicenseNameOfLicenseForRunInstances("test_license_name"),
						LicenseNum:  nifcloud.String("200"),
					},
				},
				NetworkInterface: []computing.RequestNetworkInterface{{
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
				Attribute:  computing.AttributeOfDescribeInstanceAttributeRequestDisableApiTermination,
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandStopInstancesInput(tt.args)
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
				Attribute:  "accountingType",
				Value:      "test_accounting_type",
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
				Attribute:  "description",
				Value:      "test_description",
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
				Attribute:  "disableApiTermination",
				Value:      computing.ValueOfModifyInstanceAttributeRequest("false"),
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
				Attribute:  "instanceName",
				Value:      computing.ValueOfModifyInstanceAttributeRequest("test_instance_id"),
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
				Attribute:  "instanceType",
				Value:      computing.ValueOfModifyInstanceAttributeRequest("test_instance_type"),
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
				Attribute:  "groupId",
				Value:      computing.ValueOfModifyInstanceAttributeRequest("test_security_group"),
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
				NetworkInterface: []computing.RequestNetworkInterfaceOfNiftyUpdateInstanceNetworkInterfaces{{
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
