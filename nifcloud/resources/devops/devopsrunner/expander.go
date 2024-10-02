package devopsrunner

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner/types"
)

func expandCreateRunnerInput(d *schema.ResourceData) *devopsrunner.CreateRunnerInput {
	var concurrent *int32
	if v := d.Get("concurrent").(int); v != 0 {
		concurrent = nifcloud.Int32(int32(v))
	}

	return &devopsrunner.CreateRunnerInput{
		RunnerName:       nifcloud.String(d.Get("name").(string)),
		InstanceType:     types.InstanceTypeOfCreateRunnerRequest(d.Get("instance_type").(string)),
		AvailabilityZone: types.AvailabilityZoneOfCreateRunnerRequest(d.Get("availability_zone").(string)),
		Concurrent:       concurrent,
		Description:      nifcloud.String(d.Get("description").(string)),
		NetworkConfig:    expandNetworkConfig(d),
	}
}

func expandUpdateRunnerInput(d *schema.ResourceData) *devopsrunner.UpdateRunnerInput {
	return &devopsrunner.UpdateRunnerInput{
		RunnerName:        nifcloud.String(d.Id()),
		ChangedRunnerName: nifcloud.String(d.Get("name").(string)),
		Concurrent:        nifcloud.Int32(int32(d.Get("concurrent").(int))),
		Description:       nifcloud.String(d.Get("description").(string)),
	}
}

func expandGetRunnerInput(d *schema.ResourceData) *devopsrunner.GetRunnerInput {
	return &devopsrunner.GetRunnerInput{
		RunnerName: nifcloud.String(d.Id()),
	}
}

func expandDeleteRunnerInput(d *schema.ResourceData) *devopsrunner.DeleteRunnerInput {
	return &devopsrunner.DeleteRunnerInput{
		RunnerName: nifcloud.String(d.Id()),
	}
}

func expandModifyRunnerInstanceTypeInput(d *schema.ResourceData) *devopsrunner.ModifyRunnerInstanceTypeInput {
	return &devopsrunner.ModifyRunnerInstanceTypeInput{
		RunnerName:   nifcloud.String(d.Id()),
		InstanceType: types.InstanceTypeOfModifyRunnerInstanceTypeRequest(d.Get("instance_type").(string)),
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
