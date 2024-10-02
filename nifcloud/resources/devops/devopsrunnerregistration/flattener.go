package devopsrunnerregistration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
)

func flattenRunnerName(d *schema.ResourceData, res *devopsrunner.GetRunnerOutput) error {
	if res == nil || res.Runner == nil {
		d.SetId("")
		return nil
	}

	runner := res.Runner

	if nifcloud.ToString(runner.RunnerName) != d.Get("runner_name").(string) {
		return fmt.Errorf("unable to find the DevOps Runner within: %#v", runner)
	}

	if err := d.Set("runner_name", runner.RunnerName); err != nil {
		return err
	}

	return nil
}

func flatten(d *schema.ResourceData, res *devopsrunner.ListRunnerRegistrationsOutput) error {
	if res == nil || res.Registrations == nil {
		return nil
	}

	for _, r := range res.Registrations {
		if nifcloud.ToString(r.RegistrationId) != d.Id() {
			continue
		}

		if err := d.Set("gitlab_url", r.GitlabUrl); err != nil {
			return err
		}

		if err := d.Set("parameter_group_name", r.ParameterGroupName); err != nil {
			return err
		}

		if err := d.Set("token", r.Token); err != nil {
			return err
		}
	}

	return nil
}
