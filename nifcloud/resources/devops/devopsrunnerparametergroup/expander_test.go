package devopsrunnerparametergroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateRunnerParameterGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"description": "test_description",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.CreateRunnerParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.CreateRunnerParameterGroupInput{
				ParameterGroupName: nifcloud.String("test_name"),
				Description:        nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateRunnerParameterGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateParameterGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name_changed",
		"description": "test_description",
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.UpdateRunnerParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.UpdateRunnerParameterGroupInput{
				ParameterGroupName:        nifcloud.String("test_name"),
				ChangedParameterGroupName: nifcloud.String("test_name_changed"),
				Description:               nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateRunnerParameterGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetParameterGroupsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.GetRunnerParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.GetRunnerParameterGroupInput{
				ParameterGroupName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetRunnerParameterGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteParameterGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.DeleteRunnerParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.DeleteRunnerParameterGroupInput{
				ParameterGroupName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteRunnerParameterGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateRunnerParameterInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":                                "test_name",
		"description":                         "test_description",
		"docker_disable_cache":                true,
		"docker_disable_entrypoint_overwrite": true,
		"docker_extra_host": []interface{}{
			map[string]interface{}{
				"host_name":  "test_host_name",
				"ip_address": "test_address",
			},
		},
		"docker_image":            "test_image",
		"docker_oom_kill_disable": true,
		"docker_privileged":       true,
		"docker_shm_size":         1,
		"docker_tls_verify":       true,
		"docker_volume":           []interface{}{"test_volume"},
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.UpdateRunnerParameterInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.UpdateRunnerParameterInput{
				ParameterGroupName: nifcloud.String("test_name"),
				DockerParameters: &types.RequestDockerParameters{
					DisableCache:               nifcloud.Bool(true),
					DisableEntrypointOverwrite: nifcloud.Bool(true),
					ListOfRequestExtraHosts: []types.RequestExtraHosts{
						{
							HostName:  nifcloud.String("test_host_name"),
							IpAddress: nifcloud.String("test_address"),
						},
					},
					Image:                nifcloud.String("test_image"),
					OomKillDisable:       nifcloud.Bool(true),
					Privileged:           nifcloud.Bool(true),
					ShmSize:              nifcloud.Int32(int32(1)),
					TlsVerify:            nifcloud.Bool(true),
					ListOfRequestVolumes: []string{"test_volume"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateRunnerParameterInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
