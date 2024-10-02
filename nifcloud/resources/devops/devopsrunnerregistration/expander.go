package devopsrunnerregistration

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
)

func expandGetRunnerInput(d *schema.ResourceData) *devopsrunner.GetRunnerInput {
	return &devopsrunner.GetRunnerInput{
		RunnerName: nifcloud.String(d.Get("runner_name").(string)),
	}
}

func expandListRunnerRegistrationsInput(d *schema.ResourceData) *devopsrunner.ListRunnerRegistrationsInput {
	return &devopsrunner.ListRunnerRegistrationsInput{
		RunnerName: nifcloud.String(d.Get("runner_name").(string)),
	}
}

func expandRegisterRunnerInput(d *schema.ResourceData) *devopsrunner.RegisterRunnerInput {
	return &devopsrunner.RegisterRunnerInput{
		RunnerName:          nifcloud.String(d.Get("runner_name").(string)),
		GitlabUrl:           nifcloud.String(d.Get("gitlab_url").(string)),
		ParameterGroupName:  nifcloud.String(d.Get("parameter_group_name").(string)),
		AuthenticationToken: nifcloud.String(d.Get("token").(string)),
	}
}

func expandUpdateRunnerRegistrationInput(d *schema.ResourceData) *devopsrunner.UpdateRunnerRegistrationInput {
	return &devopsrunner.UpdateRunnerRegistrationInput{
		RunnerName:         nifcloud.String(d.Get("runner_name").(string)),
		RegistrationId:     nifcloud.String(d.Id()),
		ParameterGroupName: nifcloud.String(d.Get("parameter_group_name").(string)),
	}
}

func expandUnregisterRunnerInput(d *schema.ResourceData) *devopsrunner.UnregisterRunnerInput {
	return &devopsrunner.UnregisterRunnerInput{
		RunnerName:     nifcloud.String(d.Get("runner_name").(string)),
		RegistrationId: nifcloud.String(d.Id()),
		// DisableTokenRevocation is always set to true, enabling the recreation of registrations.
		DisableTokenRevocation: nifcloud.Bool(true),
	}
}
