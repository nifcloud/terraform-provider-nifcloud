package nasinstance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateNASInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"allocated_storage":              100,
		"availability_zone":              "test_zone",
		"private_ip_address":             "192.168.0.1",
		"private_ip_address_subnet_mask": "/24",
		"master_user_password":           "test_master_user_password",
		"master_username":                "test_master_username",
		"description":                    "test_description",
		"identifier":                     "test_identifier",
		"type":                           0,
		"nas_security_group_name":        "test_group_name",
		"network_id":                     "test_network_id",
		"protocol":                       "test_protocol",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.CreateNASInstanceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &nas.CreateNASInstanceInput{
				AllocatedStorage:       nifcloud.Int32(100),
				AvailabilityZone:       nifcloud.String("test_zone"),
				MasterPrivateAddress:   nifcloud.String("192.168.0.1/24"),
				MasterUserPassword:     nifcloud.String("test_master_user_password"),
				MasterUsername:         nifcloud.String("test_master_username"),
				NASInstanceDescription: nifcloud.String("test_description"),
				NASInstanceIdentifier:  nifcloud.String("test_identifier"),
				NASInstanceType:        nifcloud.Int32(0),
				NASSecurityGroups:      []string{"test_group_name"},
				NetworkId:              nifcloud.String("test_network_id"),
				Protocol:               types.ProtocolOfCreateNASInstanceRequest("test_protocol"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateNASInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeNASInstancesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.DescribeNASInstancesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &nas.DescribeNASInstancesInput{
				NASInstanceIdentifier: nifcloud.String("test_identifier"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeNASInstancesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyNASInstanceInput(t *testing.T) {
	rdForNFS := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"protocol":                       "nfs",
		"allocated_storage":              100,
		"private_ip_address":             "192.168.0.1",
		"private_ip_address_subnet_mask": "/24",
		"description":                    "test_description",
		"identifier":                     "test_new_identifier",
		"nas_security_group_name":        "test_group_name",
		"network_id":                     "test_network_id",
		"no_root_squash":                 true,
	})
	rdForNFS.SetId("test_identifier")

	rdForCIFS := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"protocol":                                 "cifs",
		"allocated_storage":                        100,
		"authentication_type":                      1,
		"directory_service_administrator_name":     "test_directory_service_administrator_name",
		"directory_service_administrator_password": "test_directory_service_administrator_password",
		"directory_service_domain_name":            "test_directory_service_domain_name",
		"domain_controllers": []interface{}{map[string]interface{}{
			"hostname":   "test_hostname",
			"ip_address": "test_ip_address",
		}},
		"private_ip_address":             "192.168.0.1",
		"private_ip_address_subnet_mask": "/24",
		"master_user_password":           "test_master_user_password",
		"description":                    "test_description",
		"identifier":                     "test_new_identifier",
		"nas_security_group_name":        "test_group_name",
		"network_id":                     "test_network_id",
		"no_root_squash":                 true,
	})
	rdForCIFS.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.ModifyNASInstanceInput
	}{
		{
			name: "expands the resource data for nfs",
			args: rdForNFS,
			want: &nas.ModifyNASInstanceInput{
				AllocatedStorage:         nifcloud.Int32(100),
				MasterPrivateAddress:     nifcloud.String("192.168.0.1/24"),
				NASInstanceDescription:   nifcloud.String("test_description"),
				NASInstanceIdentifier:    nifcloud.String("test_identifier"),
				NASSecurityGroups:        []string{"test_group_name"},
				NetworkId:                nifcloud.String("test_network_id"),
				NewNASInstanceIdentifier: nifcloud.String("test_new_identifier"),
				NoRootSquash:             nifcloud.Bool(true),
			},
		},
		{
			name: "expands the resource data for cifs",
			args: rdForCIFS,
			want: &nas.ModifyNASInstanceInput{
				AllocatedStorage:                      nifcloud.Int32(100),
				AuthenticationType:                    nifcloud.Int32(1),
				DirectoryServiceAdministratorName:     nifcloud.String("test_directory_service_administrator_name"),
				DirectoryServiceAdministratorPassword: nifcloud.String("test_directory_service_administrator_password"),
				DirectoryServiceDomainName:            nifcloud.String("test_directory_service_domain_name"),
				DomainControllers: []types.RequestDomainControllers{
					{
						Hostname:  nifcloud.String("test_hostname"),
						IPAddress: nifcloud.String("test_ip_address"),
					},
				},
				MasterPrivateAddress:     nifcloud.String("192.168.0.1/24"),
				MasterUserPassword:       nifcloud.String("test_master_user_password"),
				NASInstanceDescription:   nifcloud.String("test_description"),
				NASInstanceIdentifier:    nifcloud.String("test_identifier"),
				NASSecurityGroups:        []string{"test_group_name"},
				NetworkId:                nifcloud.String("test_network_id"),
				NewNASInstanceIdentifier: nifcloud.String("test_new_identifier"),
				NoRootSquash:             nifcloud.Bool(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyNASInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteNASInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"directory_service_administrator_name":     "test_directory_service_administrator_name",
		"directory_service_administrator_password": "test_directory_service_administrator_password",
	})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.DeleteNASInstanceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &nas.DeleteNASInstanceInput{
				NASInstanceIdentifier:                 nifcloud.String("test_identifier"),
				DirectoryServiceAdministratorName:     nifcloud.String("test_directory_service_administrator_name"),
				DirectoryServiceAdministratorPassword: nifcloud.String("test_directory_service_administrator_password"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteNASInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
