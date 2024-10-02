package devopsrunnerregistration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/stretchr/testify/assert"
)

func TestExpandGetsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"runner_name":          "test_name",
		"gitlab_url":           "test_url",
		"parameter_group_name": "test_name_pg",
		"token":                "glrt-test_token",
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.GetRunnerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.GetRunnerInput{
				RunnerName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetRunnerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandListRunnerRegistrationsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"runner_name":          "test_name",
		"gitlab_url":           "test_url",
		"parameter_group_name": "test_name_pg",
		"token":                "glrt-test_token",
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.ListRunnerRegistrationsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.ListRunnerRegistrationsInput{
				RunnerName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandListRunnerRegistrationsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRegisterRunnerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"runner_name":          "test_name",
		"gitlab_url":           "test_url",
		"parameter_group_name": "test_name_pg",
		"token":                "glrt-test_token",
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.RegisterRunnerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.RegisterRunnerInput{
				RunnerName:          nifcloud.String("test_name"),
				GitlabUrl:           nifcloud.String("test_url"),
				ParameterGroupName:  nifcloud.String("test_name_pg"),
				AuthenticationToken: nifcloud.String("glrt-test_token"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRegisterRunnerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"runner_name":          "test_name",
		"gitlab_url":           "test_url",
		"parameter_group_name": "test_name_pg",
		"token":                "glrt-test_token",
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.UpdateRunnerRegistrationInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.UpdateRunnerRegistrationInput{
				RunnerName:         nifcloud.String("test_name"),
				RegistrationId:     nifcloud.String("test_id"),
				ParameterGroupName: nifcloud.String("test_name_pg"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateRunnerRegistrationInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUnregisterRunnerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"runner_name":          "test_name",
		"gitlab_url":           "test_url",
		"parameter_group_name": "test_name_pg",
		"token":                "glrt-test_token",
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.UnregisterRunnerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.UnregisterRunnerInput{
				RunnerName:             nifcloud.String("test_name"),
				RegistrationId:         nifcloud.String("test_id"),
				DisableTokenRevocation: nifcloud.Bool(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUnregisterRunnerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
