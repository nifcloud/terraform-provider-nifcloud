package devopsbackuprule

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateBackupRuleInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"instance_id": "test_id",
		"description": "test_description",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.CreateBackupRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.CreateBackupRuleInput{
				BackupRuleName: nifcloud.String("test_name"),
				InstanceId:     nifcloud.String("test_id"),
				Description:    nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateBackupRuleInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateBackupRuleInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name_changed",
		"description": "test_description",
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.UpdateBackupRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.UpdateBackupRuleInput{
				BackupRuleName:        nifcloud.String("test_name"),
				ChangedBackupRuleName: nifcloud.String("test_name_changed"),
				Description:           nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateBackupRuleInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetBackupRulesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.GetBackupRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.GetBackupRuleInput{
				BackupRuleName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetBackupRuleInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteBackupRuleInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.DeleteBackupRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.DeleteBackupRuleInput{
				BackupRuleName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteBackupRuleInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
