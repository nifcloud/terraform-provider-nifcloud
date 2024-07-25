package devopsrunner

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
)

func flatten(d *schema.ResourceData, res *devopsrunner.GetRunnerOutput) error {
	if res == nil || res.Runner == nil {
		d.SetId("")
		return nil
	}

	runner := res.Runner

	if nifcloud.ToString(runner.RunnerName) != d.Id() {
		return fmt.Errorf("unable to find the DevOps Runner within: %#v", runner)
	}

	if err := d.Set("name", runner.RunnerName); err != nil {
		return err
	}

	if err := d.Set("instance_type", runner.InstanceType); err != nil {
		return err
	}

	if err := d.Set("availability_zone", runner.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("concurrent", runner.Concurrent); err != nil {
		return err
	}

	if err := d.Set("description", runner.Description); err != nil {
		return err
	}

	if err := d.Set("network_id", runner.NetworkConfig.NetworkId); err != nil {
		return err
	}

	if err := d.Set("private_address", runner.NetworkConfig.PrivateAddress); err != nil {
		return err
	}

	if err := d.Set("public_ip_address", runner.PublicIpAddress); err != nil {
		return err
	}

	if err := d.Set("system_id", runner.SystemId); err != nil {
		return err
	}

	return nil
}
