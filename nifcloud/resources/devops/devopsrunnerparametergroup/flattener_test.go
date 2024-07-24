package devopsrunnerparametergroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	wantRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
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
	wantRd.SetId("test_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *devopsrunner.GetRunnerParameterGroupOutput
		d   *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want *schema.ResourceData
	}{
		{
			name: "flattens the response",
			args: args{
				d: rd,
				res: &devopsrunner.GetRunnerParameterGroupOutput{
					ParameterGroup: &types.ParameterGroup{
						ParameterGroupName: nifcloud.String("test_name"),
						Description:        nifcloud.String("test_description"),
						DockerParameters: &types.DockerParameters{
							DisableCache:               nifcloud.Bool(true),
							DisableEntrypointOverwrite: nifcloud.Bool(true),
							ExtraHosts: []types.ExtraHosts{
								{
									HostName:  nifcloud.String("test_host_name"),
									IpAddress: nifcloud.String("test_address"),
								},
							},
							Image:          nifcloud.String("test_image"),
							OomKillDisable: nifcloud.Bool(true),
							Privileged:     nifcloud.Bool(true),
							ShmSize:        nifcloud.Int32(int32(1)),
							TlsVerify:      nifcloud.Bool(true),
							Volumes:        []string{"test_volume"},
						},
					},
				},
			},
			want: wantRd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d:   wantNotFoundRd,
				res: nil,
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			assert.NoError(t, err)

			wantState := tt.want.State()
			if wantState == nil {
				tt.want.SetId("some")
				wantState = tt.want.State()
			}

			gotState := tt.args.d.State()
			if gotState == nil {
				tt.args.d.SetId("some")
				gotState = tt.args.d.State()
			}

			assert.Equal(t, wantState.Attributes, gotState.Attributes)
		})
	}
}
