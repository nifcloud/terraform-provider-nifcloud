package devopsrunnerparametergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func updateRunnerParameterGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	// lintignore:R019
	if d.HasChanges(
		"docker_disable_cache",
		"docker_disable_entrypoint_overwrite",
		"docker_extra_hosts",
		"docker_image",
		"docker_oom_kill_disable",
		"docker_privileged",
		"docker_shm_size",
		"docker_tls_verify",
		"docker_volumes",
	) {
		input := expandUpdateRunnerParameterInput(d)

		if _, err := svc.UpdateRunnerParameter(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update parameters of the DevOps Runner parameter group: %s", err))
		}
	}

	if d.HasChanges("name", "description") {
		input := expandUpdateRunnerParameterGroupInput(d)

		if _, err := svc.UpdateRunnerParameterGroup(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update a DevOps Runner parameter group: %s", err))
		}
	}

	if d.HasChange("name") {
		d.SetId(d.Get("name").(string))
	}

	return readRunnerParameterGroup(ctx, d, meta)
}
