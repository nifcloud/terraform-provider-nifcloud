package devopsinstance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
)

func flatten(d *schema.ResourceData, res *devops.GetInstanceOutput) error {
	if res == nil || res.Instance == nil {
		d.SetId("")
		return nil
	}

	instance := res.Instance

	if nifcloud.ToString(instance.InstanceId) != d.Id() {
		return fmt.Errorf("unable to find the DevOps instance within: %#v", instance)
	}

	if err := d.Set("instance_id", instance.InstanceId); err != nil {
		return err
	}

	if err := d.Set("instance_type", instance.InstanceType); err != nil {
		return err
	}

	if err := d.Set("firewall_group_name", instance.FirewallGroupName); err != nil {
		return err
	}

	if err := d.Set("parameter_group_name", instance.ParameterGroupName); err != nil {
		return err
	}

	if err := d.Set("disk_size", instance.DiskSize); err != nil {
		return err
	}

	if err := d.Set("availability_zone", instance.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("description", instance.Description); err != nil {
		return err
	}

	if err := d.Set("network_id", instance.NetworkConfig.NetworkId); err != nil {
		return err
	}

	if err := d.Set("private_address", instance.NetworkConfig.PrivateAddress); err != nil {
		return err
	}

	if err := d.Set("object_storage_account", instance.ObjectStorageConfig.Account); err != nil {
		return err
	}

	if err := d.Set("object_storage_region", instance.ObjectStorageConfig.Region); err != nil {
		return err
	}

	if instance.ObjectStorageConfig.BucketUseObjects != nil {
		if err := d.Set("lfs_bucket_name", instance.ObjectStorageConfig.BucketUseObjects.Lfs); err != nil {
			return err
		}

		if err := d.Set("packages_bucket_name", instance.ObjectStorageConfig.BucketUseObjects.Packages); err != nil {
			return err
		}

		if err := d.Set("container_registry_bucket_name", instance.ObjectStorageConfig.BucketUseObjects.ContainerRegistry); err != nil {
			return err
		}
	}

	if err := d.Set("to", instance.To); err != nil {
		return err
	}

	if err := d.Set("gitlab_url", instance.GitlabUrl); err != nil {
		return err
	}

	if err := d.Set("registry_url", instance.RegistryUrl); err != nil {
		return err
	}

	if err := d.Set("public_ip_address", instance.PublicIpAddress); err != nil {
		return err
	}

	return nil
}
