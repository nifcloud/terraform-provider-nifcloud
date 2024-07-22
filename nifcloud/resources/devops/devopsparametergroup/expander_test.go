package devopsparametergroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateParameterGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"description": "test_description",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.CreateParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.CreateParameterGroupInput{
				ParameterGroupName: nifcloud.String("test_name"),
				Description:        nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateParameterGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateParameterGroupInput(t *testing.T) {
	parameters := &types.RequestParameters{
		GitlabEmailFrom:    nifcloud.String("test_value"),
		GitlabEmailReplyTo: nifcloud.String("test_value"),
		SmtpPassword:       nifcloud.String("test_value"),
		SmtpUserName:       nifcloud.String("test_value"),
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name_changed",
		"description": "test_description",
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.UpdateParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.UpdateParameterGroupInput{
				ParameterGroupName:        nifcloud.String("test_name"),
				ChangedParameterGroupName: nifcloud.String("test_name_changed"),
				Description:               nifcloud.String("test_description"),
				Parameters:                parameters,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateParameterGroupInput(tt.args, parameters)
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
		want *devops.GetParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.GetParameterGroupInput{
				ParameterGroupName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetParameterGroupInput(tt.args)
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
		want *devops.DeleteParameterGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.DeleteParameterGroupInput{
				ParameterGroupName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteParameterGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateParameterGroupParameters(t *testing.T) {
	configured := []map[string]string{
		{
			"name":  "gitlab_email_from",
			"value": "test_value_01",
		},
		{
			"name":  "gitlab_email_reply_to",
			"value": "test_value_02",
		},
		{
			"name":  "smtp_password",
			"value": "test_value_03",
		},
		{
			"name":  "smtp_user_name",
			"value": "test_value_04",
		},
		{
			"name":  "",
			"value": "test_value_05",
		},
		{
			"name":  "omniauth_providers_saml_name_2",
			"value": "test_value_06",
		},
	}

	tests := []struct {
		name string
		args []map[string]string
		want *types.RequestParameters
	}{
		{
			name: "expands the resource data",
			args: configured,
			want: &types.RequestParameters{
				GitlabEmailFrom:            nifcloud.String("test_value_01"),
				GitlabEmailReplyTo:         nifcloud.String("test_value_02"),
				SmtpPassword:               nifcloud.String("test_value_03"),
				SmtpUserName:               nifcloud.String("test_value_04"),
				OmniauthProvidersSamlName2: nifcloud.String("test_value_06"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateParameterGroupParameters(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandParameters(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"parameter": []interface{}{
			map[string]interface{}{
				"name":  "test_name_01",
				"value": "test_value_01",
			},
			map[string]interface{}{
				"name":  "test_name_02",
				"value": "test_value_02",
			},
			map[string]interface{}{
				"value": "test_value_03",
			},
		},
	})

	tests := []struct {
		name string
		args []interface{}
		want []types.Parameters
	}{
		{
			name: "expands the resource data",
			args: rd.Get("parameter").(*schema.Set).List(),
			want: []types.Parameters{
				{
					Name:  nifcloud.String("test_name_02"),
					Value: nifcloud.String("test_value_02"),
				},
				{
					Name:  nifcloud.String("test_name_01"),
					Value: nifcloud.String("test_value_01"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandParameters(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
