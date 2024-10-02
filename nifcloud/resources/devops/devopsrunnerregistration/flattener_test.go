package devopsrunnerregistration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner/types"
	"github.com/stretchr/testify/assert"
)

func TestFlattenRunnerName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"runner_name": "test_name",
	})
	rd.SetId("test_id")

	wantRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"runner_name": "test_name",
	})
	wantRd.SetId("test_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *devopsrunner.GetRunnerOutput
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
				res: &devopsrunner.GetRunnerOutput{
					Runner: &types.Runner{
						RunnerName: nifcloud.String("test_name"),
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
			err := flattenRunnerName(tt.args.d, tt.args.res)
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

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"runner_name": "test_name",
	})
	rd.SetId("test_id")

	wantRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"runner_name":          "test_name",
		"gitlab_url":           "test_url",
		"parameter_group_name": "test_name_pg",
		"token":                "test_token",
	})
	wantRd.SetId("test_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *devopsrunner.ListRunnerRegistrationsOutput
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
				res: &devopsrunner.ListRunnerRegistrationsOutput{
					Registrations: []types.Registrations{{
						RegistrationId:     nifcloud.String("test_id"),
						GitlabUrl:          nifcloud.String("test_url/"),
						ParameterGroupName: nifcloud.String("test_name_pg"),
						Token:              nifcloud.String("test_token"),
					}},
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
