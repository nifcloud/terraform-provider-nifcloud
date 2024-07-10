package devopsrunnerparametergroup

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner/types"
)

func flatten(d *schema.ResourceData, res *devopsrunner.GetRunnerParameterGroupOutput) error {
	if res == nil || res.ParameterGroup == nil {
		d.SetId("")
		return nil
	}

	group := res.ParameterGroup

	if nifcloud.ToString(group.ParameterGroupName) != d.Id() {
		return fmt.Errorf("unable to find the DevOps Runner parameter group within: %#v", group)
	}

	if err := d.Set("name", group.ParameterGroupName); err != nil {
		return err
	}

	if err := d.Set("description", group.Description); err != nil {
		return err
	}

	if err := d.Set("docker_disable_cache", group.DockerParameters.DisableCache); err != nil {
		return err
	}

	if err := d.Set("docker_disable_entrypoint_overwrite", group.DockerParameters.DisableEntrypointOverwrite); err != nil {
		return err
	}

	if err := d.Set("docker_extra_host", flattenDockerExtraHosts(group.DockerParameters.ExtraHosts)); err != nil {
		return err
	}

	if err := d.Set("docker_image", group.DockerParameters.Image); err != nil {
		return err
	}

	if err := d.Set("docker_oom_kill_disable", group.DockerParameters.OomKillDisable); err != nil {
		return err
	}

	if err := d.Set("docker_privileged", group.DockerParameters.Privileged); err != nil {
		return err
	}

	if err := d.Set("docker_shm_size", group.DockerParameters.ShmSize); err != nil {
		return err
	}

	if err := d.Set("docker_tls_verify", group.DockerParameters.TlsVerify); err != nil {
		return err
	}

	if err := d.Set("docker_volume", group.DockerParameters.Volumes); err != nil {
		return err
	}

	return nil
}

func flattenDockerExtraHosts(extraHosts []types.ExtraHosts) []map[string]string {
	ret := make([]map[string]string, len(extraHosts))
	for i, eh := range extraHosts {
		ret[i] = map[string]string{
			"host_name":  *eh.HostName,
			"ip_address": *eh.IpAddress,
		}
	}
	return ret
}
