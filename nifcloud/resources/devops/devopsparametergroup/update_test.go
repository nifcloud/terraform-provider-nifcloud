package devopsparametergroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/stretchr/testify/assert"
)

func TestGetParametersToUpdate(t *testing.T) {
	params := []map[string]string{
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
			"value": "********",
		},
		{
			"name":  "smtp_user_name",
			"value": "test_value_04",
		},
	}
	os := schema.NewSet(schema.HashResource(newSchema()["parameter"].Elem.(*schema.Resource)), []interface{}{
		map[string]interface{}{
			"name":  "gitlab_email_from",
			"value": "test_value_01",
		},
		map[string]interface{}{
			"name":  "gitlab_email_reply_to",
			"value": "test_value_02",
		},
		map[string]interface{}{
			"name":  "smtp_password",
			"value": "test_value_03",
		},
		map[string]interface{}{
			"name":  "smtp_user_name",
			"value": "test_value_04",
		},
	})
	ns := schema.NewSet(schema.HashResource(newSchema()["parameter"].Elem.(*schema.Resource)), []interface{}{
		map[string]interface{}{
			"name":  "gitlab_email_from",
			"value": "test_value_01",
		},
		map[string]interface{}{
			"name":  "gitlab_email_reply_to",
			"value": "test_value_102",
		},
		map[string]interface{}{
			"name":  "smtp_password",
			"value": "test_value_103",
		},
	})

	type args struct {
		params []map[string]string
		os     *schema.Set
		ns     *schema.Set
	}
	tests := []struct {
		name string
		args args
		want *types.RequestParameters
	}{
		{
			name: "expands the resource data",
			args: args{params, os, ns},
			want: &types.RequestParameters{
				GitlabEmailFrom:    nifcloud.String("test_value_01"),
				GitlabEmailReplyTo: nifcloud.String("test_value_102"),
				SmtpPassword:       nifcloud.String("test_value_103"),
				SmtpUserName:       nifcloud.String(""),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getParametersToUpdate(tt.args.params, tt.args.os, tt.args.ns)
			assert.Equal(t, tt.want, got)
		})
	}
}
