package devopsrunnerparametergroup

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner/types"
)

func expandCreateRunnerParameterGroupInput(d *schema.ResourceData) *devopsrunner.CreateRunnerParameterGroupInput {
	return &devopsrunner.CreateRunnerParameterGroupInput{
		ParameterGroupName: nifcloud.String(d.Get("name").(string)),
		Description:        nifcloud.String(d.Get("description").(string)),
	}
}

func expandUpdateRunnerParameterGroupInput(d *schema.ResourceData) *devopsrunner.UpdateRunnerParameterGroupInput {
	return &devopsrunner.UpdateRunnerParameterGroupInput{
		ParameterGroupName:        nifcloud.String(d.Id()),
		ChangedParameterGroupName: nifcloud.String(d.Get("name").(string)),
		Description:               nifcloud.String(d.Get("description").(string)),
	}
}

func expandGetRunnerParameterGroupInput(d *schema.ResourceData) *devopsrunner.GetRunnerParameterGroupInput {
	return &devopsrunner.GetRunnerParameterGroupInput{
		ParameterGroupName: nifcloud.String(d.Id()),
	}
}

func expandDeleteRunnerParameterGroupInput(d *schema.ResourceData) *devopsrunner.DeleteRunnerParameterGroupInput {
	return &devopsrunner.DeleteRunnerParameterGroupInput{
		ParameterGroupName: nifcloud.String(d.Id()),
	}
}

func expandUpdateRunnerParameterInput(d *schema.ResourceData) *devopsrunner.UpdateRunnerParameterInput {
	return &devopsrunner.UpdateRunnerParameterInput{
		ParameterGroupName: nifcloud.String(d.Id()),
		DockerParameters: &types.RequestDockerParameters{
			DisableCache:               nifcloud.Bool(d.Get("docker_disable_cache").(bool)),
			DisableEntrypointOverwrite: nifcloud.Bool(d.Get("docker_disable_entrypoint_overwrite").(bool)),
			ListOfRequestExtraHosts:    expandRequestExtraHosts(d),
			Image:                      nifcloud.String(d.Get("docker_image").(string)),
			OomKillDisable:             nifcloud.Bool(d.Get("docker_oom_kill_disable").(bool)),
			Privileged:                 nifcloud.Bool(d.Get("docker_privileged").(bool)),
			ShmSize:                    nifcloud.Int32(int32(d.Get("docker_shm_size").(int))),
			TlsVerify:                  nifcloud.Bool(d.Get("docker_tls_verify").(bool)),
			ListOfRequestVolumes:       expandRequestVolumes(d),
		},
	}
}

func expandRequestExtraHosts(d *schema.ResourceData) []types.RequestExtraHosts {
	configured := d.Get("docker_extra_host").(*schema.Set)
	ret := make([]types.RequestExtraHosts, configured.Len())
	for i, raw := range configured.List() {
		v := raw.(map[string]interface{})

		ret[i] = types.RequestExtraHosts{
			HostName:  nifcloud.String(v["host_name"].(string)),
			IpAddress: nifcloud.String(v["ip_address"].(string)),
		}
	}
	return ret
}

func expandRequestVolumes(d *schema.ResourceData) []string {
	configured := d.Get("docker_volume").(*schema.Set)
	ret := make([]string, configured.Len())
	for i, raw := range configured.List() {
		ret[i] = raw.(string)
	}
	return ret
}
