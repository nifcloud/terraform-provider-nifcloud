package devopsinstance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
)

func expandCreateInstanceInput(d *schema.ResourceData) *devops.CreateInstanceInput {
	input := &devops.CreateInstanceInput{
		InstanceId:          nifcloud.String(d.Get("instance_id").(string)),
		InstanceType:        types.InstanceTypeOfCreateInstanceRequest(d.Get("instance_type").(string)),
		FirewallGroupName:   nifcloud.String(d.Get("firewall_group_name").(string)),
		ParameterGroupName:  nifcloud.String(d.Get("parameter_group_name").(string)),
		DiskSize:            nifcloud.Int32(int32(d.Get("disk_size").(int))),
		AvailabilityZone:    types.AvailabilityZoneOfCreateInstanceRequest(d.Get("availability_zone").(string)),
		Description:         nifcloud.String(d.Get("description").(string)),
		InitialRootPassword: nifcloud.String(d.Get("initial_root_password").(string)),
		NetworkConfig:       expandNetworkConfig(d),
		ObjectStorageConfig: expandObjectStorageConfig(d),
	}
	return input
}

func expandUpdateInstanceInput(d *schema.ResourceData) *devops.UpdateInstanceInput {
	return &devops.UpdateInstanceInput{
		InstanceId:        nifcloud.String(d.Id()),
		InstanceType:      types.InstanceTypeOfUpdateInstanceRequest(d.Get("instance_type").(string)),
		FirewallGroupName: nifcloud.String(d.Get("firewall_group_name").(string)),
		Description:       nifcloud.String(d.Get("description").(string)),
	}
}

func expandGetInstanceInput(d *schema.ResourceData) *devops.GetInstanceInput {
	return &devops.GetInstanceInput{
		InstanceId: nifcloud.String(d.Id()),
	}
}

func expandDeleteInstanceInput(d *schema.ResourceData) *devops.DeleteInstanceInput {
	return &devops.DeleteInstanceInput{
		InstanceId: nifcloud.String(d.Id()),
	}
}

func expandExtendDiskInput(d *schema.ResourceData) *devops.ExtendDiskInput {
	return &devops.ExtendDiskInput{InstanceId: nifcloud.String(d.Id())}
}

func expandSetupAlertInput(d *schema.ResourceData) *devops.SetupAlertInput {
	return &devops.SetupAlertInput{
		InstanceId: nifcloud.String(d.Id()),
		To:         nifcloud.String(d.Get("to").(string)),
	}
}

func expandUpdateNetworkInterfaceInput(d *schema.ResourceData) *devops.UpdateNetworkInterfaceInput {
	return &devops.UpdateNetworkInterfaceInput{
		InstanceId:    nifcloud.String(d.Id()),
		NetworkConfig: expandNetworkConfig(d),
	}
}

func expandNetworkConfig(d *schema.ResourceData) *types.RequestNetworkConfig {
	networkId := d.Get("network_id").(string)
	privateAddress := d.Get("private_address").(string)
	if networkId == "" && privateAddress == "" {
		return nil
	}

	return &types.RequestNetworkConfig{
		NetworkId:      nifcloud.String(d.Get("network_id").(string)),
		PrivateAddress: nifcloud.String(d.Get("private_address").(string)),
	}
}

func expandObjectStorageConfig(d *schema.ResourceData) *types.RequestObjectStorageConfig {
	if d.Get("object_storage_account").(string) == "" || d.Get("object_storage_region").(string) == "" {
		return nil
	}

	return &types.RequestObjectStorageConfig{
		Account: nifcloud.String(d.Get("object_storage_account").(string)),
		Region:  types.RegionOfobjectStorageConfigForCreateInstance(d.Get("object_storage_region").(string)),
		RequestBucketUseObjects: &types.RequestBucketUseObjects{
			Lfs:               nifcloud.String(d.Get("lfs_bucket_name").(string)),
			Packages:          nifcloud.String(d.Get("packages_bucket_name").(string)),
			ContainerRegistry: nifcloud.String(d.Get("container_registry_bucket_name").(string)),
		},
	}
}
