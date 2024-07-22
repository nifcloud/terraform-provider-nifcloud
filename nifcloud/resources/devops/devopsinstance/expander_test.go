package devopsinstance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":                    "test_id",
		"instance_type":                  "c-large",
		"firewall_group_name":            "test_name_fg",
		"parameter_group_name":           "test_name_pg",
		"disk_size":                      100,
		"availability_zone":              "east-11",
		"description":                    "test_description",
		"initial_root_password":          "test_password",
		"network_id":                     "test_id_nw",
		"private_address":                "192.168.1.1/24",
		"object_storage_account":         "test_account",
		"object_storage_region":          "test_region",
		"lfs_bucket_name":                "test_name_lfs",
		"packages_bucket_name":           "test_name_pkg",
		"container_registry_bucket_name": "test_name_cr",
		"to":                             "test@mail.com",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.CreateInstanceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.CreateInstanceInput{
				InstanceId:          nifcloud.String("test_id"),
				InstanceType:        types.InstanceTypeOfCreateInstanceRequest("c-large"),
				FirewallGroupName:   nifcloud.String("test_name_fg"),
				ParameterGroupName:  nifcloud.String("test_name_pg"),
				DiskSize:            nifcloud.Int32(int32(100)),
				AvailabilityZone:    types.AvailabilityZoneOfCreateInstanceRequest("east-11"),
				Description:         nifcloud.String("test_description"),
				InitialRootPassword: nifcloud.String("test_password"),
				NetworkConfig: &types.RequestNetworkConfig{
					NetworkId:      nifcloud.String("test_id_nw"),
					PrivateAddress: nifcloud.String("192.168.1.1/24"),
				},
				ObjectStorageConfig: &types.RequestObjectStorageConfig{
					Account: nifcloud.String("test_account"),
					Region:  types.RegionOfobjectStorageConfigForCreateInstance("test_region"),
					RequestBucketUseObjects: &types.RequestBucketUseObjects{
						Lfs:               nifcloud.String("test_name_lfs"),
						Packages:          nifcloud.String("test_name_pkg"),
						ContainerRegistry: nifcloud.String("test_name_cr"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":                    "test_id",
		"instance_type":                  "c-large",
		"firewall_group_name":            "test_name_fg",
		"parameter_group_name":           "test_name_pg",
		"disk_size":                      100,
		"availability_zone":              "east-11",
		"description":                    "test_description",
		"initial_root_password":          "test_password",
		"network_id":                     "test_id_nw",
		"private_address":                "192.168.1.1/24",
		"object_storage_account":         "test_account",
		"object_storage_region":          "test_region",
		"lfs_bucket_name":                "test_name_lfs",
		"packages_bucket_name":           "test_name_pkg",
		"container_registry_bucket_name": "test_name_cr",
		"to":                             "test@mail.com",
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.UpdateInstanceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.UpdateInstanceInput{
				InstanceId:        nifcloud.String("test_id"),
				InstanceType:      types.InstanceTypeOfUpdateInstanceRequest("c-large"),
				FirewallGroupName: nifcloud.String("test_name_fg"),
				Description:       nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.GetInstanceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.GetInstanceInput{
				InstanceId: nifcloud.String("test_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.DeleteInstanceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.DeleteInstanceInput{
				InstanceId: nifcloud.String("test_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandExtendDiskInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.ExtendDiskInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.ExtendDiskInput{
				InstanceId: nifcloud.String("test_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandExtendDiskInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandSetupAlertInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":                    "test_id",
		"instance_type":                  "c-large",
		"firewall_group_name":            "test_name_fg",
		"parameter_group_name":           "test_name_pg",
		"disk_size":                      100,
		"availability_zone":              "east-11",
		"description":                    "test_description",
		"initial_root_password":          "test_password",
		"network_id":                     "test_id_nw",
		"private_address":                "192.168.1.1/24",
		"object_storage_account":         "test_account",
		"object_storage_region":          "test_region",
		"lfs_bucket_name":                "test_name_lfs",
		"packages_bucket_name":           "test_name_pkg",
		"container_registry_bucket_name": "test_name_cr",
		"to":                             "test@mail.com",
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.SetupAlertInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.SetupAlertInput{
				InstanceId: nifcloud.String("test_id"),
				To:         nifcloud.String("test@mail.com"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandSetupAlertInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateNetworkInterfaceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":                    "test_id",
		"instance_type":                  "c-large",
		"firewall_group_name":            "test_name_fg",
		"parameter_group_name":           "test_name_pg",
		"disk_size":                      100,
		"availability_zone":              "east-11",
		"description":                    "test_description",
		"initial_root_password":          "test_password",
		"network_id":                     "test_id_nw",
		"private_address":                "192.168.1.1/24",
		"object_storage_account":         "test_account",
		"object_storage_region":          "test_region",
		"lfs_bucket_name":                "test_name_lfs",
		"packages_bucket_name":           "test_name_pkg",
		"container_registry_bucket_name": "test_name_cr",
		"to":                             "test@mail.com",
	})
	rd.SetId("test_id")

	noNetworkConfigRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":                    "test_id",
		"instance_type":                  "c-large",
		"firewall_group_name":            "test_name_fg",
		"parameter_group_name":           "test_name_pg",
		"disk_size":                      100,
		"availability_zone":              "east-11",
		"description":                    "test_description",
		"initial_root_password":          "test_password",
		"object_storage_account":         "test_account",
		"object_storage_region":          "test_region",
		"lfs_bucket_name":                "test_name_lfs",
		"packages_bucket_name":           "test_name_pkg",
		"container_registry_bucket_name": "test_name_cr",
		"to":                             "test@mail.com",
	})
	noNetworkConfigRd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.UpdateNetworkInterfaceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.UpdateNetworkInterfaceInput{
				InstanceId: nifcloud.String("test_id"),
				NetworkConfig: &types.RequestNetworkConfig{
					NetworkId:      nifcloud.String("test_id_nw"),
					PrivateAddress: nifcloud.String("192.168.1.1/24"),
				},
			},
		},
		{
			name: "expands the resource data (no network config)",
			args: noNetworkConfigRd,
			want: &devops.UpdateNetworkInterfaceInput{
				InstanceId:    nifcloud.String("test_id"),
				NetworkConfig: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateNetworkInterfaceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
